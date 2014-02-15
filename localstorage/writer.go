// Copyright 2014 The Cinode Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package localstorage

import (
	"io"
)

// This interface will represent an object one can use to write
// into local blob storage. Make sure that the operation ends with
// either Commit(blobId) or Rollback() function. Otherwise the
// behavior is undefined.
type Writer interface {
	io.Writer

	// Commit changes written so far, after this function is called,
	// all subsequent Write and Rollback calls will fail
	Commit(blobId string) error

	// Rollback changes to the blob, cleanup any potential garbage,
	// after this function is called, all subsequent Write and Commit
	// calls will fail
	Rollback() error
}
