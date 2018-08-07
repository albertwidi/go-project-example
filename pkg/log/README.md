# Log

Log is based on [Uber zap](https://github.com/uber-go/zap) log package

## Log level

Log level is supported in the logger, available log level is:

- Debug
- Info
- Warning
- Error
- Fatal

Log is disabled if `LogLevel` < `CurrentLogLevel`, for example `Debug` log is disabled when current level is `Info`

Example of `SetLevel`:

```go
import "github.com/albertwidi/kothak/pkg/log"

log.SetLevel(log.InfoLevel)
log.Infow("this is a log", "key1", "val1")
```

## Log to file

All logs are written to `stderr`, but we can also write the log to file by using:

```go
import "github.com/albertwidi/kothak/pkg/log"

log.SetLevel(log.InfoLevel)
err := log.SetOutputTofile("logfile.log")
if err != nil {
    panic(err)
}
log.Infow("this is a log", "key1", "val1")
```

## Key-value context in log

To add more context to log, `key-value` fields is provided. For example:

```go
import "github.com/albertwidi/kothak/pkg/log"

log.Infow("this is a log", "key1", "val1")
```

or

```go
import "github.com/albertwidi/kothak/pkg/log"

log.With("key1", "value1", "key2", "value2").Info("this is a log)
```

## Integration with Error package

Error package has a features called `errors.Fields`. This fields can be used to add more context into the error, and then we can print the fields when needed. TDK log will automatically print the fields if `error = tdkerrors.Error` by using `log.Errors`. For example:

```go
import "github.com/albertwidi/kothak/log"
import "github.com/albertwidi/kothak/errors"

err := errors.E("this is an error", errors.Fields{"field1":"value1"})
log.Errors(err)

// result is
// message=this is an error field1=value1
```