package log

import (
	"barrage-server/libs/color"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

var (
	coloredInfo  = color.Dye(color.FgGreen, infoPrefix)
	coloredWarn  = color.Dye(color.FgYellow, warnPrefix)
	coloredError = color.Dye(color.FgRed, errorPrefix)
	coloredPanic = color.Dye(color.FgHiRed, panicPrefix)
	coloredFatal = color.Dye(color.BgRed, fatalPrefix)
)

// TwoOutputLogger has two output, one is used to normal output and the other is used to error output,
// it is goroutine safe.
type TwoOutputLogger struct {
	minLevel byte
	m        sync.RWMutex
	buf      []byte // for accumulating text to write
	out      io.Writer
	err      io.Writer
}

// NewStdLogger create and return a new TwoOutputLogger which implements Logger,
// wrapped io.Stdout.
//
// all logs should be writed by logger, so before creating logger, any error will panic,
// so all New** function won't return error
func NewStdLogger(level byte) *TwoOutputLogger {
	if level > FatalLevel {
		panic("Your min level of logger is too high!")
	}
	return &TwoOutputLogger{
		minLevel: level,
		out:      os.Stdout,
		err:      os.Stderr,
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
	filename := fmt.Sprintf("%s_%d%02d%02d_%02d_%02d.log",
		prefix,
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

	return &TwoOutputLogger{
		minLevel: level,
		out:      w,
		err:      w,
	}, w
}

// NewSimpleLogger create and return a new TwoOutputLogger which implements Logger,
// this Logger use one io.Writer as output and error output.
func NewSimpleLogger(w io.Writer, level byte) *TwoOutputLogger {
	if level > FatalLevel {
		panic("Your min level of logger is too high!")
	}
	if w == nil {
		panic("Your w io.Writer should not be nil.")
	}

	return &TwoOutputLogger{
		minLevel: level,
		out:      w,
		err:      w,
	}
}

// NewLogger create and return a new TwoOutputLogger which implements Logger.
func NewLogger(out io.Writer, err io.Writer, level byte) *TwoOutputLogger {
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
		minLevel: level,
		out:      out,
		err:      err,
	}
}

// levelCheck check print whether is bigger than minLevel, if it is true return ture
func (l *TwoOutputLogger) levelCheck(printLevel byte) bool {
	if l.minLevel <= printLevel {
		return true
	}

	return false
}

// Infof print info level log.
func (l *TwoOutputLogger) Infof(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(InfoLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredInfo)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
}

// Infoln print info level log.
func (l *TwoOutputLogger) Infoln(v ...interface{}) {
	s := fmt.Sprintln(v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(InfoLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredInfo)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
}

// Fatalf print fatal level log.
func (l *TwoOutputLogger) Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(FatalLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredFatal)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

// Fatalln print fatal level log.
func (l *TwoOutputLogger) Fatalln(v ...interface{}) {
	s := fmt.Sprintln(v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(FatalLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredFatal)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

// Panicf print panic level log.
func (l *TwoOutputLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(PanicLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredPanic)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
	panic(s)
}

// Panicln print panic level log.
func (l *TwoOutputLogger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(PanicLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredPanic)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
	panic(s)
}

// Errorf print error level log.
func (l *TwoOutputLogger) Errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(ErrorLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredError)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
}

// Errorln print error level log.
func (l *TwoOutputLogger) Errorln(v ...interface{}) {
	s := fmt.Sprintln(v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(ErrorLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredError)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
}

// Warnf print warn level log.
func (l *TwoOutputLogger) Warnf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(WarnLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredWarn)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
}

// Warnln print warn level log.
func (l *TwoOutputLogger) Warnln(v ...interface{}) {
	s := fmt.Sprintln(v...)

	l.m.Lock()
	defer l.m.Unlock()
	if !l.levelCheck(WarnLevel) {
		return
	}

	generateLogHead(&l.buf, 1, coloredWarn)
	l.buf = append(l.buf, s...)

	_, err := l.out.Write(l.buf)
	if err != nil {
		panic(err)
	}
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
