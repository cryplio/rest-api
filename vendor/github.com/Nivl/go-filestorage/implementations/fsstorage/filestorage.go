// Package fsstorage is an implementation of filestorage using the filesystem
package fsstorage

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/Nivl/go-filestorage"
	"github.com/Nivl/go-filestorage/implementations"
)

// New returns a new instance of a File System Storage
func New() (*FSStorage, error) {
	dir, err := ioutil.TempDir("", "storage")
	if err != nil {
		return nil, err
	}
	return NewWithContext(context.Background(), dir)
}

// NewWithDir returns a new instance of a File System Storage with
func NewWithDir(path string) (*FSStorage, error) {
	return NewWithContext(context.Background(), path)
}

// NewWithContext returns a new GCStorage instance using a new Google Cloud
// Storage client attached to the provided context
func NewWithContext(ctx context.Context, path string) (*FSStorage, error) {
	return &FSStorage{
		path: path,
		ctx:  ctx,
	}, nil
}

// FSStorage is an implementation of the FileStorage interface for the file system
type FSStorage struct {
	path   string
	bucket string
	ctx    context.Context
}

// ID returns the unique identifier of the storage provider
func (s *FSStorage) ID() string {
	return "file_system"
}

// SetBucket is used to set the bucket
func (s *FSStorage) SetBucket(name string) error {
	s.bucket = name
	return nil
}

// Read fetches a file a returns a reader
// Will use the defaut context
func (s *FSStorage) Read(filepath string) (io.ReadCloser, error) {
	return s.ReadCtx(s.ctx, filepath)
}

// ReadCtx fetches a file a returns a reader
func (s *FSStorage) ReadCtx(ctx context.Context, filepath string) (io.ReadCloser, error) {
	return os.Open(s.fullPath(filepath))
}

// Exists check if a file exists
// Will use the defaut context
func (s *FSStorage) Exists(filepath string) (bool, error) {
	return s.ExistsCtx(s.ctx, filepath)
}

// ExistsCtx check if a file exists
func (s *FSStorage) ExistsCtx(ctx context.Context, filepath string) (bool, error) {
	fi, err := os.Stat(s.fullPath(filepath))
	if err == nil {
		return !fi.IsDir(), err
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Write copy the provided os.File to dest
// Will use the defaut context
func (s *FSStorage) Write(src io.Reader, destPath string) error {
	return s.WriteCtx(s.ctx, src, destPath)
}

// WriteCtx copy the provided os.File to dest
func (s *FSStorage) WriteCtx(ctx context.Context, src io.Reader, destPath string) error {
	fullPath := s.fullPath(destPath)

	// make sure the path exists
	if err := os.MkdirAll(path.Dir(fullPath), os.ModePerm); err != nil {
		return err
	}

	dest, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	// Copy the file
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a file, ignores files that do not exist
// Will use the defaut context
func (s *FSStorage) Delete(filepath string) error {
	return s.DeleteCtx(s.ctx, filepath)
}

// DeleteCtx removes a file, ignores files that do not exist
func (s *FSStorage) DeleteCtx(ctx context.Context, filepath string) error {
	err := os.Remove(s.fullPath(filepath))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// URL returns the URL of the file
// Will use the defaut context
func (s *FSStorage) URL(filepath string) (string, error) {
	return s.URLCtx(s.ctx, filepath)
}

// URLCtx returns the URL of the file
func (s *FSStorage) URLCtx(ctx context.Context, filepath string) (string, error) {
	return s.fullPath(filepath), nil
}

// SetAttributes sets the attributes of the file
// Will use the defaut context
func (s *FSStorage) SetAttributes(filepath string, attrs *filestorage.UpdatableFileAttributes) (*filestorage.FileAttributes, error) {
	return s.SetAttributesCtx(s.ctx, filepath, attrs)
}

// SetAttributesCtx sets the attributes of the file
func (s *FSStorage) SetAttributesCtx(ctx context.Context, filepath string, attrs *filestorage.UpdatableFileAttributes) (*filestorage.FileAttributes, error) {
	return filestorage.NewFileAttributesFromUpdatable(attrs), nil
}

// Attributes returns the attributes of the file
// Will use the defaut context
func (s *FSStorage) Attributes(filepath string) (*filestorage.FileAttributes, error) {
	return s.AttributesCtx(s.ctx, filepath)
}

// AttributesCtx returns the attributes of the file
func (s *FSStorage) AttributesCtx(ctx context.Context, filepath string) (*filestorage.FileAttributes, error) {
	return &filestorage.FileAttributes{}, nil
}

func (s *FSStorage) fullPath(filepath string) string {
	return path.Join(s.path, s.bucket, filepath)
}

// WriteIfNotExist copies the provided io.Reader to dest if the file does
// not already exist
// Returns:
//   - A boolean specifying if the file got uploaded (true) or if already
//     existed (false).
//   - A URL to the uploaded file
//   - An error if something went wrong
// Will use the defaut context
func (s *FSStorage) WriteIfNotExist(src io.Reader, destPath string) (new bool, url string, err error) {
	return s.WriteIfNotExistCtx(s.ctx, src, destPath)
}

// WriteIfNotExistCtx copies the provided io.Reader to dest if the file does
// not already exist
// Returns:
//   - A boolean specifying if the file got uploaded (true) or if already
//     existed (false).
//   - A URL to the uploaded file
//   - An error if something went wrong
func (s *FSStorage) WriteIfNotExistCtx(ctx context.Context, src io.Reader, destPath string) (new bool, url string, err error) {
	return implementations.WriteIfNotExist(ctx, s, src, destPath)
}
