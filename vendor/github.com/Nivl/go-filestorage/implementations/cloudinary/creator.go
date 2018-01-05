package cloudinary

import (
	"context"

	"github.com/Nivl/go-filestorage"
)

// Makes sure Creator is a logger.Creator
var _ filestorage.Creator = (*Creator)(nil)

// NewCreator returns a filestorage creator that will use the provided keys
// to create a new cloudinary driver for each single logger
func NewCreator(apiKey, secret, defaultBucket string) *Creator {
	return &Creator{
		apiKey:        apiKey,
		secret:        secret,
		defaultBucket: defaultBucket,
	}
}

// Creator creates new filestorage
type Creator struct {
	apiKey        string
	secret        string
	defaultBucket string
}

// New returns a new le client
func (c *Creator) New() (filestorage.FileStorage, error) {
	fs := New(c.apiKey, c.secret)
	return fs, fs.SetBucket(c.defaultBucket)
}

// NewWithContext returns a new gc storage client using the provided context as
// default context instead of the one provided during the creation of the
// Creator
func (c *Creator) NewWithContext(ctx context.Context) (filestorage.FileStorage, error) {
	fs := NewWithContext(ctx, c.apiKey, c.secret)
	return fs, fs.SetBucket(c.defaultBucket)
}
