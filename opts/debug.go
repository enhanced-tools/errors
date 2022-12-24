package opts

import "github.com/enhanced-tools/errors"

type debug struct {
	Value interface{}
}

func (debug) Type() errors.ErrorOptType {
	return "debug"
}

func (d debug) MapFormatter() map[string]interface{} {
	return map[string]interface{}{
		"debug": d,
	}
}

func (d debug) Verbosity() int {
	return 50
}

func Debug(value interface{}) debug {
	return debug{Value: value}
}
