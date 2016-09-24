package log

import (
	// "errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

// TwoOutputLogger has two output, one is used to normal output and the other is used to error output,
// it is goroutine safe.
type TwoOutputLogger struct {
	minLevel  byte
	m         sync.RWMutex
	outLogger *log.Logger
	errLogger *log.Logger
}

// NewStdLogger create and return a new TwoOutputLogger which implements Logger,
// wrapped log.Logger is created with io.Stdout.
//
// all logs should be writed by logger, so before creating logger, any error will panic,
// so all New** function won't return error
func NewStdLogger(prefix string, level byte) *TwoOutputLogger {
	if level > FatalLevel {
		panic("Your min level of logger is too high!")
	}
	return &TwoOutputLogger{
		minLevel:  level,
		outLogger: log.New(os.Stdout, prefix, defaultFlag),
		errLogger: log.New(os.Stderr, prefix, defaultFlag),
	}
}

// NewSimpleFileLogger  create and return a new TwoOutputLogger which implements Logger and
// the file created for write logs(this file pinter should only be used for testing!).
//
// it will open file in 'dirname/date.log', then use this file as output and error output.
// the date is the time of program up
func NewSimpleFileLogger(dirname string, prefix string, level byte) (*TwoOutputLogger, *os.File) {
	if dirname == "" {
		panic("Your dirname should not be \"\"!")
	}

	now := time.Now()
	filename := fmt.Sprintf("%d%02d%02d_%02d_%02d.log",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute())
	path := path.Join(dirname, filename)

	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}

	l := log.New(w, prefix, defaultFlag)

	return &TwoOutputLogger{
		minLevel:  level,
		outLogger: l,
		errLogger: l,
	}, w
}

// NewSimpleLogger create and return a new TwoOutputLogger which implements Logger,
// this Logger use one io.Writer as output and error output.
func NewSimpleLogger(w io.Writer, prefix string, level byte) *TwoOutputLogger {
	if level > FatalLevel {
		panic("Your min level of logger is too high!")
	}
	if w == nil {
		panic("Your w io.Writer should not be nil.")
	}

	l := log.New(w, prefix, defaultFlag)

	return &TwoOutputLogger{
		minLevel:  level,
		outLogger: l,
		errLogger: l,
	}
}

// NewLogger create and return a new TwoOutputLogger which implements Logger.
func NewLogger(out io.Writer, err io.Writer, prefix string, level byte) *TwoOutputLogger {
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
	return &TwoOutputLogger{
		minLevel:  level,
		outLogger: log.New(out, prefix, defaultFlag),
		errLogger: log.New(err, prefix, defaultFlag),
	}
}

// Print print base log like debug, warn, error.
func (l *TwoOutputLogger) Print(level byte, content string) {
	if level < ErrorLevel {
		l.outLogger.Print(levelMap[level], content)
	} else {
		l.errLogger.Print(levelMap[level], content)
	}
}

// Panic is equivalent to Print() followed by a call to panic().
// just call logger.Panic
func (l *TwoOutputLogger) Panic(content string) {
	l.errLogger.Panic(panicPrefix, content)
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
// just call logger.Fatal
func (l *TwoOutputLogger) Fatal(content string) {
	l.errLogger.Fatal(fatalPrefix, content)
}

// MinLevel return the minimize level logger should print.
func (l *TwoOutputLogger) MinLevel() byte {
	l.m.RLock()
	defer l.m.RUnlock()

	return l.minLevel
}

// SetMinLevel set the minLevel of logger.
func (l *TwoOutputLogger) SetMinLevel(level byte) {
	l.m.Lock()
	defer l.m.Unlock()

	if level <= FatalLevel {
		l.minLevel = level
	}
}
