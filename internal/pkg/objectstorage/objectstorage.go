package objectstorage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

	"gocloud.dev/blob"
)

// storage type of artifact
const (
	// local storage for testing
	StorageLocal = "local"
	// google cloud storage
	StorageGCS = "gcs"
	// amazon s3 storage
	StorageS3 = "s3"
	// digital ocean space storage
	StorageDO = "do"
	// minio storage
	StorageMinio = "minio"
)

// list of error
var (
	ErrByteEmpty        = errors.New("byte content is empty")
	ErrCredentialsEmpty = errors.New("credentials is empty")
)

// StorageProvider interface
type StorageProvider interface {
	Bucket() *blob.Bucket
	Name() string
	BucketName() string
	BucketURL() string
	Close() error
}

// Storage struct
type Storage struct {
	storage StorageProvider
}

// ReadOptions struct
// we might need this later on
type ReadOptions struct {
	// FileMode is an options when downloading file
	// using DownloadFile function
	FileMode os.FileMode
}

// WriteOptions struct
type WriteOptions struct {
	// BufferSize for writing many small writes concurrently
	BufferSize int
	// ContentType specifies the MIME type of the blob
	ContentType string
	// ContentDisposition specifies whether the content is displayed inline or as attachment
	ContentDisposition string
	// ContentEncoding to store specific encoding for the content
	ContentEncoding string
	// ContentLanguage for language of the content
	ContentLanguage string
	// ContentMD5 for integrity check
	ContentMD5 []byte
	// Key-value associated with the blob
	Metadata map[string]string
}

// New artifact
func New(storage StorageProvider) *Storage {
	return &Storage{storage}
}

// Attributes return information/attributes of object
func (s *Storage) Attributes(ctx context.Context, key string) (*blob.Attributes, error) {
	return s.storage.Bucket().Attributes(ctx, key)
}

// SignedURL to create a temporary URL to download a private file
func (s *Storage) SignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	return s.storage.Bucket().SignedURL(ctx, key, &blob.SignedURLOptions{Expiry: expiry})
}

// Upload file from bytes
func (s *Storage) Upload(ctx context.Context, reader io.Reader, key string, writeOptions *WriteOptions) (string, error) {
	return s.upload(ctx, key, reader, writeOptions)
}

// UploadFile file from source and destination
func (s *Storage) UploadFile(ctx context.Context, filepath, key string, writeOptions *WriteOptions) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	return s.upload(ctx, key, f, writeOptions)
}

// UploadByte is a halper function for upload receiving byte param
func (s *Storage) UploadByte(ctx context.Context, content []byte, key string, writeOptions *WriteOptions) (string, error) {
	if len(content) == 0 {
		return "", ErrByteEmpty
	}

	buff := bytes.NewBuffer(content)
	return s.upload(ctx, key, buff, writeOptions)
}

// Download file from bucket
func (s *Storage) Download(ctx context.Context, key string, readOptions *ReadOptions) (io.Reader, error) {
	return s.download(ctx, key, readOptions)
}

// DownloadFile will download and create file from object storage
func (s *Storage) DownloadFile(ctx context.Context, key, destination string, readOptions *ReadOptions) error {
	reader, err := s.download(ctx, key, readOptions)
	if err != nil {
		return err
	}
	defer reader.Close()

	// means owner able to execute and read
	// group only able to read
	// others only able to read
	var fileMode os.FileMode = 0744
	// use filemode from options if filemode is not empty
	if readOptions != nil && readOptions.FileMode != 0 {
		fileMode = readOptions.FileMode
	}

	f, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR, fileMode)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}

	return nil
}

// DownloadByte is a helper function for downloading content with return of byte
func (s *Storage) DownloadByte(ctx context.Context, key string, readOptions *ReadOptions) ([]byte, error) {
	reader, err := s.download(ctx, key, readOptions)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

// upload content to object storage
// the function return the path of uploaded object and error
func (s *Storage) upload(ctx context.Context, key string, reader io.Reader, writeOptions *WriteOptions) (string, error) {
	uploadPath := path.Join(s.storage.BucketURL(), key)
	blobBucket := s.storage.Bucket()

	var (
		result []byte
		err    error
		opts   *blob.WriterOptions
	)

	if writeOptions != nil {
		opts = &blob.WriterOptions{
			BufferSize:         writeOptions.BufferSize,
			ContentType:        writeOptions.ContentType,
			ContentDisposition: writeOptions.ContentDisposition,
			ContentEncoding:    writeOptions.ContentEncoding,
			ContentLanguage:    writeOptions.ContentLanguage,
			ContentMD5:         writeOptions.ContentMD5,
			Metadata:           writeOptions.Metadata,
		}
	}

	result, err = ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	nw, err := blobBucket.NewWriter(ctx, key, opts)
	if err != nil {
		return "", err
	}

	_, err = nw.Write(result)
	if err != nil {
		return "", err
	}

	// writer is asynchronous
	// need to close to make sure writer is error or not
	if err := nw.Close(); err != nil {
		return "", err
	}
	return uploadPath, err
}

func (s *Storage) download(ctx context.Context, key string, readOptions *ReadOptions) (*blob.Reader, error) {
	var (
		opts *blob.ReaderOptions
		err  error
	)

	if readOptions != nil {
		opts = &blob.ReaderOptions{}
	}

	bucket := s.storage.Bucket()
	reader, err := bucket.NewReader(ctx, key, opts)
	return reader, err
}

// Name of provider
// this might be useful if application has admin-port of something like that
// to retrieve the current name of the provider
func (s *Storage) Name() string {
	return s.storage.Name()
}

// BucketName of storage provider
func (s *Storage) BucketName() string {
	return s.storage.BucketName()
}

// Close will close the object storage bucket and return error
func (s *Storage) Close() error {
	return s.storage.Close()
}

// Stream create a new stream object
func (s *Storage) Stream(ctx context.Context, key string, writeOptions *WriteOptions) (*Stream, error) {
	blobBucket := s.storage.Bucket()
	st := Stream{
		bucket: blobBucket,
	}
	return &st, nil
}

// Stream struct
type Stream struct {
	bucket *blob.Bucket
}

// Reader return blob reader
func (s *Stream) Reader(ctx context.Context, key string, readOptions *ReadOptions) (*blob.Reader, error) {
	var opts *blob.ReaderOptions
	if readOptions != nil {
		opts = &blob.ReaderOptions{}
	}
	return s.bucket.NewReader(ctx, key, opts)
}

// Writer return blob writer
func (s *Stream) Writer(ctx context.Context, key string, writeOptions *WriteOptions) (*blob.Writer, error) {
	var opts *blob.WriterOptions
	if writeOptions != nil {
		opts = &blob.WriterOptions{
			BufferSize:         writeOptions.BufferSize,
			ContentType:        writeOptions.ContentType,
			ContentDisposition: writeOptions.ContentDisposition,
			ContentEncoding:    writeOptions.ContentEncoding,
			ContentLanguage:    writeOptions.ContentLanguage,
			ContentMD5:         writeOptions.ContentMD5,
			Metadata:           writeOptions.Metadata,
		}
	}
	return s.bucket.NewWriter(ctx, key, opts)
}
