package errors

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logfmt/logfmt"
	pkgerrors "github.com/pkg/errors"
)

func LogFMTFormatter(e EnhancedError, verbosityThreshold int, stackTraceFormatter StackTraceFormatter) string {
	var sb strings.Builder
	encoder := logfmt.NewEncoder(&sb)
	encoder.EncodeKeyval("errorID", e.GetErrorID())
	encoder.EncodeKeyval("errorCode", e.GetStackTraceHash())
	var wrapper Wrapper
	value, ok := e.GetOpts()[wrapper.Type()]
	if ok {
		wrapper = value.(Wrapper)
	}
	if wrapper != "" {
		wrapper = Wrapper(fmt.Sprintf("%s: ", wrapper))
	}
	encoder.EncodeKeyval("content", fmt.Sprintf("%s%s", wrapper, e.GetInternalError()))
	opts := make(map[string]interface{})
	for _, opt := range e.GetOpts() {
		for key, value := range opt.MapFormatter() {
			if opt.Verbosity() > verbosityThreshold {
				continue
			}
			opts[key] = value
		}
	}
	for opt, value := range opts {
		if reflect.ValueOf(value).Kind() == reflect.Struct {
			valueBytes, err := json.Marshal(value)
			if err != nil {
				panic(err)
			}
			value = string(valueBytes)
		}
		encoder.EncodeKeyval(opt, value)
	}
	stackTrace := e.GetStackTrace()
	stackTraceMsg := stackTraceFormatter(stackTrace)
	stackTraceMsg = strings.ReplaceAll(stackTraceMsg, "\n", "$")

	if stackTraceMsg != "" {
		encoder.EncodeKeyval("stackTrace", stackTraceMsg)
	}

	if err := encoder.EndRecord(); err != nil {
		return "Error in logfmt formatter"
	}
	return sb.String()
}

func JSONStackTraceFormatter(st pkgerrors.StackTrace) string {
	type stackTraceLine struct {
		SourceFile   string `json:"file"`
		SourceLine   string `json:"line"`
		FunctionName string `json:"func"`
	}
	stackTrace := make([]stackTraceLine, 0, len(st))
	for _, f := range st {
		stackTrace = append(stackTrace, stackTraceLine{
			SourceFile:   fmt.Sprintf("%s", f),
			SourceLine:   fmt.Sprintf("%d", f),
			FunctionName: fmt.Sprintf("%n", f),
		})
	}
	msg, err := json.Marshal(stackTrace)
	if err != nil {
		return "Error with marshaling stack trace"
	}
	return string(msg)
}
