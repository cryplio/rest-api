package dependencies

import (
	"context"
	"errors"
	"sync"

	filestorage "github.com/Nivl/go-filestorage"
	logger "github.com/Nivl/go-logger"
	mailer "github.com/Nivl/go-mailer"
	reporter "github.com/Nivl/go-reporter"
	db "github.com/Nivl/go-sqldb"
)

var _ Dependencies = (*AppDependencies)(nil)

// AppDependencies is a structure implementing the Dependencies interface
type AppDependencies struct {
	sync.Mutex

	defaultSQLCon db.Connection
	defaultLogger logger.Logger
	defaultMailer mailer.Mailer

	reporterCreator    reporter.Creator
	loggerCreator      logger.Creator
	filestorageCreator filestorage.Creator
}

// SetDB creates a connection to a SQL database
func (deps *AppDependencies) SetDB(con db.Connection) {
	deps.defaultSQLCon = con
}

// DB returns the current SQL connection
func (deps *AppDependencies) DB() db.Connection {
	return deps.defaultSQLCon
}

// SetLoggerCreator sets a logger creator used to generate new loggers
func (deps *AppDependencies) SetLoggerCreator(creator logger.Creator) {
	deps.loggerCreator = creator
}

// DefaultLogger returns the default logger to use app wide
func (deps *AppDependencies) DefaultLogger() (logger.Logger, error) {
	deps.Lock()
	defer deps.Unlock()

	if deps.loggerCreator == nil {
		return nil, errors.New("no logger creator has been set")
	}
	if deps.defaultLogger == nil {
		var err error
		deps.defaultLogger, err = deps.loggerCreator.New()
		return deps.defaultLogger, err
	}
	return deps.defaultLogger, nil
}

// NewLogger returns a new logger using the logger Creator
func (deps *AppDependencies) NewLogger() (logger.Logger, error) {
	if deps.loggerCreator == nil {
		return nil, errors.New("no logger creator has been set")
	}
	return deps.loggerCreator.New()
}

// SetMailer sets the mailer to be used to send emails
func (deps *AppDependencies) SetMailer(m mailer.Mailer) {
	deps.defaultMailer = m
}

// Mailer returns the mailer set with SetMailer
func (deps *AppDependencies) Mailer() (mailer.Mailer, error) {
	if deps.defaultMailer == nil {
		return nil, errors.New("no mailer has been set")
	}
	return deps.defaultMailer, nil
}

// SetReporterCreator sets a reporter creator used to generate new reporters
func (deps *AppDependencies) SetReporterCreator(creator reporter.Creator) {
	deps.reporterCreator = creator
}

// NewReporter creates a new reporter using the provided reporter Creator
func (deps *AppDependencies) NewReporter() (reporter.Reporter, error) {
	if deps.reporterCreator == nil {
		return nil, errors.New("no reporter creator has been set")
	}
	return deps.reporterCreator.New()
}

// SetFileStorageCreator returns the default filestorage following this order
func (deps *AppDependencies) SetFileStorageCreator(creator filestorage.Creator) {
	deps.filestorageCreator = creator
}

// NewFileStorage creates a new filestorage using the provided reporter Creator
func (deps *AppDependencies) NewFileStorage(ctx context.Context) (filestorage.FileStorage, error) {
	if deps.reporterCreator == nil {
		return nil, errors.New("no filestorage creator has been set")
	}
	return deps.filestorageCreator.NewWithContext(ctx)
}
