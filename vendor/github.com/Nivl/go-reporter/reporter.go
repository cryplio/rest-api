// Package reporter contains stuctures and interfaces to deal with reporters
package reporter

//go:generate mockgen -destination implementations/mockreporter/reporter.go -package mockreporter github.com/Nivl/go-reporter Reporter

// Reporter describes a struct used to report errors
type Reporter interface {
	SetUser(u *User)
	AddTag(key, value string)
	AddTags(tags map[string]string)
	ReportError(err error)
	ReportErrorAndWait(err error)
}
