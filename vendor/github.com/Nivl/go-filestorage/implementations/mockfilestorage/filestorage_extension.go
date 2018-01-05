package mockfilestorage

import (
	"errors"
	"os"
	"path"

	"github.com/Nivl/go-filestorage"
	gomock "github.com/golang/mock/gomock"
)

var (
	// StringType represent a string argument
	StringType = gomock.Eq("string")

	// UpdatableAttrType represents a *filestorage.UpdatableFileAttributes argument
	UpdatableAttrType = gomock.Eq("*filestorage.UpdatableFileAttributes")

	// AnyType represents am argument that can accept anything
	AnyType = gomock.Any()
)

// WriteIfNotExistSuccess is an helper that expects a WriteIfNotExist to
// succeed with the provided params
func (mr *MockFileStorageMockRecorder) WriteIfNotExistSuccess(isNew bool, url string) *gomock.Call {
	return mr.WriteIfNotExist(AnyType, StringType).Return(isNew, url, nil).Times(1)
}

// WriteIfNotExistError is an helper that expects a WriteIfNotExist to
// fail
func (mr *MockFileStorageMockRecorder) WriteIfNotExistError() *gomock.Call {
	call := mr.WriteIfNotExist(AnyType, StringType)
	call.Return(false, "", errors.New("server unreachable"))
	return call.Times(1)
}

// ReadSuccess is an helper that expects a Read that succeed
func (mr *MockFileStorageMockRecorder) ReadSuccess(cwd, filename string) *gomock.Call {
	filePath := path.Join(cwd, "fixtures", filename)
	return mr.Read(StringType).Return(os.Open(filePath)).Times(1)
}

// ExistsSuccess is an helper that expects Exists() to return true
func (mr *MockFileStorageMockRecorder) ExistsSuccess() *gomock.Call {
	return mr.Exists(StringType).Return(true, nil).Times(1)
}

// NotExistsSuccess is an helper that expects Exists() to return false
func (mr *MockFileStorageMockRecorder) NotExistsSuccess() *gomock.Call {
	return mr.Exists(StringType).Return(false, nil).Times(1)
}

// URLSuccess is an helper that expects URL() to return given param
func (mr *MockFileStorageMockRecorder) URLSuccess(url string) *gomock.Call {
	return mr.URL(StringType).Return(url, nil).Times(1)
}

// SetAttributesSuccess is an helper that expects SetAttributes() to work,
// and to return an empty content
func (mr *MockFileStorageMockRecorder) SetAttributesSuccess() *gomock.Call {
	attrs := &filestorage.FileAttributes{}
	return mr.SetAttributes(StringType, UpdatableAttrType).Return(attrs, nil).Times(1)
}

// SetAttributesRetSuccess is an helper that expects SetAttributes() to
// return the provided object
func (mr *MockFileStorageMockRecorder) SetAttributesRetSuccess(attrs *filestorage.FileAttributes) *gomock.Call {
	return mr.SetAttributes(StringType, UpdatableAttrType).Return(attrs, nil).Times(1)
}

// AttributesSuccess is an helper that expects Attributes() to
// return the provided object
func (mr *MockFileStorageMockRecorder) AttributesSuccess(attrs *filestorage.FileAttributes) *gomock.Call {
	return mr.Attributes(StringType).Return(attrs, nil).Times(1)
}

// DeleteSuccess is an helper that expects Delete() to succeed
func (mr *MockFileStorageMockRecorder) DeleteSuccess() *gomock.Call {
	return mr.Delete(StringType).Return(nil).Times(1)
}
