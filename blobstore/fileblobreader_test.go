package blobstore

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func putBlob(storage BlobStorage, bid string, content []byte) {
	writer, _ := storage.NewBlobWriter(bid)
	writer.Write(content)
	writer.Finalize()
}

func TestEmptyFileBlob(t *testing.T) {

	storage := NewMemoryBlobStorage()

	content := []byte{}

	putBlob(storage,
		"b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
		[]byte{0x01, 0xeb})

	rdr := FileBlobReader{Storage: storage}
	err := rdr.Open(
		"b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
		"017b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00")
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(&rdr)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(content, data) {
		t.Error("Invalid blob content")
	}
}
