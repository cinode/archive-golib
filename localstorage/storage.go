// Copyright 2013-2014 The Cinode Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package localstorage

// Storage is an abstraction around local storage for blobs
type Storage interface {

	// GetBlobReader creates new reader for existing blob or fails with an error
	GetBlobReader(blobId string) (reader Reader, err error)

	// GetBlobWriter creates a writer for new blob
	GetBlobWriter() (writer Writer, err error)
}
