package errors

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type ErrorFormatter func(e EnhancedError, verbosityThreshold int, stackTraceFormatter StackTraceFormatter) string

type StackTraceFormatter func(errors.StackTrace) string

type JSONFormatter func(e EnhancedError, verbosityThreshold int) json.RawMessage

func NoStackTrace(stackTrace errors.StackTrace) string {
	return ""
}

func AsJSON(e EnhancedError, verbosityThreshold ...int) json.RawMessage {
	threshold := 100
	if len(verbosityThreshold) > 0 {
		threshold = verbosityThreshold[0]
	}

	outputMap := make(map[string]interface{})
	for _, opt := range e.GetOpts() {
		if opt.Verbosity() <= threshold {
			for key, value := range opt.MapFormatter() {
				outputMap[key] = value
			}
		}
	}
	outputMap["errorCode"] = e.GetStackTraceHash()
	outputMap["errorID"] = e.GetErrorID()
	output, _ := json.Marshal(outputMap)
	return output
}
