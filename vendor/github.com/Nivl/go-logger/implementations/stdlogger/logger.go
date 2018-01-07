package stdlogger

import (
	"fmt"
	"log"

	"github.com/Nivl/go-logger"
	"github.com/Nivl/go-logger/implementations"
)

// Makes sure the Logger object implements logger.Logger
var _ logger.Logger = (*Logger)(nil)

// Logger represents a logger that just uses the default log package
type Logger struct {
	staticData []string
}

// AddStaticData is used to add static data to the logs.
// static data will be added to all logs
func (bl *Logger) AddStaticData(msg string, args ...interface{}) {
	bl.staticData = append(bl.staticData, fmt.Sprintf(msg, args...))
}

// Errorf is used to log an error with formating
func (bl *Logger) Errorf(msg string, args ...interface{}) {
	bl.Send(implementations.FormatError(bl.staticData, fmt.Sprintf(msg, args...)))
}

// Error is used to log an error
func (bl *Logger) Error(msg string) {
	bl.Send(implementations.FormatError(bl.staticData, msg))
}

// Send is used to log a raw message
func (bl *Logger) Send(msg string) {
	log.Println(msg)
}

// Close closes the connection to Logger
func (bl *Logger) Close() error {
	return nil
}

// New creates an instance of Logger
func New() *Logger {
	return &Logger{}
}
