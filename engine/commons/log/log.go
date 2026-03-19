package log

import (
	"errors"
	"sync"
)

type Logger string

const defaultPath = ".logs"

var (
	log  Log = newBootstrapLogger()
	once sync.Once
)

type Provider interface {
	Build() (Log, error)
}

func InitFromProvider(provider Provider) error {
	lg, err := provider.Build()
	if err != nil {
		return err
	}
	return InitFromLog(lg)
}

func InitFromLog(lg Log) error {
	if lg == nil {
		return errors.New("nil logger")
	}

	init := false

	once.Do(func() {
		old := log
		if b, ok := old.(bootstrap); ok {
			b.flush(lg)
		}

		log = lg
		log.Messagef("Logging is configured to use the %s instance.", lg.Name())
		old.Close()

		init = true
	})

	if !init {
		return errors.New("logger already initialized")
	}

	return nil
}

func OnClose() {
	log.Close()
}

type bootstrap interface {
	flush(Log)
}

type Log interface {
	Name() Logger
	Records() []Record
	Custom(string, string) Record
	Custome(string, error) Record
	Customf(string, string, ...any) Record
	Message(string) Record
	Messagef(string, ...any) Record
	Warning(string) Record
	Warningf(string, ...any) Record
	Error(error) Record
	Errors(string) Record
	Errorf(string, ...any) Record
	Record(...Record) []Record
	Close() []Record
}

func Name() Logger {
	return log.Name()
}

func Records() []Record {
	return log.Records()
}

func Custom(category string, message string) {
	log.Custom(category, message)
}

func Custome(category string, err error) {
	log.Custome(category, err)
}

func Customf(category string, format string, args ...any) {
	log.Customf(category, format, args...)
}

func Message(message string) {
	log.Message(message)
}

func Messagef(format string, args ...any) {
	log.Messagef(format, args...)
}

func Warning(message string) {
	log.Warning(message)
}

func Warningf(format string, args ...any) {
	log.Warningf(format, args...)
}

func Error(err error) {
	log.Error(err)
}

func Errors(message string) {
	log.Errors(message)
}

func Errorf(format string, args ...any) {
	log.Errorf(format, args...)
}
