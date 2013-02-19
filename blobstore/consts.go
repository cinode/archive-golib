// Copyright 2013 The Cinode Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blobstore

const (
	blobTypeSimpleStaticFile = 0x01
	blobTypeSplitStaticFile  = 0x02
	blobTypeSimpleStaticDir  = 0x11
	blobTypeSplitStaticDir   = 0x12

	cipherAES256    = 0x01
	cipherAES256Hex = "01"

	maxSimpleFileDataSize = 16 * 1024 * 1024
	maxSimpleDirEntries   = 1024

	maxSaneSplitFileParts = 1024 * 1024
	maxSaneBidLength      = 1024
	maxSaneKeyLength      = 16 * 1024

	validationMethodHash = 0x01
)
