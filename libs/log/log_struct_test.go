package log

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogColor(t *testing.T) {
	t.Log(levelMap[InfoLevel])
	t.Log(levelMap[WarnLevel])
	t.Log(levelMap[ErrorLevel])
	t.Log(levelMap[PanicLevel])
	t.Log(levelMap[FatalLevel])
}

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
	info := "testing_info 1 2 3\n"
	if tbs := testBuffer.String(); !strings.Contains(tbs, info) || !strings.Contains(tbs, infoPrefix) {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", info, tbs)
	}

	testBuffer.Reset()
	l.Warnln("testing_warn", 1, 2, 3)
	warn := "testing_warn 1 2 3\n"
	if tbs := testBuffer.String(); !strings.Contains(tbs, warn) || !strings.Contains(tbs, warnPrefix) {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", warn, tbs)
	}

	testBuffer.Reset()
	l.Errorln("testing_error", 1, 2, 3)
	error := "testing_error 1 2 3\n"
	if tbs := testBuffer.String(); !strings.Contains(tbs, error) || !strings.Contains(tbs, errorPrefix) {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", error, tbs)
	}

	testBuffer.Reset()
	l.Infof("testing_info %d %d %d\n", 1, 2, 3)
	if tbs := testBuffer.String(); !strings.Contains(tbs, info) || !strings.Contains(tbs, infoPrefix) {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", info, tbs)
	}

	testBuffer.Reset()
	l.Warnf("testing_warn %d %d %d\n", 1, 2, 3)
	if tbs := testBuffer.String(); !strings.Contains(tbs, warn) || !strings.Contains(tbs, warnPrefix) {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", warn, tbs)
	}

	testBuffer.Reset()
	l.Errorf("testing_error %d %d %d\n", 1, 2, 3)
	if tbs := testBuffer.String(); !strings.Contains(tbs, error) || !strings.Contains(tbs, errorPrefix) {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", error, tbs)
	}
}

func TestNewSimpleFileLogger(t *testing.T) {
	l, w := NewSimpleFileLogger("./", "Testing", InfoLevel)
	warn := "testing_warn\n"

	l.Warnln("testing_warn")

	filename := w.Name()
	t.Logf("logger file name: %s.", filename)
	w.Close()

	w, _ = os.Open(filename)
	wReader := bufio.NewReader(w)
	if tbs, _ := wReader.ReadString('\n'); !strings.Contains(tbs, warn) || !strings.Contains(tbs, warnPrefix) {
		t.Errorf("Print: the end of printed string should be %s, but get %s.", warn, tbs)
	}

	w.Close()
	os.Remove(filename)
}
