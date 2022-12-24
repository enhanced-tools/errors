package main

import "github.com/enhanced-tools/errors"

// Implements ErrorOpt interface
type Argument string

func (a Argument) Type() errors.ErrorOptType {
	return "argument"
}

func (a Argument) MapFormatter() map[string]interface{} {
	return map[string]interface{}{
		"argument": a,
	}
}

func (Argument) Verbosity() int {
	return 0
}

// Implements ErrorOpt interface
type ArgCount int

func (a ArgCount) Type() errors.ErrorOptType {
	return "argCount"
}

func (a ArgCount) MapFormatter() map[string]interface{} {
	return map[string]interface{}{
		"argCount": a,
	}
}

func (ArgCount) Verbosity() int {
	return 0
}
