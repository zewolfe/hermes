package log

import (
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

type Logger = logr.Logger

func NewStdoutLogger() *Logger {
	newlogger := log.New(os.Stdout, "", log.Lshortfile)
	stdLogger := stdr.New(newlogger)

	return &stdLogger
}

func NewNoopLogger() *Logger {
	l := logr.Discard()
	return &l
}
