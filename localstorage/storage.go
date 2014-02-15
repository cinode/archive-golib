// Copyright 2013-2014 The Cinode Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package localstorage

// An interface usefull for blob storage operations
type Storage interface {

	// Create new writer for blobs
	GetBlobWriter() (writer Writer, err error)

	// Create new reader for existing blob
	GetBlobReader(blobId string) (reader Reader, err error)
}
