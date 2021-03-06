

# stdlogr
`import "github.com/asifjalil/stdlogr"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package stdlogr  implements github.com/go-logr/logr.Logger in terms of
Go's standard log package.

Copied from <a href="https://github.com/go-logr/stdr">https://github.com/go-logr/stdr</a>




## <a name="pkg-index">Index</a>
* [func New(std StdLogger) logr.Logger](#New)
* [func SetVerbosity(v int) int](#SetVerbosity)
* [type DefaultFormatter](#DefaultFormatter)
  * [func (f DefaultFormatter) Format(e Entry) string](#DefaultFormatter.Format)
* [type Entry](#Entry)
* [type Formatter](#Formatter)
* [type Logger](#Logger)
  * [func (l Logger) Enabled() bool](#Logger.Enabled)
  * [func (l Logger) Error(err error, msg string, kvList ...interface{})](#Logger.Error)
  * [func (l Logger) Info(msg string, kvList ...interface{})](#Logger.Info)
  * [func (l Logger) V(level int) logr.InfoLogger](#Logger.V)
  * [func (l Logger) WithName(name string) logr.Logger](#Logger.WithName)
  * [func (l Logger) WithValues(kvList ...interface{}) logr.Logger](#Logger.WithValues)
* [type StdLogger](#StdLogger)


#### <a name="pkg-files">Package files</a>
[stdr.go](/src/github.com/asifjalil/stdlogr/stdr.go) 





## <a name="New">func</a> [New](/src/target/stdr.go?s=5569:5604#L199)
``` go
func New(std StdLogger) logr.Logger
```
New returns a logr.Logger which is implemented by Go's standard log package.

Example: stdr.New(log.New(os.Stderr, "", log.LstdFlags))

You can also just instantiate your own:


	var log = stdr.Logger{
	        Std: log.New(os.Stderr, "", log.LstdFlags),
	        Depth: 0
	        Formatter: stdr.DefaultFormatter{},
	}



## <a name="SetVerbosity">func</a> [SetVerbosity](/src/target/stdr.go?s=2071:2099#L73)
``` go
func SetVerbosity(v int) int
```
SetVerbosity sets the global level against which all info logs will be
compared.  If this is greater than or equal to the "V" of the logger, the
message will be logged.  A higher value here means more logs will be written.
The previous verbosity value is returned.

SetVerbosity is concurrent-safe.




## <a name="DefaultFormatter">type</a> [DefaultFormatter](/src/target/stdr.go?s=2145:2592#L77)
``` go
type DefaultFormatter struct {
    // TimestampFormat sets the format used to print log timestamp.
    // If TimestampFormat is set, make sure timestamp flags are off
    // for the std. logger. Otherwise there will be duplicate timestamp
    // in the output.
    TimestampFormat string
    // ForceQuote true quotes all fieldKeys and fieldValues using %q.
    ForceQuote bool
    // HideKeys true will show [fieldValue] instead of [fieldKey=fieldValue]
    HideKeys bool
}
```









### <a name="DefaultFormatter.Format">func</a> (DefaultFormatter) [Format](/src/target/stdr.go?s=2617:2665#L90)
``` go
func (f DefaultFormatter) Format(e Entry) string
```
Format a log entry.




## <a name="Entry">type</a> [Entry](/src/target/stdr.go?s=350:1297#L22)
``` go
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
```
Entry is a log entry.
It is used by `Formatter` to log format output.










## <a name="Formatter">type</a> [Formatter](/src/target/stdr.go?s=1339:1389#L43)
``` go
type Formatter interface {
    Format(Entry) string
}
```
Formatter formats a log entry Entry.










## <a name="Logger">type</a> [Logger](/src/target/stdr.go?s=4469:5197#L170)
``` go
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
    // contains filtered or unexported fields
}
```









### <a name="Logger.Enabled">func</a> (Logger) [Enabled](/src/target/stdr.go?s=7517:7547#L275)
``` go
func (l Logger) Enabled() bool
```



### <a name="Logger.Error">func</a> (Logger) [Error](/src/target/stdr.go?s=7607:7674#L279)
``` go
func (l Logger) Error(err error, msg string, kvList ...interface{})
```



### <a name="Logger.Info">func</a> (Logger) [Info](/src/target/stdr.go?s=7234:7289#L264)
``` go
func (l Logger) Info(msg string, kvList ...interface{})
```



### <a name="Logger.V">func</a> (Logger) [V](/src/target/stdr.go?s=5680:5724#L206)
``` go
func (l Logger) V(level int) logr.InfoLogger
```



### <a name="Logger.WithName">func</a> (Logger) [WithName](/src/target/stdr.go?s=6021:6070#L215)
``` go
func (l Logger) WithName(name string) logr.Logger
```
WithName returns a new logr.Logger with the specified name appended.  stdr
uses '/' characters to separate name elements.  Callers should not pass '/'
in the provided name string, but this library does not actually enforce that.




### <a name="Logger.WithValues">func</a> (Logger) [WithValues](/src/target/stdr.go?s=6183:6244#L224)
``` go
func (l Logger) WithValues(kvList ...interface{}) logr.Logger
```



## <a name="StdLogger">type</a> [StdLogger](/src/target/stdr.go?s=4334:4467#L165)
``` go
type StdLogger interface {
    // Output is the same as log.Output and log.Logger.Output.
    Output(calldepth int, logline string) error
}
```
StdLogger is the subset of the Go stdlib log.Logger API that is needed for
this adapter.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
