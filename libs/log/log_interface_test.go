package log

import (
	"testing"
)

var testValue string = ""

type testLogger struct {
	minLevel byte
}

func (tl *testLogger) Print(prefix, content string) {
	testValue = prefix + content
}
func (tl *testLogger) Panic(prefix, content string) {
	testValue = prefix + content
}
func (tl *testLogger) Fatal(prefix, content string) {
	testValue = prefix + content
}
func (tl *testLogger) MinLevel() byte {
	return tl.minLevel
}
func (tl *testLogger) SetMinLevel(n byte) {
	tl.minLevel = n
}

// TestLoggerBaseOutPut test the method of interface
// and base of all function for Logger interface.
func TestLoggerBaseOutPut(t *testing.T) {
	tl := &testLogger{minLevel: DebugLevel}

	//MinLevel
	if ml := tl.MinLevel(); ml != DebugLevel {
		t.Errorf("MinLevel: the minLevel should be 1, but get %d.\n", ml)
	}

	//SetMinLevel
	tl.SetMinLevel(FatalLevel)
	t.Log(tl.MinLevel())
	if ml := tl.MinLevel(); ml != FatalLevel {
		t.Error(FatalLevel)
		t.Errorf("SetMinLevel: the minLevel should be 4, but get %d.\n", ml)
	}
	tl.SetMinLevel(DebugLevel)

	//Print
	tl.Print("1", "1")
	if testValue != "11" {
		t.Errorf("Print: the testValue should be \"11\", but get %s.\n", testValue)
	}

	//Panic
	tl.Panic("2", "2")
	if testValue != "22" {
		t.Errorf("Panic: the testValue should be \"22\", but get %s.\n", testValue)
	}

	//Fatal
	tl.Fatal("3", "3")
	if testValue != "33" {
		t.Errorf("Fatal: the testValue should be \"33\", but get %s.\n", testValue)
	}

	testValue = ""
	//Debugf
	Debugf(tl, "1%s%d", "1", 1)
	if testValue != "[Debug]111" {
		t.Errorf("Debugf: the testValue should be \"[Debug]111\", but get %s.\n", testValue)
	}
	//Warnf
	Warnf(tl, "1%s%d", "1", 1)
	if testValue != "[Warn ]111" {
		t.Errorf("Warnf: the testValue should be \"[Warn ]111\", but get %s.\n", testValue)
	}
	//Debugf
	Errorf(tl, "1%s%d", "1", 1)
	if testValue != "[Error]111" {
		t.Errorf("Errorf: the testValue should be \"[Error]111\", but get %s.\n", testValue)
	}
	//Panicf
	Panicf(tl, "1%s%d", "1", 1)
	if testValue != "[Panic]111" {
		t.Errorf("Panicf: the testValue should be \"[Panic]111\", but get %s.\n", testValue)
	}
	//Fatalf
	Fatalf(tl, "1%s%d", "1", 1)
	if testValue != "[Fatal]111" {
		t.Errorf("Fatalf: the testValue should be \"[Fatal]111\", but get %s.\n", testValue)
	}

}

// TestLoggerLevelFilter test the feature of level filter
// of all function for Logger interface.
func TestLoggerLevelFilter(t *testing.T) {
	tl := &testLogger{minLevel: DebugLevel}

	//DebugLevel
	testValue = ""
	Debugln(tl, "222")
	if testValue != "[Debug]222\n" {
		t.Errorf("Debugf: the testValue should be \"[Debug]222\\n\", but get %s.\n", testValue)
	}
	Warnln(tl, "222")
	if testValue != "[Warn ]222\n" {
		t.Errorf("Warnf: the testValue should be \"[Warn ]222\\n\", but get %s.\n", testValue)
	}
	Errorln(tl, "222")
	if testValue != "[Error]222\n" {
		t.Errorf("Errorf: the testValue should be \"[Error]222\\n\", but get %s.\n", testValue)
	}
	Panicln(tl, "222")
	if testValue != "[Panic]222\n" {
		t.Errorf("Panicf: the testValue should be \"[Panic]222\\n\", but get %s.\n", testValue)
	}
	Fatalln(tl, "222")
	if testValue != "[Fatal]222\n" {
		t.Errorf("Fatalf: the testValue should be \"[Fatal]222\\n\", but get %s.\n", testValue)
	}

	//WarnLevel
	tl.SetMinLevel(WarnLevel)
	testValue = ""
	Debugln(tl, "222")
	if testValue == "[Debug]222\n" {
		t.Errorf("Debugf: the testValue should not be \"[Debug]222\\n\"\n")
	}
	Warnln(tl, "222")
	if testValue != "[Warn ]222\n" {
		t.Errorf("Warnf: the testValue should be \"[Warn ]222\\n\", but get %s.\n", testValue)
	}
	Errorln(tl, "222")
	if testValue != "[Error]222\n" {
		t.Errorf("Errorf: the testValue should be \"[Error]222\\n\", but get %s.\n", testValue)
	}
	Panicln(tl, "222")
	if testValue != "[Panic]222\n" {
		t.Errorf("Panicf: the testValue should be \"[Panic]222\\n\", but get %s.\n", testValue)
	}
	Fatalln(tl, "222")
	if testValue != "[Fatal]222\n" {
		t.Errorf("Fatalf: the testValue should be \"[Fatal]222\\n\", but get %s.\n", testValue)
	}

	//ErrorLevel
	tl.SetMinLevel(ErrorLevel)
	testValue = ""
	Debugln(tl, "222")
	if testValue == "[Debug]222\n" {
		t.Errorf("Debugf: the testValue should not be \"[Debug]222\\n\"\n")
	}
	Warnln(tl, "222")
	if testValue == "[Warn ]222\n" {
		t.Errorf("Warnf: the testValue should not be \"[Warn ]222\\n\"\n")
	}
	Errorln(tl, "222")
	if testValue != "[Error]222\n" {
		t.Errorf("Errorf: the testValue should be \"[Error]222\\n\", but get %s.\n", testValue)
	}
	Panicln(tl, "222")
	if testValue != "[Panic]222\n" {
		t.Errorf("Panicf: the testValue should be \"[Panic]222\\n\", but get %s.\n", testValue)
	}
	Fatalln(tl, "222")
	if testValue != "[Fatal]222\n" {
		t.Errorf("Fatalf: the testValue should be \"[Fatal]222\\n\", but get %s.\n", testValue)
	}

	//PanicLevel
	tl.SetMinLevel(PanicLevel)
	testValue = ""
	Debugln(tl, "222")
	if testValue == "[Debug]222\n" {
		t.Errorf("Debugf: the testValue should not be \"[Debug]222\\n\"\n")
	}
	Warnln(tl, "222")
	if testValue == "[Warn ]222\n" {
		t.Errorf("Warnf: the testValue should not be \"[Warn ]222\\n\"\n")
	}
	Errorln(tl, "222")
	if testValue == "[Error]222\n" {
		t.Errorf("Errorf: the testValue should not be \"[Error]222\\n\"\n")
	}
	Panicln(tl, "222")
	if testValue != "[Panic]222\n" {
		t.Errorf("Panicf: the testValue should be \"[Panic]222\\n\", but get %s.\n", testValue)
	}
	Fatalln(tl, "222")
	if testValue != "[Fatal]222\n" {
		t.Errorf("Fatalf: the testValue should be \"[Fatal]222\\n\", but get %s.\n", testValue)
	}

	//FatalLevel
	tl.SetMinLevel(FatalLevel)
	testValue = ""
	Debugln(tl, "222")
	if testValue == "[Debug]222\n" {
		t.Errorf("Debugf: the testValue should not be \"[Debug]222\\n\"\n")
	}
	Warnln(tl, "222")
	if testValue == "[Warn ]222\n" {
		t.Errorf("Warnf: the testValue should not be \"[Warn ]222\\n\"\n")
	}
	Errorln(tl, "222")
	if testValue == "[Error]222\n" {
		t.Errorf("Errorf: the testValue should not be \"[Error]222\\n\"\n")
	}
	Panicln(tl, "222")
	if testValue == "[Panic]222\n" {
		t.Errorf("Panicf: the testValue should not be \"[Panic]222\\n\"\n")
	}
	Fatalln(tl, "222")
	if testValue != "[Fatal]222\n" {
		t.Errorf("Fatalf: the testValue should be \"[Fatal]222\\n\", but get %s.\n", testValue)
	}
}
