//Package log provides a power but simple logger
package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// levels
const (
	InfoLevel = iota
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel

	infoPrefix  = "[Info ]"
	warnPrefix  = "[Warn ]"
	errorPrefix = "[Error]"
	panicPrefix = "[Panic]"
	fatalPrefix = "[Fatal]"
)

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

// getInvokerLocation get filename and line according to skipNumber.
func getInvokerLocation(skipNumber int) string {
	//Get the file and line of the invoker
	_, file, line, ok := runtime.Caller(skipNumber)
	if !ok {
		return ""
	}

	var simpleFileName string
	//Only get the file basename(the same as os.path.basename in python3)
	if index := strings.LastIndex(file, "/"); index > 0 {
		simpleFileName = file[index+1 : len(file)]
	} else {
		simpleFileName = file
	}

	return fmt.Sprintf("%s:%d", simpleFileName, line)
}

// generateLogContent create log content.
// its formate is "prefix date clock filePosition - levelPrefix"
//
// pos (position) is the relative position where calling log method relatives to
// generateLogContent in function invocation stack.
func generateLogContent(
	pos uint,
	levelPrefix,
	format string,
	v ...interface{}) string {

	t := time.Now()
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	//TODO: to be faster!
	// fmt.Sprintf is too slow, to use more underline api.
	// -----
	//calculate the aim function position on the calling goroutine's stack
	skipNumber := int(pos) + 2
	baseInfo := fmt.Sprintf(" %4d/%02d/%02d %02d:%02d:%02d %s - %s ",
		year, month, day,
		hour, min, sec,
		getInvokerLocation(skipNumber),
		levelPrefix,
	)
	// -----

	var result string
	if len(format) > 0 {
		//generate accroding to format
		result = fmt.Sprintf((baseInfo + format), v...)
	} else {
		//generate directly
		vLen := len(v)
		params := make([]interface{}, (vLen + 1))
		params[0] = baseInfo
		copy(params[1:], v)
		result = fmt.Sprintln(params...)
	}
	return result
}
