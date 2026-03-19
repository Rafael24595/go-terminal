package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Rafael24595/go-terminal/engine/commons"
	"github.com/Rafael24595/go-terminal/engine/commons/file"
	"github.com/Rafael24595/go-terminal/engine/platform/clock"
)

const loggerBootstrap Logger = "Bootstrap"

type bootstrapLogger struct {
	mu        sync.RWMutex
	clock     clock.Clock
	timestamp int64
	records   []Record
}

func newBootstrapLogger() *bootstrapLogger {
	return &bootstrapLogger{
		clock:     clock.UnixMilliClock,
		timestamp: clock.UnixMilliClock(),
		records:   make([]Record, 0),
	}
}

func (l *bootstrapLogger) Name() Logger {
	return loggerBootstrap
}

func (l *bootstrapLogger) Records() []Record {
	l.mu.RLock()
	defer l.mu.RUnlock()

	out := make([]Record, len(l.records))
	copy(out, l.records)

	return out
}

func (l *bootstrapLogger) Custom(category string, message string) Record {
	upperCategory := strings.ToUpper(category)
	return l.write(Category(upperCategory), message)
}

func (l *bootstrapLogger) Custome(category string, err error) Record {
	return l.Custom(category, err.Error())
}

func (l *bootstrapLogger) Customf(category string, format string, a ...any) Record {
	return l.Custom(category, fmt.Sprintf(format, a...))
}

func (l *bootstrapLogger) Message(message string) Record {
	return l.write(MESSAGE, message)
}

func (l *bootstrapLogger) Messagef(format string, a ...any) Record {
	return l.Message(fmt.Sprintf(format, a...))
}

func (l *bootstrapLogger) Warning(message string) Record {
	return l.write(WARNING, message)
}

func (l *bootstrapLogger) Warningf(format string, a ...any) Record {
	return l.Warning(fmt.Sprintf(format, a...))
}

func (l *bootstrapLogger) Error(err error) Record {
	return l.Errors(err.Error())
}

func (l *bootstrapLogger) Errors(message string) Record {
	return l.write(ERROR, message)
}

func (l *bootstrapLogger) Errorf(format string, a ...any) Record {
	return l.Errors(fmt.Sprintf(format, a...))
}

func (l *bootstrapLogger) Record(records ...Record) []Record {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.records = append(l.records, records...)

	return records
}

func (l *bootstrapLogger) Close() []Record {
	if len(l.records) == 0 {
		return l.records
	}
	
	l.mu.Lock()
	defer l.mu.Unlock()

	jsonData, err := json.MarshalIndent(l.records, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}

	name := fmt.Sprintf("log-unsigned-%s", commons.FormatMillisecondsCompact(l.timestamp))
	path := fmt.Sprintf("%s/%s.json", defaultPath, name)

	file.WriteFileSafe(path, string(jsonData))

	records := l.records
	l.records = make([]Record, 0)

	return records
}

func (l *bootstrapLogger) write(category Category, message string) Record {
	l.mu.Lock()
	defer l.mu.Unlock()

	record := Record{
		Category:  category,
		Message:   message,
		Timestamp: l.clock(),
	}

	l.records = append(l.records, record)

	return record
}

func (l *bootstrapLogger) flush(logger Log) {
	l.mu.Lock()
	records := l.records
	l.records = make([]Record, 0)
	l.mu.Unlock()

	logger.Record(records...)
}
