package lelogger

import (
	"fmt"
	"log"

	"github.com/Nivl/go-logger"
	"github.com/Nivl/go-logger/implementations"
	"github.com/bsphere/le_go"
)

// Makes sure the Logger object implements logger.Logger
var _ logger.Logger = (*Logger)(nil)

// New creates a new logger using a new tcp connection
func New(token string) (*Logger, error) {
	driver, err := le_go.Connect(token)
	if err != nil {
		return nil, err
	}
	return &Logger{
		driver: driver,
	}, nil
}

// NewWithDriver creates a new logger using an existing driver, the driver
// will not be closed when Close() is called
func NewWithDriver(driver *le_go.Logger) *Logger {
	return &Logger{
		driver:    driver,
		keepAlive: true,
	}
}

// Logger represents a logger that sends data to logentries
type Logger struct {
	driver     *le_go.Logger
	keepAlive  bool
	staticData []string
}

// AddStaticData is used to add static data to the logs.
// static data will be added to all logs
func (le *Logger) AddStaticData(msg string, args ...interface{}) {
	le.staticData = append(le.staticData, fmt.Sprintf(msg, args...))
}

// Errorf is used to log an error with formating
func (le *Logger) Errorf(msg string, args ...interface{}) {
	le.Send(implementations.FormatError(le.staticData, fmt.Sprintf(msg, args...)))
}

// Error is used to log an error
func (le *Logger) Error(msg string) {
	le.Send(implementations.FormatError(le.staticData, msg))
}

// Send is used to log a raw message
func (le *Logger) Send(msg string) {
	go le.driver.Println(msg)
	log.Println(msg)
}

// Close closes the connection to logentries
func (le *Logger) Close() error {
	if le.keepAlive {
		return nil
	}
	return le.driver.Close()
}
