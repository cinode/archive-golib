package blobstore

import (
	"testing"
)

func genTestDirData() (BlobStorage, *DirBlobWriter, DirBlobReader) {

	storage := NewMemoryBlobStorage()
	writer := DirBlobWriter{Storage: storage}
	reader := NewDirBlobReader(storage)

	return storage, &writer, reader
}

func testMultipleEntriesDir(t *testing.T, entries []DirEntry) {

	entryMap := make(map[string]DirEntry)
	for _, entry := range entries {
		entryMap[entry.Name] = entry
	}

	_, w, r := genTestDirData()

	for _, entry := range entryMap {
		w.AddEntry(entry)
	}

	bid, key, err := w.Finalize()
	if err != nil {
		t.Error(err)
	}

	if err = r.Open(bid, key); err != nil {
		t.Error(err)
	}

	for len(entryMap) > 0 {

		if !r.IsNextEntry() {
			t.Error("Missing directory entries")
		}

		entry, err := r.NextEntry()
		if err != nil {
			t.Error(err)
		}

		entry2, ok := entryMap[entry.Name]
		if !ok {
			t.Error("Read unknown entry: " + entry.Name)
		}

		if entry != entry2 {
			t.Error("Entries do not match: " + entry.Name)
		}

		delete(entryMap, entry.Name)
	}

	if r.IsNextEntry() {
		t.Error("Extra entries found in the directory")
	}

	if _, err = r.NextEntry(); err != ErrNoMoreDirEntries {
		t.Error("Invalid error returned when trying to read entry past the end of directory", err)
	}
}

var testVector = [][]DirEntry{
	// Empty vector
	{},
	// Single entry
	{
		{Name: "test.txt", MimeType: "mime", Key: "key", Bid: "bid"},
	},
	// Multiple entries
	{
		{Name: "test.txt", MimeType: "mime", Key: "key", Bid: "bid"},
		{Name: "test2.txt", MimeType: "mime2", Key: "key2", Bid: "bid2"},
		{Name: "test3.txt", MimeType: "mime3", Key: "key3", Bid: "bid3"},
		{Name: "test4.txt", MimeType: "mime4", Key: "key4", Bid: "bid4"},
	},
	// Multiple entries, different order
	{
		{Name: "test3.txt", MimeType: "mime3", Key: "key3", Bid: "bid3"},
		{Name: "test2.txt", MimeType: "mime2", Key: "key2", Bid: "bid2"},
		{Name: "test.txt", MimeType: "mime", Key: "key", Bid: "bid"},
		{Name: "test4.txt", MimeType: "mime4", Key: "key4", Bid: "bid4"},
	},
}

func TestDirVectors(t *testing.T) {
	for _, data := range testVector {
		testMultipleEntriesDir(t, data)
	}
}
