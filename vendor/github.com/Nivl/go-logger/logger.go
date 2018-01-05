// Package logger contains interfaces to deal with loggers
package logger

//go:generate mockgen -destination implementations/mocklogger/logger.go -package mocklogger github.com/Nivl/go-logger Logger

// Logger is an interface used for all loggers
type Logger interface {
	// AddStaticData is used to add static data to the logs.
	// static data will be added to all logs
	AddStaticData(msg string, args ...interface{})

	// Errorf is used to log an error with formating
	Errorf(msg string, args ...interface{})

	// Error is used to log an error
	Error(msg string)

	// Close closes the connection to BasicLogger
	Close() error
}
