package errors

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ErrorOptType string

type ErrorOpt interface {
	// Type returns the type of the option. It must be unique withing single enhanced error. If not, the last option will be used.
	Type() ErrorOptType
	// MapFormatter returns a map that will be used to format the error. The map will be merged with the map returned by the error.
	MapFormatter() map[string]interface{}
	// Verbosity returns the verbosity of the option. The higher the value, the more specific the option is. The default value is 0.
	Verbosity() int
}

type enhancedError struct {
	error
	TemplateID string
	ErrorID    string
	Opts       map[ErrorOptType]ErrorOpt
}

type EnhancedError interface {
	error
	// With adds an option to the error. If the option already exists (by checking its type), it will be overwritten.
	With(opts ...ErrorOpt) EnhancedError
	// From returns a enhanced error from common one. If the error is already enhanced, it will be returned as copy of previous one.
	From(err error) EnhancedError
	// FromEmpty returns a new enhanced error from nothing. It is used with Template() function to mark proper stack trace.
	FromEmpty() EnhancedError

	// Log logs the error using the default logger or the loggers specified in the function. Loggers must be registered before either the program will panic.
	Log(loggers ...LogName)
	// Is checks if the error is of the same type as the one specified. It will check the template ID and the error ID if comparing enhanced errors.
	Is(err error) bool
	// Wrap wraps the error with a message. Is is used as a shortcut for With(Wrapper(msg))
	Wrap(msg string) EnhancedError

	// GetStackTrace returns the stack trace of the error.
	GetStackTrace() errors.StackTrace
	// GetStackTraceHash returns the hash of the stack trace of the error.
	GetStackTraceHash() string
	// GetInternalError returns the internal error.
	GetInternalError() error
	// GetOpts returns the options of the error.
	GetOpts() map[ErrorOptType]ErrorOpt
	// GetOpt sets error opt value
	GetOpt(opt ErrorOpt) bool
	// GetErrorID returns the error ID of the error.
	GetErrorID() string
}

func copyOpts(opts map[ErrorOptType]ErrorOpt) map[ErrorOptType]ErrorOpt {
	targetOpts := make(map[ErrorOptType]ErrorOpt)
	for k, v := range opts {
		targetOpts[k] = v
	}
	return targetOpts
}

func Enhance(err error) EnhancedError {
	if err == nil {
		return nil
	}
	if enErr, ok := err.(EnhancedError); ok {
		return enErr
	}
	return &enhancedError{
		ErrorID: uuid.NewString(),
		error:   errors.WithStack(err),
		Opts:    make(map[ErrorOptType]ErrorOpt),
	}
}

func New(msg string) EnhancedError {
	return &enhancedError{
		ErrorID: uuid.NewString(),
		error:   errors.New(msg),
		Opts:    make(map[ErrorOptType]ErrorOpt),
	}
}

func Newf(msg string, data ...interface{}) EnhancedError {
	return &enhancedError{
		ErrorID: uuid.NewString(),
		error:   errors.New(fmt.Sprintf(msg, data...)),
		Opts:    make(map[ErrorOptType]ErrorOpt),
	}
}

func (e enhancedError) With(opts ...ErrorOpt) EnhancedError {
	newOpts := copyOpts(e.Opts)
	for _, opt := range opts {
		if existingWrapper, ok := e.Opts["wrapper"]; opt.Type() == "wrapper" && ok {
			existingWrapper := existingWrapper.(Wrapper)
			newWrapper := opt.(Wrapper)
			newOpts["wrapper"] = newWrapper + existingWrapper
		} else {
			newOpts[opt.Type()] = opt
		}
	}
	return &enhancedError{
		ErrorID:    uuid.NewString(),
		TemplateID: e.TemplateID,
		error:      e.error,
		Opts:       newOpts,
	}
}

func Wrap(err error, msg string) EnhancedError {
	return Enhance(err).Wrap(msg)
}

func Wrapf(err error, msg string, data ...interface{}) EnhancedError {
	return Wrap(err, fmt.Sprintf(msg, data...))
}

func Template() EnhancedError {
	return &enhancedError{
		TemplateID: uuid.NewString(),
		Opts:       make(map[ErrorOptType]ErrorOpt),
	}
}

func (e enhancedError) From(err error) EnhancedError {
	if enErr, ok := err.(*enhancedError); ok {
		return &enhancedError{
			ErrorID:    uuid.NewString(),
			TemplateID: e.TemplateID,
			error:      enErr.error,
			Opts:       copyOpts(e.Opts),
		}
	}
	return &enhancedError{
		ErrorID:    uuid.NewString(),
		TemplateID: e.TemplateID,
		error:      errors.WithStack(err),
		Opts:       copyOpts(e.Opts),
	}
}

func (e enhancedError) FromEmpty() EnhancedError {
	return e.From(fmt.Errorf("error"))
}

func (e enhancedError) Log(loggers ...LogName) {
	if len(loggers) == 0 {
		loggers = []LogName{DefaultLog}
	}
	for _, logger := range loggers {
		logF, ok := errManager.loggers[logger]
		if !ok {
			panic(fmt.Sprintf("Logger %s not found", logger))
		}
		logF(e)
	}
}

func (e enhancedError) Is(err error) bool {
	if err == nil {
		return false
	}
	if errn, ok := err.(*enhancedError); ok {
		if errn.TemplateID != "" {
			return e.TemplateID == errn.TemplateID
		} else {
			return e.ErrorID == errn.ErrorID
		}
	}
	return errors.Is(e.error, err)
}

func (e enhancedError) Unwrap() error {
	return e.GetInternalError()
}

type Wrapper string

func (Wrapper) Type() ErrorOptType {
	return "wrapper"
}

func (Wrapper) MapFormatter() map[string]interface{} {
	return make(map[string]interface{})
}

func (Wrapper) Verbosity() int {
	return 0
}

func (e enhancedError) Wrap(msg string) EnhancedError {
	opts := copyOpts(e.Opts)
	wrapper := Wrapper(msg)
	if existingWrapper, ok := e.Opts[wrapper.Type()]; ok {
		existingWrapper := existingWrapper.(Wrapper)
		opts[wrapper.Type()] = wrapper + existingWrapper
	} else {
		opts[wrapper.Type()] = wrapper
	}
	return &enhancedError{
		ErrorID:    uuid.NewString(),
		TemplateID: e.TemplateID,
		error:      e.error,
		Opts:       opts,
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (e enhancedError) GetStackTrace() errors.StackTrace {
	if err, ok := e.error.(stackTracer); ok {
		return err.StackTrace()
	}
	return nil
}

const sizeOfUintPtr = unsafe.Sizeof(uintptr(0))

func uintptrToBytes(u *errors.Frame) []byte {
	return (*[sizeOfUintPtr]byte)(unsafe.Pointer(u))[:]
}

func (e enhancedError) GetStackTraceHash() string {
	st := e.GetStackTrace()
	buffer := bytes.Buffer{}
	for _, f := range st {
		buffer.Write(uintptrToBytes((&f)))
	}
	return fmt.Sprintf("%x", md5.Sum(buffer.Bytes()))
}

func (e enhancedError) GetInternalError() error {
	return e.error
}

func (e enhancedError) GetOpts() map[ErrorOptType]ErrorOpt {
	return copyOpts(e.Opts)
}

func (e enhancedError) GetOpt(opt ErrorOpt) bool {
	value, ok := e.GetOpts()[(opt).Type()]
	if !ok {
		return false
	}
	v := reflect.ValueOf(opt)
	if v.Kind() != reflect.Ptr {
		panic("opt must be a pointer")
	}
	v = v.Elem()
	v.Set(reflect.ValueOf(value))
	return true
}

func (e enhancedError) GetErrorID() string {
	return e.ErrorID
}

func (e enhancedError) Error() string {
	var wrapper Wrapper
	value, ok := e.GetOpts()[wrapper.Type()]
	if !ok {
		return e.GetInternalError().Error()
	}
	return fmt.Sprintf("%s: %s", value, e.GetInternalError())
}
