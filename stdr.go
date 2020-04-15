// Package stdlogr  implements github.com/go-logr/logr.Logger in terms of
// Go's standard log package.
//
// Copied from https://github.com/go-logr/stdr
//
package stdlogr

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
)

// Entry is a log entry.
// It is used by `Formatter` to log format output.
type Entry struct {
	// Err is an error being logged using `logr Error(err error ...)` method.
	Err error
	// Name is logger name. It constists of a series of name "segments"
	// added by successive calls to `logr WithName(name string)` method.
	Name string
	// Verbosity represents how little a log matters.
	// Level zero, the default, matters most. Increasing levels matter less and less.
	// Verbosity is specified using `logr V(level int)` method.
	Verbosity int
	// Message consists of a constant message attached to the the log line.
	// This should generally be a simple description of what's occuring,
	// and should never be a format string.
	// Message is spcified using `logr Info` method and `logr Error` method.
	Message string
	// KeyAndValues where keys are arbitrary strings, while values may be any Go value.
	// KeyAndValues specified using `logr WithValues`, `logr Info`, and `logr Error` method.
	KeysAndValues []interface{}
}

// Formatter formats a log entry Entry.
type Formatter interface {
	Format(Entry) string
}

type verbosity struct {
	mu sync.Mutex //ensures atomic writes; protects the following fields
	v  int        // v is verbosity of this log package
}

func (vb *verbosity) setVerbosity(v int) int {
	vb.mu.Lock()
	old := vb.v
	vb.v = v
	vb.mu.Unlock()
	return old
}

func (vb *verbosity) getVerbosity() int {
	return vb.v
}

var globalVerbosity = new(verbosity)

// SetVerbosity sets the global level against which all info logs will be
// compared.  If this is greater than or equal to the "V" of the logger, the
// message will be logged.  A higher value here means more logs will be written.
// The previous verbosity value is returned.
//
// SetVerbosity is concurrent-safe.
//
func SetVerbosity(v int) int {
	return globalVerbosity.setVerbosity(v)
}

type DefaultFormatter struct {
	// TimestampFormat sets the format used to print log timestamp.
	// If TimestampFormat is set, make sure timestamp flags are off
	// for the std. logger. Otherwise there will be duplicate timestamp
	// in the output.
	TimestampFormat string
	HideKeys        bool // show [fieldValue] instead [fieldKey=fieldValue]
}

// Format a log entry.
func (f DefaultFormatter) Format(e Entry) string {
	var b strings.Builder

	if f.TimestampFormat != "" {
		fmt.Fprintf(&b, "%s ", time.Now().Format(f.TimestampFormat))
	}
	if e.Err != nil {
		fmt.Fprintf(&b, "[Error=%v] ", e.Err)
	}
	if e.Name != "" {
		fmt.Fprintf(&b, "[name=%s] ", e.Name)
	}
	fmt.Fprintf(&b, "[verbosity=%d]", e.Verbosity)

	// write fields (keys/values)
	if len(e.KeysAndValues) > 0 {
		b.WriteString(" ")
	}
	f.writeFields(&b, e.KeysAndValues)
	// log message
	fmt.Fprintf(&b, " %s\n", e.Message)
	return b.String()
}

func (f DefaultFormatter) writeFields(b *strings.Builder, kvList []interface{}) {
	keys := make([]string, 0, len(kvList))
	vals := make(map[string]interface{}, len(kvList))
	for i := 0; i < len(kvList); i += 2 {
		k, ok := kvList[i].(string)
		if !ok {
			fmt.Fprintf(b, "**key is not a string: %[1]v(type=%[1]T**)", kvList[i])
			return
		}
		var v interface{}
		if i+1 < len(kvList) {
			v = kvList[i+1]
		}
		keys = append(keys, k)
		vals[k] = v
	}
	sort.Strings(keys)
	for i, k := range keys {
		v := vals[k]
		if i > 0 {
			b.WriteString(" ")
		}
		if f.HideKeys {
			fmt.Fprintf(b, "[%v]", v)
		} else {
			fmt.Fprintf(b, "[%s=%v]", k, v)
		}
	}
	return
}

// StdLogger is the subset of the Go stdlib log.Logger API that is needed for
// this adapter.
type StdLogger interface {
	// Output is the same as log.Output and log.Logger.Output.
	Output(calldepth int, logline string) error
}

type Logger struct {
	Std StdLogger
	// DepthOffset biases the assumed number of call frames to the "true"
	// caller.  This is useful when the calling code calls a function which then
	// calls glogr (e.g. a logging shim to another API).  Values less than zero
	// will be treated as zero.
	Depth int
	// All log entries pass through the formatter before logged to Out. The
	// included formatter is `DefaultFormatter`.
	// You can easily implement your Formatter implements the `Formatter` interface, see the
	// or included formatters for example.
	Formatter Formatter

	verbosity int           // set using logr.V(n)
	prefix    string        // set using logr.WithName
	values    []interface{} // set using logr.WithValues
}

// New returns a logr.Logger which is implemented by Go's standard log package.
//
// Example: stdr.New(log.New(os.Stderr, "", log.LstdFlags))
//
// You can also just instantiate your own:
//
//      var log = stdr.Logger{
//              Std: log.New(os.Stderr, "", log.LstdFlags),
//              Depth: 0
//              Formatter: stdr.DefaultFormatter{},
//      }
func New(std StdLogger) logr.Logger {
	return Logger{
		Std:       std,
		Formatter: DefaultFormatter{},
	}
}

func (l Logger) V(level int) logr.InfoLogger {
	new := l.clone()
	new.verbosity = level
	return new
}

// WithName returns a new logr.Logger with the specified name appended.  stdr
// uses '/' characters to separate name elements.  Callers should not pass '/'
// in the provided name string, but this library does not actually enforce that.
func (l Logger) WithName(name string) logr.Logger {
	new := l.clone()
	if len(l.prefix) > 0 {
		new.prefix = l.prefix + "/"
	}
	new.prefix += name
	return new
}

func (l Logger) WithValues(kvList ...interface{}) logr.Logger {
	new := l.clone()
	new.values = append(new.values, kvList...)
	return new
}

func (l Logger) clone() Logger {
	return Logger{
		Std:       l.Std,
		Depth:     l.Depth,
		Formatter: l.Formatter,
		verbosity: l.verbosity,
		prefix:    l.prefix,
		values:    copySlice(l.values),
	}
}

func copySlice(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	copy(out, in)
	return out
}

// Magic string for intermediate frames that we should ignore.
const autogeneratedFrameName = "<autogenerated>"

// Discover how many frames we need to climb to find the caller. This approach
// was suggested by Ian Lance Taylor of the Go team, so it *should* be safe
// enough (famous last words).
func framesToCaller() int {
	// 1 is the immediate caller.  3 should be too many.
	for i := 1; i < 3; i++ {
		_, file, _, _ := runtime.Caller(i + 1) // +1 for this function's frame
		if file != autogeneratedFrameName {
			return i
		}
	}
	return 1 // something went wrong, this is safe
}

func (l Logger) Info(msg string, kvList ...interface{}) {
	if l.Enabled() {
		l.output(framesToCaller()+l.Depth,
			l.Formatter.Format(Entry{
				Name:          l.prefix,
				Verbosity:     l.verbosity,
				Message:       msg,
				KeysAndValues: append(l.values, kvList...)}))
	}
}

func (l Logger) Enabled() bool {
	return globalVerbosity.getVerbosity() >= l.verbosity
}

func (l Logger) Error(err error, msg string, kvList ...interface{}) {
	l.output(framesToCaller()+l.Depth,
		l.Formatter.Format(Entry{
			Err:           err,
			Name:          l.prefix,
			Verbosity:     l.verbosity,
			Message:       msg,
			KeysAndValues: append(l.values, kvList...)}))
}

func (l Logger) output(calldepth int, s string) {
	depth := calldepth + 2 // offset for this adapter

	// ignore errors - what can we really do about them?
	if l.Std != nil {
		_ = l.Std.Output(depth, s)
	} else {
		_ = log.Output(depth, s)
	}
}
