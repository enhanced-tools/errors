package opts

import "github.com/enhanced-tools/errors"

type Type string

func (Type) Type() errors.ErrorOptType {
	return "type"
}

func (t Type) MapFormatter() map[string]interface{} {
	return map[string]interface{}{
		"error": t,
	}
}

func (t Type) Verbosity() int {
	return 0
}

const (
	ErrNameResources      Type = "resource"
	ErrNameParameter      Type = "parameter"
	ErrNameAuthorization  Type = "authorization"
	ErrNameOutsideService Type = "outsideService"
	ErrNameHeaders        Type = "headers"
	ErrNamePermissions    Type = "permissions"
	ErrNameInternal       Type = "internal"
)
