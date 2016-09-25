package log

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestLoggerMinLevelAndSetMinLevel(t *testing.T) {
	l := NewStdLogger("Testing", InfoLevel)

	if level := l.MinLevel(); level != InfoLevel {
		t.Errorf("MinLevel: the minLevel of logger should be %d, but get %d.", InfoLevel, level)
	}

	l.SetMinLevel(ErrorLevel)
	if level := l.MinLevel(); level != ErrorLevel {
		t.Errorf("SetMinLevel: the minLevel of logger should be %d, but get %d.", ErrorLevel, level)
	}
}

func TestLoggerFormat(t *testing.T) {
	var testBuffer bytes.Buffer
	l := NewSimpleLogger(&testBuffer, "Testing", InfoLevel)

	l.Infoln("testing_info", 1, 2, 3)
	info := "[Info ]testing_info 1 2 3\n"
	if tbs := testBuffer.String()[46:]; tbs != info {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", info, tbs)
	}

	testBuffer.Reset()
	l.Warnln("testing_warn", 1, 2, 3)
	warn := "[Warn ]testing_warn 1 2 3\n"
	if tbs := testBuffer.String()[46:]; tbs != warn {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", warn, tbs)
	}

	testBuffer.Reset()
	l.Errorln("testing_error", 1, 2, 3)
	error := "[Error]testing_error 1 2 3\n"
	if tbs := testBuffer.String()[46:]; tbs != error {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", error, tbs)
	}

	testBuffer.Reset()
	l.Infof("testing_info %d %d %d\n", 1, 2, 3)
	if tbs := testBuffer.String()[46:]; tbs != info {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", info, tbs)
	}

	testBuffer.Reset()
	l.Warnf("testing_warn %d %d %d\n", 1, 2, 3)
	if tbs := testBuffer.String()[46:]; tbs != warn {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", warn, tbs)
	}

	testBuffer.Reset()
	l.Errorf("testing_error %d %d %d\n", 1, 2, 3)
	if tbs := testBuffer.String()[46:]; tbs != error {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", error, tbs)
	}
}

func TestNewSimpleFileLogger(t *testing.T) {
	l, w := NewSimpleFileLogger("./", "Testing", InfoLevel)
	warn := "[Warn ]testing_warn\n"

	l.Warnln("testing_warn")

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
