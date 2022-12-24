package opts

import "github.com/enhanced-tools/errors"

type RequestID string

func (RequestID) Type() errors.ErrorOptType {
	return "request_id"
}

func (r RequestID) MapFormatter() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r,
	}
}

func (RequestID) Verbosity() int {
	return 0
}
