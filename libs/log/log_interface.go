//Package log provides a power but simple logger
package log

import (
	"fmt"
	"log"
)

// levels
const (
	DebugLevel = iota
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
	defaultFlag = log.LstdFlags | log.Lshortfile

	debugPrefix = "[Debug]"
	warnPrefix  = "[Warn ]"
	errorPrefix = "[Error]"
	panicPrefix = "[Panic]"
	fatalPrefix = "[Fatal]"
)

var levelMap = map[byte]string{
	DebugLevel: debugPrefix,
	WarnLevel:  warnPrefix,
	ErrorLevel: errorPrefix,
	PanicLevel: panicPrefix,
	FatalLevel: fatalPrefix,
}

// Logger defines the base interface of logger.
// logger should be open while the program running, so it has no Close method.
//
// its base feature:
//   * add log level
//   * filter low level logs
type Logger interface {
	// Print base log like debug, warn, error
	Print(level byte, content string)
	// Panic is equivalent to Print() followed by a call to panic().
	Panic(content string)
	// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
	Fatal(content string)
	// MinLevel return the minimize level logger should print.
	MinLevel() byte
	// SetMinLevel set the minLevel of logger.
	SetMinLevel(byte)
}

// Debugf print debug level logs.
//
// Pure Debug function or other type print log function is useless, so it does not
// define them.
func Debugf(logger Logger, format string, v ...interface{}) {
	if logger.MinLevel() <= DebugLevel {
		logger.Print(DebugLevel, fmt.Sprintf(format, v...))
	}
}

// Debugln print debug level logs.
func Debugln(logger Logger, v ...interface{}) {
	if logger.MinLevel() <= DebugLevel {
		logger.Print(DebugLevel, fmt.Sprintln(v...))
	}
}

// Warnf print warn level logs.
func Warnf(logger Logger, format string, v ...interface{}) {
	if logger.MinLevel() <= WarnLevel {
		logger.Print(WarnLevel, fmt.Sprintf(format, v...))
	}
}

// Warnln print warn level logs.
func Warnln(logger Logger, v ...interface{}) {
	if logger.MinLevel() <= WarnLevel {
		logger.Print(WarnLevel, fmt.Sprintln(v...))
	}
}

// Errorf print error level logs.
func Errorf(logger Logger, format string, v ...interface{}) {
	if logger.MinLevel() <= ErrorLevel {
		logger.Print(ErrorLevel, fmt.Sprintf(format, v...))
	}
}

// Errorln print error level logs.
func Errorln(logger Logger, v ...interface{}) {
	if logger.MinLevel() <= ErrorLevel {
		logger.Print(ErrorLevel, fmt.Sprintln(v...))
	}
}

// Panicf print panic level logs.
func Panicf(logger Logger, format string, v ...interface{}) {
	if logger.MinLevel() <= PanicLevel {
		logger.Panic(fmt.Sprintf(format, v...))
	}
}

// Panicln print panic level logs.
func Panicln(logger Logger, v ...interface{}) {
	if logger.MinLevel() <= PanicLevel {
		logger.Panic(fmt.Sprintln(v...))
	}
}

// Fatalf print fatal level logs.
func Fatalf(logger Logger, format string, v ...interface{}) {
	if logger.MinLevel() <= FatalLevel {
		logger.Fatal(fmt.Sprintf(format, v...))
	}
}

// Fatalln print fatal level logs.
func Fatalln(logger Logger, v ...interface{}) {
	if logger.MinLevel() <= FatalLevel {
		logger.Fatal(fmt.Sprintln(v...))
	}
}
