// Copyright 2014 The Cinode Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package localstorage

import (
	"io"
)

// This interface will represent an object one can use to read
// data from a local blob storage
type Reader interface {
	io.ReadCloser
}
