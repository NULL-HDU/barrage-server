package log

import (
	// "errors"
	"io"
	"log"
	"os"
)

type logger struct {
	minLevel  byte
	outLogger *log.Logger
	errLogger *log.Logger
}

// NewStdLogger create and return a new Logger which is implemented by logger,
// wrapped log.Logger is created with io.Stdout.
//
// all logs should be writed by logger, so before creating logger, any error will panic,
// so all New** function won't return error
func NewStdLogger(prefix string, level byte) Logger {
	if level > FatalLevel {
		panic("Your min level of logger is too high!")
	}
	return &logger{
		minLevel:  level,
		outLogger: log.New(os.Stdout, prefix, defaultFlag),
		errLogger: log.New(os.Stderr, prefix, defaultFlag),
	}
}

// NewSimpleFileLogger create and return a new Logger which is implemented by logger,
// it will open file in 'path', then use this file as output and error output.
func NewSimpleFileLogger(path string, prefix string, level byte) Logger {
	if path == "" {
		panic("Your path should not be \"\"!")
	}
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}

	l := log.New(w, prefix, defaultFlag)

	return &logger{
		minLevel:  level,
		outLogger: l,
		errLogger: l,
	}
}

// NewSimpleLogger create and return a new Logger which is implemented by logger,
// this Logger use one io.Writer as output and error output.
func NewSimpleLogger(w io.Writer, prefix string, level byte) Logger {
	if level > FatalLevel {
		panic("Your min level of logger is too high!")
	}
	if w == nil {
		panic("Your w io.Writer should not be nil.")
	}

	l := log.New(w, prefix, defaultFlag)

	return &logger{
		minLevel:  level,
		outLogger: l,
		errLogger: l,
	}
}

// NewLogger create and return a new Logger which is implemented by logger.
func NewLogger(out io.Writer, err io.Writer, prefix string, level byte) Logger {
	if level > FatalLevel {
		panic("Your min level of logger is too high!")
	}
	if out == nil && level < ErrorLevel {
		panic("Your output io.Writer is nil, " +
			"but your min level of logger is bigger than ErrorLevel!")
	}
	if err == nil {
		panic("Your error output io.Writer should not be nil.")
	}
	return &logger{
		minLevel:  level,
		outLogger: log.New(out, prefix, defaultFlag),
		errLogger: log.New(err, prefix, defaultFlag),
	}
}

// Print print base log like debug, warn, error.
func (l *logger) Print(level byte, content string) {
	if level < ErrorLevel {
		l.outLogger.Print(levelMap[level], content)
	} else {
		l.errLogger.Print(levelMap[level], content)
	}
}
