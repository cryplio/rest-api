package implementations

import (
	"context"
	"io"

	filestorage "github.com/Nivl/go-filestorage"
)

// WriteIfNotExist is an helper to implements FileStorage.WriteIfNotExist
// Since this function is the same for all provider, let's not rewrite it
func WriteIfNotExist(ctx context.Context, s filestorage.FileStorage, src io.Reader, destPath string) (new bool, url string, err error) {
	exists, err := s.ExistsCtx(ctx, destPath)
	if err != nil {
		return false, "", err
	}

	if !exists {
		if err = s.WriteCtx(ctx, src, destPath); err != nil {
			return false, "", err
		}
	}

	url, err = s.URLCtx(ctx, destPath)
	if err != nil {
		return false, "", err
	}

	return !exists, url, nil
}
