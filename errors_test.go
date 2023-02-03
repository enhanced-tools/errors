package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/enhanced-tools/errors"
	"github.com/enhanced-tools/errors/opts"
	"github.com/stretchr/testify/assert"
)

func TestNewf(t *testing.T) {
	err := errors.Newf("%s %d", "t", 1)
	msg := err.Error()
	assert.Equal(t, "t 1", msg)
}

func TestGetOption(t *testing.T) {
	var title opts.Title
	err := errors.Newf("%s %d", "t", 1).With(opts.Title("title"))
	err.GetOpt(&title)
	assert.Equal(t, opts.Title("title"), title)
}

func TestTemplateIs2Enhanced(t *testing.T) {
	err := errors.Template()
	err1 := err.FromEmpty()
	err2 := err.From(fmt.Errorf("error"))
	assert.True(t, err1.Is(err2))
	assert.True(t, err2.Is(err1))
	assert.True(t, stderrors.Is(err1, err))
	assert.True(t, stderrors.Is(err2, err))
}

func TestTemplate1Enhance1Common(t *testing.T) {
	err := errors.Template()
	errCommon := fmt.Errorf("error")
	err1 := err.From(errCommon)
	assert.True(t, err1.Is(errCommon))
	assert.True(t, stderrors.Is(err1, errCommon), "standard package errors.Is should be true for wrapped error")
	assert.False(t, stderrors.Is(err1, fmt.Errorf("other error")), "standard package errors.Is should be false for other new error")
}

func TestUnwrapReturnsInternalError(t *testing.T) {
	err := fmt.Errorf("Orginal error")
	enhanced := errors.Enhance(err)
	// Additional Unwrap is for *withStack wrapper form "github.com/pkg/errors"
	assert.Equal(t, err, stderrors.Unwrap(enhanced.GetInternalError()), "Enhanced error GetInternalError() should return original err")
	assert.Equal(t, err, stderrors.Unwrap(stderrors.Unwrap(enhanced)), "Enhanced error should return original on errors.Unwrap")
}

func TestMultilineError(t *testing.T) {
	err := fmt.Errorf("Orginal error")
	enhanced := errors.Enhance(err)
	enhanced.Log()
}

func TestLogFMTError(t *testing.T) {
	err := fmt.Errorf("Orginal error")
	errors.Manager().Setup("./stacks.txt")
	errors.Manager().SetDefaultLogger(errors.CustomLogger(
		errors.WithErrorFormatter(errors.LogFMTFormatter),
		errors.WithVerbosity(200),
		errors.WithStackTraceFormatter(errors.NoStackTrace),
		errors.WithSaveStack(true),
	))
	enhanced := errors.Enhance(err).With(opts.Debug("John"), opts.Title("Smoth"))
	enhanced.Log()
}
