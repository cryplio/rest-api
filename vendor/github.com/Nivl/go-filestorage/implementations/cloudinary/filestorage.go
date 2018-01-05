package cloudinary

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Nivl/go-filestorage"
	"github.com/Nivl/go-filestorage/implementations"
)

var fileTypes = []string{"image", "raw", "video"}

// cdnrErrorResponse represents the response returned by Cloudinary upon error
type cdnrErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// cdnrDeleteResponse represents the response returned by cloudinary when
// deleting a file
type cdnrDeleteResponse struct {
	Result string `json:"result"`
}

// cdnrUploadSuccessResponse represents the response returned by Cloudinary
// upon successful upload
type cdnrUploadSuccessResponse struct {
	SecureURL string `json:"secure_url"`
}

// New returns a new instance of a Cloudinary Storage
func New(apiKey, secret string) *Cloudinary {
	return NewWithContext(context.Background(), apiKey, secret)
}

// NewWithContext returns a new instance of a Cloudinary Storage
// attached to the provided context
func NewWithContext(ctx context.Context, apiKey, secret string) *Cloudinary {
	return &Cloudinary{
		apiKey: apiKey,
		secret: secret,
		ctx:    ctx,
		cache:  map[string]string{},
	}
}

// signature generates and returns a request signature as well as
// the associated timestamp
func (s *Cloudinary) signature(publicID string, invalidate bool) (string, string) {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	signature := "public_id=" + publicID + "&timestamp=" + timestamp + s.secret
	if invalidate {
		signature = "invalidate=true&" + signature
	}
	hash := sha1.Sum([]byte(signature))
	return fmt.Sprintf("%x", hash), timestamp
}

// Cloudinary is an implementation of the FileStorage interface for Cloudinary
type Cloudinary struct {
	cache     map[string]string
	apiKey    string
	secret    string
	ctx       context.Context
	cloudName string // bucket
}

var _ filestorage.FileStorage = (*Cloudinary)(nil)

// ID returns the unique identifier of the storage provider
func (s *Cloudinary) ID() string {
	return "cloudinary"
}

// SetBucket is used to set the bucket
// Always return nil
func (s *Cloudinary) SetBucket(name string) error {
	s.cloudName = name
	return nil
}

// apiBaseURL returns the base URL for an API call
func (s *Cloudinary) apiBaseURL(typ string) string {
	return fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s", s.cloudName, typ)
}

// resBaseURL returns the base URL for a resource
func (s *Cloudinary) resBaseURL(typ string) string {
	return fmt.Sprintf("https://res.cloudinary.com/%s/%s", s.cloudName, typ)
}

// URL returns the URL of the file
// Will use the defaut context
func (s *Cloudinary) URL(filepath string) (string, error) {
	return s.URLCtx(s.ctx, filepath)
}

// URLCtx returns the URL of the file
// Because Cloudinary forces to have the file type in the URL, this
// method tries to download the file using each types until it finds the
// right URL
func (s *Cloudinary) URLCtx(ctx context.Context, filepath string) (string, error) {
	url, found := s.cache[filepath]
	if found {
		return url, nil
	}

	for _, typ := range fileTypes {
		url := fmt.Sprintf("%s/upload/%s", s.resBaseURL(typ), filepath)
		resp, err := s.read(url)
		if err == os.ErrNotExist {
			// if we get a os.ErrNotExist, then we try again with the next type
			continue
		}

		// We don't care about the content so we close it right away if we got any
		if err == nil {
			resp.Close()
			s.cache[filepath] = url
		}
		return url, err
	}
	// we tried all the types and didn't found anything
	return "", os.ErrNotExist
}

// Read fetches a file a returns a reader
// Will use the defaut context
func (s *Cloudinary) Read(filepath string) (io.ReadCloser, error) {
	return s.ReadCtx(s.ctx, filepath)
}

// ReadCtx fetches a file a returns a reader
// Because Cloudinary forces to have the file type in the URL, this
// method brutforces on all the possible types
func (s *Cloudinary) ReadCtx(ctx context.Context, filepath string) (io.ReadCloser, error) {
	for _, typ := range fileTypes {
		url := fmt.Sprintf("%s/upload/%s", s.resBaseURL(typ), filepath)
		resp, err := s.read(url)
		if err == os.ErrNotExist {
			// if we get a os.ErrNotExist, then we try again with the next type
			continue
		}
		return resp, err
	}
	// we tried all the types and didn't found anything
	return nil, os.ErrNotExist
}

// read fetches a file using a URL
func (s *Cloudinary) read(url string) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	// If the file does not exist we return an error
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, os.ErrNotExist
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		var pld *cdnrErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&pld); err != nil {
			return nil, err
		}
		return nil, errors.New(pld.Error.Message)
	}

	return resp.Body, nil
}

// Exists check if a file exists
// Will use the defaut context
func (s *Cloudinary) Exists(filepath string) (bool, error) {
	return s.ExistsCtx(s.ctx, filepath)
}

// ExistsCtx check if a file exists
func (s *Cloudinary) ExistsCtx(ctx context.Context, filepath string) (bool, error) {
	r, err := s.Read(filepath)
	if err == nil {
		r.Close()
		return true, nil
	}
	if err == os.ErrNotExist {
		return false, nil
	}
	return false, err
}

// Delete removes a file, ignores files that do not exist
// Will use the defaut context
func (s *Cloudinary) Delete(filepath string) error {
	return s.DeleteCtx(s.ctx, filepath)
}

// DeleteCtx removes a file
// Because Cloudinary forces to have the file type in the URL, this
// method brutforces on all the possible types
func (s *Cloudinary) DeleteCtx(ctx context.Context, filepath string) error {
	signature, timestamp := s.signature(filepath, true)
	form := url.Values{
		"api_key":    []string{s.apiKey},
		"public_id":  []string{filepath},
		"timestamp":  []string{timestamp},
		"signature":  []string{signature},
		"invalidate": []string{"true"},
	}

	for _, typ := range fileTypes {
		endpointURL := fmt.Sprintf("%s/destroy", s.apiBaseURL(typ))
		err := s.execDelete(endpointURL, strings.NewReader(form.Encode()))
		if err == os.ErrNotExist {
			continue
		}
		return err
	}

	return os.ErrNotExist
}

// execDelete makes a delete request
func (s *Cloudinary) execDelete(endpointURL string, body io.Reader) error {
	// Make the request
	req, err := http.NewRequest("POST", endpointURL, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// parse the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// For that endpoint they actually send a 200 for errors. But because the
	// documentation sucks and they don't talk about errors at all
	// let's assume that each full moon they decide to return a real error

	switch resp.StatusCode {
	case http.StatusOK:
		var pld *cdnrDeleteResponse
		if err := json.NewDecoder(resp.Body).Decode(&pld); err != nil {
			return err
		}
		switch pld.Result {
		case "ok":
			return nil
		case "not found":
			return os.ErrNotExist
		default:
			// no clue what else can be returned
			return errors.New(pld.Result)
		}
	case http.StatusNotFound:
		return os.ErrNotExist
	default:
		var pld *cdnrErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&pld); err != nil {
			return err
		}
		return errors.New(pld.Error.Message)
	}
}

// Write copy the provided os.File to dest
// Will use the defaut context
func (s *Cloudinary) Write(src io.Reader, destPath string) error {
	return s.WriteCtx(s.ctx, src, destPath)
}

// WriteCtx copy the provided os.File to dest
func (s *Cloudinary) WriteCtx(ctx context.Context, src io.Reader, destPath string) error {
	endpointURL := fmt.Sprintf("%s/upload", s.apiBaseURL("auto"))

	// REQUEST
	body := bytes.NewBufferString("")
	// attach the file
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, err := writer.CreateFormFile("file", path.Base(destPath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, src)
	// Attach the fields
	signature, timestamp := s.signature(destPath, false)
	fields := map[string]string{
		"api_key":   s.apiKey,
		"public_id": destPath,
		"timestamp": timestamp,
		"signature": signature,
	}
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return err
		}
	}
	// Close the writer
	if err := writer.Close(); err != nil {
		return err
	}
	// Make the request
	req, err := http.NewRequest("POST", endpointURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Parse the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// Parse the error
		var pld *cdnrErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&pld); err != nil {
			return err
		}
		return errors.New(pld.Error.Message)
	}
	// If the upload succeed we parse the response to cache the URL
	var pld *cdnrUploadSuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&pld); err != nil {
		return err
	}
	s.cache[destPath] = pld.SecureURL
	return nil
}

// SetAttributes sets the attributes of the file
// Will use the defaut context
func (s *Cloudinary) SetAttributes(filepath string, attrs *filestorage.UpdatableFileAttributes) (*filestorage.FileAttributes, error) {
	return s.SetAttributesCtx(s.ctx, filepath, attrs)
}

// SetAttributesCtx sets the attributes of the file
func (s *Cloudinary) SetAttributesCtx(ctx context.Context, filepath string, attrs *filestorage.UpdatableFileAttributes) (*filestorage.FileAttributes, error) {
	return filestorage.NewFileAttributesFromUpdatable(attrs), nil
}

// Attributes returns the attributes of the file
// Will use the defaut context
func (s *Cloudinary) Attributes(filepath string) (*filestorage.FileAttributes, error) {
	return s.AttributesCtx(s.ctx, filepath)
}

// AttributesCtx returns the attributes of the file
func (s *Cloudinary) AttributesCtx(ctx context.Context, filepath string) (*filestorage.FileAttributes, error) {
	return &filestorage.FileAttributes{}, nil
}

// WriteIfNotExist copies the provided io.Reader to dest if the file does
// not already exist
// Returns:
//   - A boolean specifying if the file got uploaded (true) or if already
//     existed (false).
//   - A URL to the uploaded file
//   - An error if something went wrong
// Will use the defaut context
func (s *Cloudinary) WriteIfNotExist(src io.Reader, destPath string) (new bool, url string, err error) {
	return s.WriteIfNotExistCtx(s.ctx, src, destPath)
}

// WriteIfNotExistCtx copies the provided io.Reader to dest if the file does
// not already exist
// Returns:
//   - A boolean specifying if the file got uploaded (true) or if already
//     existed (false).
//   - A URL to the uploaded file
//   - An error if something went wrong
func (s *Cloudinary) WriteIfNotExistCtx(ctx context.Context, src io.Reader, destPath string) (new bool, url string, err error) {
	return implementations.WriteIfNotExist(ctx, s, src, destPath)
}
