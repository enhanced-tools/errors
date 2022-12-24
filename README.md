# Enhanced Error

The library that allows boost up debugging abilities in your app with customizable errors with stack traces and parameters.

## Installation

```bash
go get -u github.com/enhanced-tools/errors
```

## Usage

```go
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/enhanced-tools/errors"
	"github.com/enhanced-tools/errors/opts"
)

// Define error templates as global variables
var (
	ErrNotANumber = errors.Template().With(
		opts.Title("Not a number"),
	)
)

func main() {
	if len(os.Args) != 3 {
		// Create error inline and log it
		errors.New("wrong number of arguments").Log()
		return
	}
	v1, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		// Create enhanced error from common error and log it
		errors.Enhance(err).Log()
		return
	}
	v2, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		// Add option that describes better the error
		errors.Enhance(err).With(errors.Title("Bad second argument")).Log()
		return
	}
	log.Printf("%.2f + %f = %f", v1, v2, v1+v2)
}
```

## Predefined errors

To create predefined errors you can use `errors.Template` function. It returns a template that can be used to create errors with the same options.

```go
var (
    ErrNotANumber = errors.Template().With(
        opts.Title("Not a number"),
        opts.Status(http.StatusBadRequest),
    )
)
```
And use it like this 
```go
if err := strconv.ParseFloat(os.Args[1], 64); err != nil {
    ErrNotANumber.From(err).Log()
    return
}
```
You can also create it without existing error
```go
if err := strconv.ParseFloat(os.Args[1], 64); err != nil {
    ErrNotANumber.FromEmpty().Log()
    return
}
```
`FromEmpty()` is required in order to save stack trace from the place where error was created.

## Available options

- `Title: string` - error title
- `StatusCode: int` - HTTP status code
- `Debug: any` - debug value to include
- `Type: string` - type for error
- `RequestID: string` - request id

You can add any option to error using `With` method

```go
errors.New("some error").With(errors.Title("Some title"))
```

Yoy can also declare your own options by implementing `ErrorOpt` interface


```go
type ErrorOpt interface {
	// Type returns the type of the option
	Type() string
	// MapFormatter returns the map representation of the option
	MapFormatter() map[string]interface{}
	// Verbosity returns the verbosity of the option
	Verbosity() int
}
```

```go
type MyOption struct {
    Value string
}

func (o MyOption) Type() string {
    return "my-option"
}

func (o MyOption) MapFormatter() map[string]interface{} {
    return map[string]interface{}{
        "value": o.Value,
    }
}

func (o MyOption) Verbosity() int {
    return 0
}
```

## Logging

To log an error you can use `Log` method

```go
errors.New("some error").Log()
```

`Log(... LogName)` takes optional log names. If you don't specify it, it will use `default` log name.

To register new logger use `errors.RegisterLogger` function
```go
errors.Manager().RegisterLogger("custom-name", loggerFunc)
```
Where `LoggerFunc` satisfies
```go
type LoggerFunc func(err EnhancedError)
```
To use this logger just envoke 
```go
errors.New("some error").Log("custom-name")
```
You can call multiple loggers in `Log` function.

You can create your custom logger from components using `CustomLogger`

```go
errors.Manager().RegisterLogger(
	// name for the logger
	"custom-name",
    errors.CustomLogger(
		// log error in LogFMT format
		errors.WithErrorFormatter(errors.LogFMTFormatter),
		// print no stack trace
		errors.WithStackTraceFormatter(errors.NoStackTrace),
		// but include stack trace in file
		errors.WithSaveStack(true),
        // set verbosity for messages (only 30 or below verbosity options will be printed out)
		errors.WithVerbosity(30),
	)
)
```
To replace default logger you can use
```
errors.Manager().SetDefaultLogger(yourLoggerImplementation)
```
See examples for more info

## Manager instance

Manger is an singleton managing errors and stack traces. You can get it using `errors.Manager()` function. It has the following methods available:

```go
type ErrorsManager interface {
	// SaveStack saves the stack trace of the error in the stack trace file. Format is optional
	SaveStack(err EnhancedError, format ...StackTraceFormatter) error
	// RegisterLogger registers a logger for a specific log name
	RegisterLogger(name LogName, logger LoggerFunc)
	// SetDefaultLogger sets the default logger
	SetDefaultLogger(logger LoggerFunc)
	// Setup creates a new stack trace file and reads the existing one
	Setup(stackTracePath string) error
}
```
To save stack traces you need first to Setup the manager with the path to the stack trace file. You can use `errors.Setup` function to do it.  
After that you can use `SaveStack` method to save stack traces. You can also customize loggers to save it while logging.

Example output of the stack trace file
```
>>> 1f56e93b2cf89835ce9f1b33f2d88662
	github.com/enhanced-tools/errors.Enhance
	/Users/vachingmachine/projects/enhanced-errors/errors.go:70
	main.getIntVar
	/Users/vachingmachine/projects/enhanced-errors/example/server/main.go:54
	main.main.func2
	/Users/vachingmachine/projects/enhanced-errors/example/server/main.go:97
	net/http.HandlerFunc.ServeHTTP
	/usr/local/go/src/net/http/server.go:2084
	github.com/go-chi/chi/v5.(*Mux).routeHTTP
	/Users/vachingmachine/go/pkg/mod/github.com/go-chi/chi/v5@v5.0.8/mux.go:444
	net/http.HandlerFunc.ServeHTTP
	/usr/local/go/src/net/http/server.go:2084
	main.requestID.func1
	/Users/vachingmachine/projects/enhanced-errors/example/server/main.go:41
	net/http.HandlerFunc.ServeHTTP
	/usr/local/go/src/net/http/server.go:2084
	github.com/go-chi/chi/v5/middleware.RequestLogger.func1.1
	/Users/vachingmachine/go/pkg/mod/github.com/go-chi/chi/v5@v5.0.8/middleware/logger.go:54
	net/http.HandlerFunc.ServeHTTP
	/usr/local/go/src/net/http/server.go:2084
	github.com/go-chi/chi/v5.(*Mux).ServeHTTP
	/Users/vachingmachine/go/pkg/mod/github.com/go-chi/chi/v5@v5.0.8/mux.go:90
	net/http.serverHandler.ServeHTTP
	/usr/local/go/src/net/http/server.go:2916
	net/http.(*conn).serve
	/usr/local/go/src/net/http/server.go:1966
	runtime.goexit
	/usr/local/go/src/runtime/asm_arm64.s:1263
```

