package local

import (
	"io"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/model/record"
)

func WriterErrorHandler(w io.Writer, f func() error) {
	err := f();
	if  err == nil {
		return
	}

	_, err = w.Write([]byte(err.Error()))
	if  err == nil {
		return
	}

	println(err.Error())
}

func LogErrorHandler(f func() error) {
	WriterErrorHandler(log.WriterFromCategory(record.ERROR), f)
}
