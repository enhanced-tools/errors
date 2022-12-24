package opts

import "github.com/enhanced-tools/errors"

type Title string

func (Title) Type() errors.ErrorOptType {
	return "title"
}

func (t Title) MapFormatter() map[string]interface{} {
	return map[string]interface{}{
		"message": t,
	}
}

func (Title) Verbosity() int {
	return 0
}
