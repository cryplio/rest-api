// Package gcstorage is an implementation of filestorage for Google Cloud
package gcstorage

import (
	"context"
	"io"

	"github.com/Nivl/go-filestorage"
	"github.com/Nivl/go-filestorage/implementations"
	"google.golang.org/api/option"

	"fmt"

	"cloud.google.com/go/storage"
)

var _ filestorage.FileStorage = (*GCStorage)(nil)

// New returns a new GCStorage instance using a new Google Cloud Storage client
func New(apiKey string) (*GCStorage, error) {
	return NewWithContext(context.Background(), apiKey)
}

// NewWithContext returns a new GCStorage instance using a new Google Cloud
// Storage client attached to the provided context
func NewWithContext(ctx context.Context, apiKey string) (*GCStorage, error) {
	client, err := storage.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GCStorage{
		ctx:    ctx,
		client: client,
	}, nil
}

// NewWithClient returns a new instance of a Google Cloud Storage using the
// provided client
func NewWithClient(defaultContext context.Context, client *storage.Client) *GCStorage {
	return &GCStorage{
		ctx:    defaultContext,
		client: client,
	}
}

// GCStorage is an implementation of the FileStorage interface for Google Cloud
type GCStorage struct {
	ctx        context.Context
	client     *storage.Client
	bucket     *storage.BucketHandle
	bucketName string
}

// ID returns the unique identifier of the storage provider
func (s *GCStorage) ID() string {
	return "google_cloud_storage"
}

// SetBucket is used to set the bucket
// Always return nil
func (s *GCStorage) SetBucket(name string) error {
	s.bucketName = name
	s.bucket = s.client.Bucket(name)
	return nil
}

// Read fetches a file a returns a reader
// Will use the defaut context
func (s *GCStorage) Read(filepath string) (io.ReadCloser, error) {
	return s.ReadCtx(s.ctx, filepath)
}

// ReadCtx fetches a file a returns a reader
func (s *GCStorage) ReadCtx(ctx context.Context, filepath string) (io.ReadCloser, error) {
	return s.bucket.Object(filepath).NewReader(ctx)
}

// Write copy the provided os.File to dest
// Will use the defaut context
func (s *GCStorage) Write(src io.Reader, destPath string) error {
	return s.WriteCtx(s.ctx, src, destPath)
}

// WriteCtx copy the provided os.File to dest
func (s *GCStorage) WriteCtx(ctx context.Context, src io.Reader, destPath string) error {
	obj := s.bucket.Object(destPath)
	dest := obj.NewWriter(ctx)

	// Copy the file
	_, err := io.Copy(dest, src)
	if err != nil {
		dest.CloseWithError(err)
		return err
	}

	// Send the changes to GCP
	return dest.Close()
}

// SetAttributes sets the attributes of the file
// Will use the defaut context
func (s *GCStorage) SetAttributes(filepath string, attrs *filestorage.UpdatableFileAttributes) (*filestorage.FileAttributes, error) {
	return s.SetAttributesCtx(s.ctx, filepath, attrs)
}

// SetAttributesCtx sets the attributes of the file
func (s *GCStorage) SetAttributesCtx(ctx context.Context, filepath string, attrs *filestorage.UpdatableFileAttributes) (*filestorage.FileAttributes, error) {
	gcsAttrs, err := s.bucket.Object(filepath).Update(ctx, *updatableAttrsToGCStorage(attrs))
	if err != nil {
		return nil, err
	}
	return newAttributes(gcsAttrs), nil
}

// Attributes returns the attributes of the file
// Will use the defaut context
func (s *GCStorage) Attributes(filepath string) (*filestorage.FileAttributes, error) {
	return s.AttributesCtx(s.ctx, filepath)
}

// AttributesCtx returns the attributes of the file
func (s *GCStorage) AttributesCtx(ctx context.Context, filepath string) (*filestorage.FileAttributes, error) {
	gcsAttrs, err := s.bucket.Object(filepath).Attrs(ctx)
	if err != nil {
		return nil, err
	}
	return newAttributes(gcsAttrs), nil
}

// Exists check if a file exists
// Will use the defaut context
func (s *GCStorage) Exists(filepath string) (bool, error) {
	return s.ExistsCtx(s.ctx, filepath)
}

// ExistsCtx check if a file exists
func (s *GCStorage) ExistsCtx(ctx context.Context, filepath string) (bool, error) {
	_, err := s.AttributesCtx(ctx, filepath)
	if err == nil {
		return true, nil
	}
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	return false, err
}

// URL returns the URL of the file
// Will use the defaut context
func (s *GCStorage) URL(filepath string) (string, error) {
	return s.URLCtx(s.ctx, filepath)
}

// URLCtx returns the URL of the file
func (s *GCStorage) URLCtx(ctx context.Context, filepath string) (string, error) {
	return fmt.Sprintf("https://%s.storage.googleapis.com/%s", s.bucketName, filepath), nil
}

// Delete removes a file, ignores files that do not exist
// Will use the defaut context
func (s *GCStorage) Delete(filepath string) error {
	return s.DeleteCtx(s.ctx, filepath)
}

// DeleteCtx removes a file, ignores files that do not exist
func (s *GCStorage) DeleteCtx(ctx context.Context, filepath string) error {
	return s.bucket.Object(filepath).Delete(ctx)
}

// WriteIfNotExist copies the provided io.Reader to dest if the file does
// not already exist
// Returns:
//   - A boolean specifying if the file got uploaded (true) or if already
//     existed (false).
//   - A URL to the uploaded file
//   - An error if something went wrong
// Will use the defaut context
func (s *GCStorage) WriteIfNotExist(src io.Reader, destPath string) (new bool, url string, err error) {
	return s.WriteIfNotExistCtx(s.ctx, src, destPath)
}

// WriteIfNotExistCtx copies the provided io.Reader to dest if the file does
// not already exist
// Returns:
//   - A boolean specifying if the file got uploaded (true) or if already
//     existed (false).
//   - A URL to the uploaded file
//   - An error if something went wrong
func (s *GCStorage) WriteIfNotExistCtx(ctx context.Context, src io.Reader, destPath string) (new bool, url string, err error) {
	return implementations.WriteIfNotExist(ctx, s, src, destPath)
}

// updatableAttrsToGCStorage converts a *UpdatableFileAttributes into a *storage.ObjectAttrsToUpdate
func updatableAttrsToGCStorage(attrs *filestorage.UpdatableFileAttributes) *storage.ObjectAttrsToUpdate {
	return &storage.ObjectAttrsToUpdate{
		ContentType:        attrs.ContentType,
		ContentDisposition: attrs.ContentDisposition,
		ContentLanguage:    attrs.ContentLanguage,
		ContentEncoding:    attrs.ContentEncoding,
		CacheControl:       attrs.CacheControl,
		Metadata:           attrs.Metadata,
	}
}

// newAttributes converts a *storage.ObjectAttrs to *FileAttributes
func newAttributes(attrs *storage.ObjectAttrs) *filestorage.FileAttributes {
	return &filestorage.FileAttributes{
		ContentType:        attrs.ContentType,
		ContentDisposition: attrs.ContentDisposition,
		ContentLanguage:    attrs.ContentLanguage,
		ContentEncoding:    attrs.ContentEncoding,
		CacheControl:       attrs.CacheControl,
		Metadata:           attrs.Metadata,
	}
}
