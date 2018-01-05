package logger

//go:generate mockgen -destination implementations/mocklogger/creator.go -package mocklogger github.com/Nivl/go-logger Creator

// Creator describes a struct used to create loggers
type Creator interface {
	New() (Logger, error)
}
