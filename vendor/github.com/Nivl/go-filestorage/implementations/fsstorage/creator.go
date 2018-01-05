package fsstorage

import (
	"context"
	"io/ioutil"

	"github.com/Nivl/go-filestorage"
)

// Makes sure Creator is a logger.Creator
var _ filestorage.Creator = (*Creator)(nil)

// NewCreator returns a filestorage creator that will use the provided keys
// to create a new cloudinary driver for each single logger
func NewCreator(defaultBucket string) *Creator {
	return &Creator{
		defaultBucket: defaultBucket,
	}
}

// Creator creates new filestorage
type Creator struct {
	defaultBucket string
}

// New returns a new fs client
func (c *Creator) New() (filestorage.FileStorage, error) {
	fs, err := NewWithDir(c.defaultBucket)
	if err != nil {
		return nil, err
	}
	return fs, fs.SetBucket(c.defaultBucket)
}

// NewWithContext returns a new gc storage client using the provided context as
// default context
func (c *Creator) NewWithContext(ctx context.Context) (filestorage.FileStorage, error) {
	dir, err := ioutil.TempDir("", "storage")
	if err != nil {
		return nil, err
	}

	fs, err := NewWithContext(ctx, dir)
	if err != nil {
		return nil, err
	}
	return fs, fs.SetBucket(c.defaultBucket)
}
