package log

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestLoggerMinLevelAndSetMinLevel(t *testing.T) {
	l := NewStdLogger("Testing", DebugLevel)

	if level := l.MinLevel(); level != DebugLevel {
		t.Errorf("MinLevel: the minLevel of logger should be %d, but get %d.", DebugLevel, level)
	}

	l.SetMinLevel(ErrorLevel)
	if level := l.MinLevel(); level != ErrorLevel {
		t.Errorf("SetMinLevel: the minLevel of logger should be %d, but get %d.", ErrorLevel, level)
	}
}

func TestLoggerPrint(t *testing.T) {
	var testBuffer bytes.Buffer
	l := NewSimpleLogger(&testBuffer, "Testing", DebugLevel)

	l.Print(DebugLevel, "testing_debug")
	debug := "[Debug]testing_debug\n"
	if tbs := testBuffer.String()[46:]; tbs != debug {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", debug, tbs)
	}

	testBuffer.Reset()
	l.Print(WarnLevel, "testing_warn")
	warn := "[Warn ]testing_warn\n"
	if tbs := testBuffer.String()[46:]; tbs != warn {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", warn, tbs)
	}

	testBuffer.Reset()
	l.Print(ErrorLevel, "testing_error")
	error := "[Error]testing_error\n"
	if tbs := testBuffer.String()[46:]; tbs != error {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", error, tbs)
	}
}

func TestNewSimpleFileLogger(t *testing.T) {
	l, w := NewSimpleFileLogger("./", "Testing", DebugLevel)
	warn := "[Warn ]testing_warn\n"

	l.Print(WarnLevel, "testing_warn")

	filename := w.Name()
	t.Logf("logger file name: %s.", filename)
	w.Close()

	w, _ = os.Open(filename)
	wReader := bufio.NewReader(w)
	if tbs, _ := wReader.ReadString('\n'); tbs[46:] != warn {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", warn, tbs[46:])
	}

	w.Close()
	os.Remove(filename)
}
