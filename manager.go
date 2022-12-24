package errors

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type ErrorsManager interface {
	// SaveStack saves the stack trace of the error in the stack trace file. Format is optional and can be used to format the stack trace
	SaveStack(err EnhancedError, format ...StackTraceFormatter) error
	// RegisterLogger registers a logger for a specific log name
	RegisterLogger(name LogName, logger LoggerFunc)
	// SetDefaultLogger sets the default logger
	SetDefaultLogger(logger LoggerFunc)
	// Setup creates a new stack trace file and reads the existing one
	Setup(stackTracePath string) error
}

type errorsManager struct {
	stackTracePath string
	stackFile      *os.File
	stackWriter    *Writer

	stacks map[string]bool

	loggers map[LogName]LoggerFunc
}

var errManager errorsManager

func init() {
	errManager.stacks = make(map[string]bool)
	errManager.loggers = map[LogName]LoggerFunc{
		DefaultLog: DefaultLogger(),
	}
}

func (e *errorsManager) SaveStack(err EnhancedError, format ...StackTraceFormatter) error {
	if e.stackTracePath == "" {
		return fmt.Errorf("stack trace path not set")
	}
	formatter := MultilineStackTraceFormatter
	if len(format) > 0 {
		formatter = format[0]
	}
	stackTraceHash := err.GetStackTraceHash()
	if _, ok := e.stacks[stackTraceHash]; !ok {
		e.stacks[stackTraceHash] = true
		message := fmt.Sprintf(">>> %s\n%s\n", stackTraceHash, formatter(err.GetStackTrace()))
		if _, err := e.stackWriter.WriteString(message); err != nil {
			return err
		}
	}
	return nil
}

func Manager() ErrorsManager {
	return &errManager
}

func (m *errorsManager) Setup(stackTracePath string) error {
	if m.stackTracePath != "" {
		return fmt.Errorf("duplicated error initialization")
	}
	err := func() error {
		if _, err := os.Stat(stackTracePath); os.IsNotExist(err) {
			return nil
		}
		file, err := os.Open(stackTracePath)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if text := scanner.Text(); strings.HasPrefix(text, ">>>") {
				parts := strings.Split(text, " ")
				errManager.stacks[parts[1]] = true
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return errors.Wrap(err, "reading error stack file")
	}
	stFile, err := os.OpenFile(stackTracePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.Wrap(err, "Log Setup")
	}

	m.stackTracePath = stackTracePath
	m.stackFile = stFile
	m.stackWriter = NewWriter(stFile)

	return nil
}

type LoggerFunc func(err EnhancedError)

type LogName string

const DefaultLog LogName = "default"

func (m *errorsManager) RegisterLogger(name LogName, logger LoggerFunc) {
	m.loggers[name] = logger
}

func (m *errorsManager) SetDefaultLogger(logger LoggerFunc) {
	m.loggers[DefaultLog] = logger
}

type Writer struct {
	w  io.Writer
	mu sync.Mutex
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	n, err := w.w.Write(p)
	w.mu.Unlock()

	return n, err
}

func (w *Writer) WriteString(str string) (int, error) {
	w.mu.Lock()
	n, err := w.w.Write([]byte(str))
	w.mu.Unlock()

	return n, err
}
