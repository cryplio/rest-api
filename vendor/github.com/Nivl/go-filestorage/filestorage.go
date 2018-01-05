// Package filestorage describes interfaces to deal with filestorage
package filestorage

import (
	"context"
	"io"
)

// FileStorage is a interface to store and retrieve files
//go:generate mockgen -destination implementations/mockfilestorage/filestorage.go -package mockfilestorage github.com/Nivl/go-filestorage FileStorage
type FileStorage interface {
	// ID returns the unique identifier of the storage provider
	ID() string

	// SetBucket sets the bucket that will contains the files
	SetBucket(bucket string) error

	// Read fetches a file a returns a reader
	// Returns os.ErrNotExist if the file does not exists
	// Will use the defaut context
	Read(filepath string) (io.ReadCloser, error)

	// ReadCtx fetches a file a returns a reader
	// Returns os.ErrNotExist if the file does not exists
	ReadCtx(ctx context.Context, filepath string) (io.ReadCloser, error)

	// Write copies the provided io.Reader to dest
	// Will use the defaut context
	Write(src io.Reader, destPath string) error

	// WriteCtx copies the provided io.Reader to dest
	WriteCtx(ctx context.Context, src io.Reader, destPath string) error

	// Delete removes a file, ignores files that do not exist
	// Will use the defaut context
	Delete(filepath string) error

	// DeleteCtx removes a file, ignores files that do not exist
	DeleteCtx(ctx context.Context, filepath string) error

	// URL returns the URL of the file
	// Will use the defaut context
	URL(filepath string) (string, error)

	// URLCtx returns the URL of the file
	URLCtx(ctx context.Context, filepath string) (string, error)

	// SetAttributes sets the attributes of the file
	// Will use the defaut context
	SetAttributes(filepath string, attrs *UpdatableFileAttributes) (*FileAttributes, error)

	// SetAttributesCtx sets the attributes of the file
	SetAttributesCtx(ctx context.Context, filepath string, attrs *UpdatableFileAttributes) (*FileAttributes, error)

	// Attributes returns the attributes of the file
	// Will use the defaut context
	Attributes(filepath string) (*FileAttributes, error)

	// AttributesCtx returns the attributes of the file
	AttributesCtx(ctx context.Context, filepath string) (*FileAttributes, error)

	// Exists check if a file exists
	// Will use the defaut context
	Exists(filepath string) (bool, error)

	// ExistsCtx check if a file exists
	ExistsCtx(ctx context.Context, filepath string) (bool, error)

	// WriteIfNotExist copies the provided io.Reader to dest if the file does
	// not already exist
	// Returns:
	//   - A boolean specifying if the file got uploaded (true) or if already
	//     existed (false).
	//   - A URL to the uploaded file
	//   - An error if something went wrong
	// Will use the defaut context
	WriteIfNotExist(src io.Reader, destPath string) (new bool, url string, err error)

	// WriteIfNotExistCtx copies the provided io.Reader to dest if the file does
	// not already exist
	// Returns:
	//   - A boolean specifying if the file got uploaded (true) or if already
	//     existed (false).
	//   - A URL to the uploaded file
	//   - An error if something went wrong
	WriteIfNotExistCtx(ctx context.Context, src io.Reader, destPath string) (new bool, url string, err error)
}

// FileAttributes represents the attributes a file can have
type FileAttributes struct {
	ContentType        string
	ContentLanguage    string
	ContentEncoding    string
	ContentDisposition string
	CacheControl       string
	Metadata           map[string]string
}

// NewFileAttributesFromUpdatable returns a FileAttributes from a UpdatableFileAttributes
func NewFileAttributesFromUpdatable(attrs *UpdatableFileAttributes) *FileAttributes {
	fa := &FileAttributes{
		Metadata: attrs.Metadata,
	}
	if val, ok := attrs.ContentType.(string); ok {
		fa.ContentType = val
	}
	if val, ok := attrs.ContentDisposition.(string); ok {
		fa.ContentDisposition = val
	}
	if val, ok := attrs.ContentLanguage.(string); ok {
		fa.ContentLanguage = val
	}
	if val, ok := attrs.ContentEncoding.(string); ok {
		fa.ContentEncoding = val
	}
	if val, ok := attrs.CacheControl.(string); ok {
		fa.CacheControl = val
	}
	return fa
}

// UpdatableFileAttributes represents the updatable attributes a file can have
type UpdatableFileAttributes struct {
	ContentType        interface{}
	ContentLanguage    interface{}
	ContentEncoding    interface{}
	ContentDisposition interface{}
	CacheControl       interface{}
	Metadata           map[string]string // set to map[string]string{} to delete
}
