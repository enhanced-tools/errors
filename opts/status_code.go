package opts

import "github.com/enhanced-tools/errors"

type StatusCode int64

func (StatusCode) Type() errors.ErrorOptType {
	return "status_code"
}

func (s StatusCode) MapFormatter() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": s,
	}
}

func (s StatusCode) Verbosity() int {
	return 0
}
