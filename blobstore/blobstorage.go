package blobstore

import "io"

// An interface usefull for blob storage operations
type BlobStorage interface {

	// Create new writer for blobs
	NewBlobWriter(blobId string) io.WriteCloser
}
