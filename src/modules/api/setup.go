package api

import (
	"github.com/Nivl/go-filestorage/implementations/cloudinary"
	"github.com/Nivl/go-filestorage/implementations/fsstorage"
	"github.com/Nivl/go-filestorage/implementations/gcstorage"
	"github.com/Nivl/go-logger/implementations/lelogger"
	mailer "github.com/Nivl/go-mailer"
	"github.com/Nivl/go-mailer/implementations/printmailer"
	"github.com/Nivl/go-mailer/implementations/sendgridmailer"
	reporter "github.com/Nivl/go-reporter"
	"github.com/Nivl/go-reporter/implementations/mailerreporter"
	"github.com/Nivl/go-reporter/implementations/noopreporter"
	"github.com/Nivl/go-reporter/implementations/sentryreporter"
	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-sqldb/implementations/sqlxdb"
	"github.com/kelseyhightower/envconfig"
)

// Args represents the app args
type Args struct {
	Port                   string `default:"5000"`
	PostgresURI            string `required:"true" envconfig:"postgres_uri"`
	LogEntriesToken        string `envconfig:"logentries_token"`
	SendgridAPIKey         string `envconfig:"sendgrid_api_key"`
	EmailFrom              string `envconfig:"email_default_from"`
	EmailTo                string `envconfig:"email_default_to"`
	SendgridStacktraceUUID string `envconfig:"sendgrid_stacktrace_uuid"`
	CloudinaryAPIKey       string `envconfig:"cloudinary_api_key"`
	CloudinarySecret       string `envconfig:"cloudinary_secret"`
	CloudinaryBucket       string `envconfig:"cloudinary_bucket"`
	LocalFSBucket          string `envconfig:"local_fs_bucket"`
	GCPAPIKey              string `envconfig:"gcp_api_key"`
	GCPProject             string `envconfig:"gcp_project"`
	GCPBucket              string `envconfig:"gcp_bucket"`
	Debug                  bool   `default:"false"`
	SentryDSN              string `envconfig:"sentry_dsn"`
}

// DefaultSetup parses the env and returns the args and dependencies
func DefaultSetup() (*Args, dependencies.Dependencies, error) {
	params := &Args{}
	if err := envconfig.Process("", params); err != nil {
		return nil, nil, err
	}

	deps := &dependencies.AppDependencies{}
	err := Setup(params, deps)
	return params, deps, err
}

// Setup parses the env, sets the app globals and returns the params
func Setup(params *Args, deps dependencies.Dependencies) error {
	sqlClient, err := sqlxdb.New(params.PostgresURI)
	if err != nil {
		return err
	}
	deps.SetDB(sqlClient)

	if params.LogEntriesToken != "" {
		creator, err := lelogger.NewSharedCreator(params.LogEntriesToken)
		if err != nil {
			return err
		}
		deps.SetLoggerCreator(creator)
	}

	var m mailer.Mailer = &printmailer.Mailer{}
	if params.SendgridAPIKey != "" {
		m = sendgridmailer.New(params.SendgridAPIKey, params.EmailFrom, params.EmailTo, params.SendgridStacktraceUUID)
	}
	deps.SetMailer(m)

	if params.GCPAPIKey != "" {
		s, err := gcstorage.NewCreator(params.GCPAPIKey, params.GCPBucket)
		if err != nil {
			return err
		}
		deps.SetFileStorageCreator(s)
	} else if params.CloudinaryAPIKey != "" {
		deps.SetFileStorageCreator(cloudinary.NewCreator(params.CloudinaryAPIKey, params.CloudinarySecret, params.CloudinaryBucket))
	} else {
		deps.SetFileStorageCreator(fsstorage.NewCreator(params.LocalFSBucket))
	}

	var creator reporter.Creator
	creator, _ = noopreporter.NewCreator()
	if params.SentryDSN != "" {
		creator, err = sentryreporter.NewCreator(params.SentryDSN)
		if err != nil {
			return err
		}
	} else if true {
		creator, _ = mailerreporter.NewCreator(m)
	}
	deps.SetReporterCreator(creator)

	return nil
}
