//Package log provides a power but simple logger
package log

import (
	"barrage-server/libs/color"
	"log"
)

// levels
const (
	InfoLevel = iota
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
	defaultFlag = log.LstdFlags | log.Lshortfile

	infoPrefix  = "[Info ]"
	warnPrefix  = "[Warn ]"
	errorPrefix = "[Error]"
	panicPrefix = "[Panic]"
	fatalPrefix = "[Fatal]"
)

var levelMap = map[byte]string{
	InfoLevel:  color.Dye(color.FgGreen, infoPrefix),
	WarnLevel:  color.Dye(color.FgYellow, warnPrefix),
	ErrorLevel: color.Dye(color.FgRed, errorPrefix),
	PanicLevel: color.Dye(color.FgHiRed, panicPrefix),
	FatalLevel: color.Dye(color.BgRed, fatalPrefix),
}

// Logger defines the base interface of logger.
// logger should be open while the program running, so it has no Close method.
//
// its base feature:
//   * add log level
//   * filter low level logs
type Logger interface {
	// Infof print info level log.
	Infof(format string, v ...interface{})
	// Warnf print warn level log.
	Warnf(format string, v ...interface{})
	// Errorf print error level log.
	Errorf(format string, v ...interface{})
	// Panic is equivalent to Print() followed by a call to panic().
	Panicf(format string, v ...interface{})
	// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
	Fatalf(format string, v ...interface{})
	// Infoln print info level log.
	Infoln(v ...interface{})
	// Warnln print wran level log.
	Warnln(v ...interface{})
	// Errorln print error level log.
	Errorln(v ...interface{})
	// Panic is equivalent to Print() followed by a call to panic().
	Panicln(v ...interface{})
	// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
	Fatalln(v ...interface{})
	// MinLevel return the minimize level logger should print.
	MinLevel() byte
	// SetMinLevel set the minLevel of logger.
	SetMinLevel(byte)
}
