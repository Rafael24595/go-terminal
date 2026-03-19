package log

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Rafael24595/go-terminal/engine/commons"
	"github.com/Rafael24595/go-terminal/engine/commons/file"
	"github.com/Rafael24595/go-terminal/engine/platform/clock"
)

const LoggerFile Logger = "File"

const defaultBuff = 100

const ArgFilePath = "LOGGER_FILE_PATH"
const ArgFileBuff = "LOGGER_FILE_BUFF"

type FileProvider struct {
	Session    string
	Path       string
	BufferSize int
}

func (p FileProvider) Build() (Log, error) {
	if p.Path == "" {
		p.Path = defaultPath
	}

	if p.BufferSize <= 0 {
		p.BufferSize = defaultBuff
	}

	return newFileLogger(p.Session, p.Path, p.BufferSize)
}

type fileLogger struct {
	mu sync.RWMutex

	ch     chan Record
	done   chan struct{}
	closed atomic.Bool

	file *file.File

	session   string
	timestamp int64
	clock     clock.Clock

	format  Format
	records []Record
}

func newFileLogger(session string, path string, buff int) (*fileLogger, error) {
	timestamp := clock.UnixMilliClock()

	name := fmt.Sprintf("log-%s-%s", session, commons.FormatMillisecondsCompact(timestamp))
	source := fmt.Sprintf("%s/%s.%s", path, name, DefaultFormat.Extension)

	file, err := file.NewFile(source, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	logger := &fileLogger{
		ch:        make(chan Record, buff),
		done:      make(chan struct{}),
		file:      file,
		session:   session,
		clock:     clock.UnixMilliClock,
		timestamp: timestamp,
		format:    DefaultFormat,
		records:   make([]Record, 0),
	}

	go logger.loop()

	return logger, nil
}

func (l *fileLogger) Session(session string) *fileLogger {
	l.session = session
	return l
}

func (l *fileLogger) Formatter(format Format) *fileLogger {
	l.format = format
	return l
}

func (l *fileLogger) Name() Logger {
	return LoggerFile
}

func (l *fileLogger) Records() []Record {
	l.mu.RLock()
	defer l.mu.RUnlock()

	out := make([]Record, len(l.records))
	copy(out, l.records)

	return out
}

func (l *fileLogger) Custom(category string, message string) Record {
	upperCategory := strings.ToUpper(category)
	return l.write(Category(upperCategory), message)
}

func (l *fileLogger) Custome(category string, err error) Record {
	return l.Custom(category, err.Error())
}

func (l *fileLogger) Customf(category string, format string, a ...any) Record {
	return l.Custom(category, fmt.Sprintf(format, a...))
}

func (l *fileLogger) Message(message string) Record {
	return l.write(MESSAGE, message)
}

func (l *fileLogger) Messagef(format string, a ...any) Record {
	return l.Message(fmt.Sprintf(format, a...))
}

func (l *fileLogger) Warning(message string) Record {
	return l.write(WARNING, message)
}

func (l *fileLogger) Warningf(format string, a ...any) Record {
	return l.Warning(fmt.Sprintf(format, a...))
}

func (l *fileLogger) Error(err error) Record {
	return l.Errors(err.Error())
}

func (l *fileLogger) Errors(message string) Record {
	return l.write(ERROR, message)
}

func (l *fileLogger) Errorf(format string, a ...any) Record {
	return l.Errors(fmt.Sprintf(format, a...))
}

func (l *fileLogger) Record(records ...Record) []Record {
	for _, record := range records {
		l.ch <- record
	}
	return records
}

func (l *fileLogger) Close() []Record {
	if l.closed.CompareAndSwap(false, true) {
		close(l.ch)
		<-l.done
	}
	return l.records
}

func (l *fileLogger) write(category Category, message string) Record {
	record := Record{
		Category:  category,
		Message:   message,
		Timestamp: l.clock(),
	}

	if l.closed.Load() {
		return record
	}

	l.ch <- record

	return record
}

func (l *fileLogger) loop() {
	for record := range l.ch {
		l.mu.Lock()
		l.records = append(l.records, record)
		l.mu.Unlock()

		data := l.format.Format(record)
		_ = l.file.Append([]byte(data))
	}

	close(l.done)
}
