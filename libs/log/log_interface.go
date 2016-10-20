//Package log provides a power but simple logger
package log

import (
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
func getInvokerLocation(skipNumber int) (simpleFileName string, line int) {
	//Get the file and line of the invoker
	_, file, line, ok := runtime.Caller(skipNumber)
	if !ok {
		return
	}

	//Only get the file basename(the same as os.path.basename in python3)
	if index := strings.LastIndex(file, "/"); index > 0 {
		simpleFileName = file[index+1 : len(file)]
	} else {
		simpleFileName = file
	}

	return
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// generateLogHead create log Head.
// its formate is " date clock filePosition - levelPrefix "
//
// pos (position) is the relative position where calling log method relatives to
// generateLogHead in function invocation stack.
func generateLogHead(buf *[]byte, pos uint, levelPrefix string) {
	*buf = (*buf)[:0]

	*buf = append(*buf, ' ')

	// timestamp
	t := time.Now()
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)

	*buf = append(*buf, ' ')

	// filename and line
	filename, line := getInvokerLocation(int(pos) + 2)
	*buf = append(*buf, filename...)
	*buf = append(*buf, ':')
	itoa(buf, line, -1)

	*buf = append(*buf, " - "...)
	*buf = append(*buf, levelPrefix...)
	*buf = append(*buf, ' ')
}
