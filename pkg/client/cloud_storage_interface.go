package client

import (
	"context"

	"cloud.google.com/go/storage"
)

type ICloudStorageClient interface {
	Bucket(name string) ICloudStorageBucketHandle
}

type ICloudStorageBucketHandle interface {
	Object(name string) ICloudStorageObjectHandle
}

type ICloudStorageObjectHandle interface {
	NewWriter(ctx context.Context) ICloudStorageWriter
}

type ICloudStorageWriter interface {
	Attrs() *storage.ObjectAttrs
	Close() error
	Write(p []byte) (n int, err error)
	SetContentType(value string)
}
