// Copyright 2013 The Cinode Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blobstore

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
)

var (
	ErrInvalidValidationMethod = errors.New("Invalid validation method")
)

func createHashValidatedBlobFromReaderGenerator(readerGenerator func() io.Reader, storage BlobStorage) (bid string, key string, err error) {

	// Generate the key
	hasher := sha512.New()
	io.Copy(hasher, readerGenerator())
	keySource := hasher.Sum(nil)

	// Generate the encrypted content
	encryptedBuffer := bytes.Buffer{}
	encryptedWriter, key, err := createEncryptor(keySource, nil, &encryptedBuffer)
	if err != nil {
		return
	}
	io.Copy(encryptedWriter, readerGenerator())

	// Generate blob id
	hasher.Reset()
	io.Copy(hasher, bytes.NewReader(encryptedBuffer.Bytes()))
	bid = hex.EncodeToString(hasher.Sum(nil))

	// Finally generate the blob itself
	blobWriter, err := storage.NewBlobWriter(bid)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			blobWriter.Cancel()
		}
	}()
	if _, err = blobWriter.Write([]byte{validationMethodHash}); err != nil {
		return
	}
	if _, err = io.Copy(blobWriter, &encryptedBuffer); err != nil {
		return
	}
	if err = blobWriter.Finalize(); err != nil {
		return
	}

	// Ok, we're done here
	return
}

func createReaderForHashBlobData(reader io.Reader, key string) (rawReader io.Reader, err error) {
	return createDecryptor(key, nil, reader)
}

func createReaderForHashBlob(bid string, key string, storage BlobStorage) (rawReader io.Reader, err error) {

	// Get the reader
	encryptedReader, err := storage.NewBlobReader(bid)
	if err != nil {
		return
	}

	// Test the validation method
	validationType, err := deserializeInt(encryptedReader)
	if err != nil {
		return
	}
	if validationType != validationMethodHash {
		return nil, ErrInvalidValidationMethod
	}

	// Get the encryptor
	return createReaderForHashBlobData(encryptedReader, key)
}
