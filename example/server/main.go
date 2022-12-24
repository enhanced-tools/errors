package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/enhanced-tools/errors"
	"github.com/enhanced-tools/errors/opts"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// display debug messages in server logs
const LogVerbosity = 50

// display only error core values in response
const JSONVerbosity = 0

func handleErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// if status code not provieded with error assume it is internal error
	var statusCode opts.StatusCode = http.StatusInternalServerError
	enhancedErr := errors.Enhance(err)
	if value, ok := enhancedErr.GetOpts()[statusCode.Type()]; ok {
		statusCode = value.(opts.StatusCode)
	}
	// retrievs requestID from context and add it to error
	requestID := r.Context().Value("requestID")
	enhancedErr = enhancedErr.With(opts.RequestID(requestID.(string)))
	w.WriteHeader(int(statusCode))
	// send error to client in JSON format
	w.Write(errors.AsJSON(enhancedErr, JSONVerbosity))
	// log error with diffrent verbosity level
	enhancedErr.Log()
}

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString()
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIntVar(r *http.Request, name string) (int64, error) {
	strValue := chi.URLParam(r, name)
	if strValue == "" {
		// you can create enhanced error in place and return it
		return 0, errors.Newf("missing variable %s", name).With(opts.StatusCode(http.StatusBadRequest))
	}
	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		// or you can enhance existing error and return it
		return 0, errors.Enhance(err).With(
			opts.StatusCode(http.StatusBadRequest),
			opts.Title(fmt.Sprintf("cannot parse %s as int", name)),
			opts.Debug(strValue),
		)
	}
	return value, nil
}

// ErrDivideByZero is an example of predefined error
var ErrDivideByZero = errors.Template().With(
	opts.Title("divide by zero"),
	opts.StatusCode(http.StatusBadRequest),
)

func main() {
	r := chi.NewRouter()
	// Setting up default logger for all errors
	errors.Manager().Setup("stacks.txt")
	errors.Manager().SetDefaultLogger(errors.CustomLogger(
		// log error in LogFMT format
		errors.WithErrorFormatter(errors.LogFMTFormatter),
		// print no stack trace
		errors.WithStackTraceFormatter(errors.NoStackTrace),
		// but include stack trace in file
		errors.WithSaveStack(true),
		// print debug messages
		errors.WithVerbosity(LogVerbosity),
	))
	r.Use(middleware.Logger)
	r.Use(requestID)
	r.Get("/add/{a}/{b}", func(w http.ResponseWriter, r *http.Request) {
		a, err := getIntVar(r, "a")
		if err != nil {
			handleErrorResponse(w, r, err)
			return
		}
		b, err := getIntVar(r, "b")
		if err != nil {
			handleErrorResponse(w, r, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("%d", a+b)))
	})
	r.Get("/divide/{a}/{b}", func(w http.ResponseWriter, r *http.Request) {
		a, err := getIntVar(r, "a")
		if err != nil {
			handleErrorResponse(w, r, err)
			return
		}
		b, err := getIntVar(r, "b")
		if err != nil {
			handleErrorResponse(w, r, err)
			return
		}
		if b == 0 {
			handleErrorResponse(w, r, ErrDivideByZero.FromEmpty())
			return
		}
		w.Write([]byte(fmt.Sprintf("%d", a/b)))
	})
	log.Println("server started on port 3000")
	http.ListenAndServe(":3000", r)
}
