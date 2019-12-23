package local

import (
	"context"
	"os"
	"path/filepath"

	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

// Local storage struct
type Local struct {
	path            string
	localBlobBucket *blob.Bucket
	options         *Options
}

// Options of local storage
type Options struct {
	// delete bucket when closing the storage
	DeleteOnClose bool
}

// New local storage
func New(ctx context.Context, bucketpath string, opts *Options) (*Local, error) {
	err := os.MkdirAll(filepath.Dir(bucketpath), 0744)
	if err != nil && err != os.ErrExist {
		return nil, err
	}

	// if !path.IsAbs(bucketpath) {
	// 	// get current working directory for full path
	// 	currentDir, err := os.Getwd()
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	bucketpath = path.Join(currentDir, bucketpath)
	// }

	b, err := fileblob.OpenBucket(bucketpath, nil)
	if err != nil {
		return nil, err
	}

	l := Local{
		path:            bucketpath,
		localBlobBucket: b,
		options:         opts,
	}
	return &l, nil
}

// Bucket of local storage
func (l *Local) Bucket() *blob.Bucket {
	return l.localBlobBucket
}

// Name return the name of provider
func (l *Local) Name() string {
	return objectstorage.StorageLocal
}

// BucketName return name of path in local
func (l *Local) BucketName() string {
	return l.path
}

// BucketURL return the url of blob storage bucket
func (l *Local) BucketURL() string {
	// return l.path
	return ""
}

// Close will close the local bucket
func (l *Local) Close() error {
	return l.localBlobBucket.Close()
}
