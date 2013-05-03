package blobstore

import (
	"io"
)

type baseBlobReader struct {
}

// Internal function, try to open a blob having it's bid and key,
// don't interpret anything but blob's type
func (r *baseBlobReader) openInternal(
	storage BlobStorage, bid, key string, requiredValidationMethod int64) (
	reader io.Reader, blobType int64, err error) {

	// Get the raw blob reader
	if reader, err = storage.NewBlobReader(bid); err != nil {
		return
	}

	// Find out the validation method
	validationMethod, err := deserializeInt(reader)
	if err != nil {
		return
	}

	// File blobs must use the hash-based validation
	// TODO: We may relax this if we start using links and decide to dereference links here
	if validationMethod != requiredValidationMethod {
		return nil, 0, ErrInvalidValidationMethod
	}

	// Get the unencrypted stream
	if reader, err = createReaderForHashBlobData(reader, key); err != nil {
		return
	}

	// See what type of a blob this is
	blobType, err = deserializeInt(reader)
	if err != nil {
		return
	}

	return
}
