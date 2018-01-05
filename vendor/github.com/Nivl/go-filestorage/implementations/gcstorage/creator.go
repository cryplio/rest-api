package gcstorage

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/Nivl/go-filestorage"
	"google.golang.org/api/option"
)

// Makes sure Creator is a logger.Creator
var _ filestorage.Creator = (*Creator)(nil)

// NewCreator returns a filestorage creator that will use the provided key
// to create a new google storage client that will be reused for every
// new gcstorage instance
func NewCreator(apiKey string, defaultBucket string) (*Creator, error) {
	return NewCreatorWithContext(context.Background(), apiKey, defaultBucket)
}

// NewCreatorWithContext returns a filestorage creator that will
// use the provided context and key to create a new google storage
// client that will be reused for every new gcstorage instance
// The provided context will be used as default context for all new FileStorage
// instance
func NewCreatorWithContext(ctx context.Context, apiKey string, defaultBucket string) (*Creator, error) {
	client, err := storage.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &Creator{
		client:        client,
		defaultCtx:    ctx,
		defaultBucket: defaultBucket,
	}, nil
}

// Creator creates new filestorage
type Creator struct {
	client        *storage.Client
	defaultCtx    context.Context
	defaultBucket string
}

// New returns a new le client
func (c *Creator) New() (filestorage.FileStorage, error) {
	fs := NewWithClient(c.defaultCtx, c.client)
	err := fs.SetBucket(c.defaultBucket)
	return fs, err
}

// NewWithContext returns a new gc storage client using the provided context as
// default context instead of the one provided during the creation of the
// Creator
func (c *Creator) NewWithContext(ctx context.Context) (filestorage.FileStorage, error) {
	fs := NewWithClient(ctx, c.client)
	err := fs.SetBucket(c.defaultBucket)
	return fs, err
}
