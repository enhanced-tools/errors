package errors

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
	pkgerrors "github.com/pkg/errors"
)

func MultilineFormatter(e EnhancedError, verbosityThreshold int, stackTraceFormatter StackTraceFormatter) string {
	var sb strings.Builder
	stackTraceHash := e.GetStackTraceHash()
	sb.WriteString(fmt.Sprintf("--- %s --- %s --- %s \n", aurora.Red("ERROR"), aurora.Blue(stackTraceHash), e.GetErrorID()))
	var wrapper Wrapper
	value, ok := e.GetOpts()[wrapper.Type()]
	if ok {
		wrapper = value.(Wrapper)
	}
	if wrapper != "" {
		wrapper = Wrapper(fmt.Sprintf("%s: ", wrapper))
	}
	sb.WriteString(fmt.Sprintf("\tCONTENT: %s%s \n", wrapper, e.GetInternalError()))

	opts := make(map[string]interface{})
	for _, opt := range e.GetOpts() {
		for key, value := range opt.MapFormatter() {
			if opt.Verbosity() > verbosityThreshold {
				continue
			}
			opts[key] = value
		}
	}
	optBytes, err := json.MarshalIndent(opts, "\t", "  ")
	if err != nil {
		return "Error in Marshaling Error"
	}
	if len(opts) > 0 {
		sb.WriteString("\t")
		sb.Write(optBytes)
		sb.WriteString("\n")
	}
	stackTrace := e.GetStackTrace()
	msg := stackTraceFormatter(stackTrace)
	if msg != "" {
		sb.WriteString(fmt.Sprintf("\tSTACK TRACE: \n%s", msg))
	}
	return sb.String()
}

func MultilineStackTraceFormatter(st pkgerrors.StackTrace) string {
	var sb strings.Builder
	for _, f := range st {
		sb.WriteString(fmt.Sprintf("\t%+s:%d\n", f, f))
	}
	return sb.String()
}
