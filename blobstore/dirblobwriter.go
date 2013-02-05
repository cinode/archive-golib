// Copyright 2013 The Cinode Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blobstore

import (
	"bytes"
	"io"
	"sort"
)

// Helper structure for holding one directory entry
type dirEntry struct {
	name, bid, key string
}

// Helper for sorting by name
type sortByName []*dirEntry

func (s sortByName) Len() int {
	return len(s)
}

func (s sortByName) Less(i, j int) bool {
	return s[i].name < s[j].name
}

func (s sortByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Structure used for creating new directory blobs
type DirBlobWriter struct {

	// Storage Object
	Storage BlobStorage

	// A list of currently handled entries
	entries []*dirEntry
}

// Adds a new entry to the directory
// TODO: Don't allow adding duplicated entries 
func (d *DirBlobWriter) AddEntry(name, bid, key string) error {
	d.entries = append(d.entries,
		&dirEntry{
			name: name,
			bid:  bid,
			key:  key})
	return nil
}

func (d *DirBlobWriter) Finalize() (bid string, key string, err error) {
	if len(d.entries) <= maxSimpleDirEntries {
		return d.finalizeSimple()
	}
	return d.finalizeSplit()
}

func (d *DirBlobWriter) finalizeSimple() (bid string, key string, err error) {

	// Sort entries by name
	sort.Sort(sortByName(d.entries))

	// Serialize the data
	var buffer bytes.Buffer
	buffer.WriteByte(blobTypeSimpleStaticDir)

	// Number of entries first
	serializeInt(int64(len(d.entries)), &buffer)

	// All entries right after
	for _, entry := range d.entries {
		serializeString(entry.name, &buffer)
		serializeString(entry.bid, &buffer)
		serializeString(entry.key, &buffer)
	}

	// Create blob out of the data
	return createHashValidatedBlobFromReaderGenerator(
		func() io.Reader { return bytes.NewReader(buffer.Bytes()) },
		d.Storage)
}

func (d *DirBlobWriter) finalizeSplit() (bid string, key string, err error) {
	panic("Unimplemented: split dir blob")
}
