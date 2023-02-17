package client

import (
	"context"

	"cloud.google.com/go/storage"
)

// Creating new structs to avoid making changes to the GoCommon package
// TO DO: Move these into the common package at a later date

type CloudStorageClient storage.Client

type CloudStorageBucketHandle storage.BucketHandle

type CloudStorageObjectHandle storage.ObjectHandle

type CloudStorageWriter storage.Writer

// Executes *storage.Client.BucketHandle() on the Client instance
func (cs CloudStorageClient) Bucket(name string) ICloudStorageBucketHandle {
	bucketHandle := (*storage.Client)(&cs).Bucket(name)

	return CloudStorageBucketHandle(*bucketHandle)
}

// Executes Object() on the BucketHandle instance
func (bh CloudStorageBucketHandle) Object(name string) ICloudStorageObjectHandle {
	objectHandle := (*storage.BucketHandle)(&bh).Object(name)

	return CloudStorageObjectHandle(*objectHandle)
}

// Executes NewWriter() on the ObjectHandle instance
func (oh CloudStorageObjectHandle) NewWriter(ctx context.Context) ICloudStorageWriter {
	writer := (*storage.ObjectHandle)(&oh).NewWriter(ctx)

	return CloudStorageWriter(*writer)
}

// Executes Attrs() on the Writer instance
func (w CloudStorageWriter) Attrs() *storage.ObjectAttrs {
	return (*storage.Writer)(&w).Attrs()
}

// Executes Close() on the Writer instance
func (w CloudStorageWriter) Close() error {
	return (*storage.Writer)(&w).Close()
}

// Executes Write() on the Writer instance
func (w CloudStorageWriter) Write(p []byte) (n int, err error) {
	return (*storage.Writer)(&w).Write(p)
}

// Adds a new method to the Writer instance to set the content type of
// the bytes being written
func (w CloudStorageWriter) SetContentType(value string) {
	(*storage.Writer)(&w).ContentType = value
}
