# **Go errors**

## Overview

The main motto of the module is to handle errors as first-class citizens in Go. It provides the way to attach inspectable context to the error (mainly by merging several errors into an instance of module’s internal error type) and to provide utility functions to inspect the errors for the attached context. Main focus was on providing the way to make errors testable and to convey some meaningful information for the upper layers on error propagation. 

This module was inspired by the github.com/pkg/errors module which solves some error handling issues in Go, but still is hardly usable for the unit testing and has constraints for working with typed errors. 

All errors returned from the utility functions **Wrap**, **WithMessage** and **WrapWithMessage** return an object conforming to the **error** interface with a stacktrace created at the moment of function invocation. Use:

```go
err := errors.New("error message")
fmt.Printf("%s", err) // Prints the error message
fmt.Printf("%+v", err) // Prints the error message with stacktrace
```

The error message consists of all error messages joined with " : " sequence. 

## **Installation**

```
go get github.com/ameteiko/errors
```

## **New(msg string) error**

This function is a replacement for the **errors.New** function to prevent from importing several packages for error handling. 

## **Wrap(err ...error) error**

Utility function Wrap merges several errors passed as a variadic parameter into one. All nil errors will be filtered out of the resulting error. Best usage is to attach an application error to the error returned from the 3rd-party module.

```go
var errValidation = errors.New("validation error")

func parseResponse(r []byte) error {
    var resp Response
    if err := json.Unmarshal(r, &resp); err != nil {
        return errors.Wrap(err, errValidation)
    }
}
```

## WithMessage(err error, format string, args ...interface{}) error

Attaches a formatted message for the error. If message is empty, the original error will be returned. If error is nil, then nil will be returned. Best used to add some error context information to the error, like parameters that resulted with an error. 

```go
func parseResponse(r []byte) error {
    var resp Response
    if err := json.Unmarshal(r, &resp); err != nil {
        return errors.WithMessage(err, "unmarshalling error for %v", r)
    }
}
```



## WrapWithMessage(erra, errb error, format string, args ...interface{}) error

Wraps several errors and attaches a message to them. This function is a combination of Wrap and WithMessage. 

```go
var errValidation = errors.New("validation error")

func parseResponse(r []byte) error {
    var resp Response
    if err := json.Unmarshal(r, &resp); err != nil {
        return errors.WrapWithMessage(err, errValidation, "unmarshalling error for %v", r)
    }
}
```



## Fetch(source error, target error) error

Inspects the source error and returns matched target error object from it. If target is nil, or a pointer to nil, then result is nil.

```go
parseErr := parseResponse
if err := errors.Fetch(parseErr, errValidation); err != nil {
    // There was a parsing error.
}
```

## FetchByType(source error, target interface{}) error

Inspects the error and returns the entries matched by the type. Returns nil if If target is nil or not a pointer.

```go
type LoggableErr struct {
    msg string
}

func (e LoggableErr) Error() string {
    return e.msg
}

func (e LoggableErr) Log() string {
    log.Print(e.msg)
}

errResponseParsing := LoggableErr{msg: "response parsing error"}
  
func main() {
    parsingErr := parseResponse("}")
    if err := errors.FetchByType(parsingErr, (*LoggableErr)(nil)); err != nil {
        loggableErr := err.(LoggableErr)
        loggableErr.Log()
    }
}
```

