package main

import (
	"log"
	"os"
	"strconv"

	"github.com/enhanced-tools/errors"
	"github.com/enhanced-tools/errors/opts"
)

// Define error templates as global variables
var (
	ErrWrongNumberOfArguments = errors.Template().With(
		opts.Title("Wrong number of arguments"),
	)
	ErrNotANumber = errors.Template().With(
		opts.Title("Not a number"),
	)
)

func main() {
	errors.Manager().Setup("./stacks.txt")
	errors.Manager().SetDefaultLogger(errors.CustomLogger(
		errors.WithErrorFormatter(errors.MultilineFormatter),
		// errors.WithStackTraceFormatter(errors.NoStackTrace),
	))

	if len(os.Args) != 3 {
		// Use error template to create an error with FromEmpty to mark stack trace
		// and With to add options. The error is then logged with default logger.
		ErrWrongNumberOfArguments.FromEmpty().With(
			ArgCount(len(os.Args)),
		).Log()
		return
	}
	v1, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		ErrNotANumber.From(err).With(Argument(os.Args[1])).Log()
		return
	}
	v2, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		ErrNotANumber.From(err).With(Argument(os.Args[2])).Log()
		return
	}
	log.Printf("%.2f + %f = %f", v1, v2, v1+v2)
}
