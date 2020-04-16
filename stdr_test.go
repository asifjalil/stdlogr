package stdlogr

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"testing"
)

func TestInfo(t *testing.T) {
	const testString = "test\ntest"
	var b bytes.Buffer
	l := New(log.New(&b, "", 0))
	l.Info(testString)
	if expect := fmt.Sprintf("[verbosity=0] %s\n", testString); b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestQuoteInfo(t *testing.T) {
	const testString = "test\ntest"
	var b bytes.Buffer
	l := Logger{
		Std: log.New(&b, "", 0),
		Formatter: DefaultFormatter{
			ForceQuote: true,
		},
	}
	l.Info(testString)
	if expect := fmt.Sprintf("[verbosity=0] %q\n", testString); b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestError(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := New(log.New(&b, "", 0))
	l.Error(errors.New(testString), testString)
	if expect := fmt.Sprintf("[Error=%[1]s] [verbosity=0] %[1]s\n", testString); b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestName(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := New(log.New(&b, "", 0))
	l.WithName(testString).Info(testString)
	if expect := fmt.Sprintf("[name=%[1]s] [verbosity=0] %[1]s\n", testString); b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestNameAppend(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := New(log.New(&b, "", 0))
	l.WithName(testString).WithName(testString).Info(testString)
	if expect := fmt.Sprintf("[name=%[1]s/%[1]s] [verbosity=0] %[1]s\n", testString); b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestInfoKV(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := New(log.New(&b, "", 0))
	l.Info(testString, testString, testString)
	if expect := fmt.Sprintf("[verbosity=0] [%[1]s=%[1]v] %[1]s\n", testString); b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestVerbosity(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := New(log.New(&b, "", 0))
	l.V(1).Info(testString, testString, testString)
	if b.String() != "" {
		t.Errorf("Expected 0 output but got %q", b.String())
	}
}

func TestTimestampFormat(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer

	l := Logger{
		Std: log.New(&b, "", 0),
		Formatter: DefaultFormatter{
			TimestampFormat: "2006-01-02 03:04:05 -0700",
		},
	}
	l.Info(testString)
	t.Log(b.String())
}
