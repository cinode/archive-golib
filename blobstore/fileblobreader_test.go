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

func genericFileBlobTest(t *testing.T, bid, key string, blobContent, fileContent []byte) {

	storage := NewMemoryBlobStorage()

	putBlob(storage, bid, blobContent)

	rdr := NewFileBlobReader(storage)
	err := rdr.Open(bid, key)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(fileContent, data) {
		t.Error("Invalid blob content")
	}
}

func TestEmptyFileBlob(t *testing.T) {
	genericFileBlobTest(t,
		"b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
		"017b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00",
		[]byte{0x01, 0xeb},
		[]byte{})
}

func TestSingleAFileBlob(t *testing.T) {
	genericFileBlobTest(t,
		"c9d30a9938ecea16bed58efe5ad5b998927a56da7c8c36c1ee13292dec79aa50c5613fc90d80c37a77a5a422691d1967693a1236892e228ad95ed6fe4b505d85",
		"01504ce2f6de7e33389deb73b21f765570ad2b9f2aa8aaec8328f47b48bc3e841f",
		[]byte{0x01, 0x8f, 0x14},
		[]byte{0x61})
}

func TestHelloWorldFileBlob(t *testing.T) {
	genericFileBlobTest(t,
		"82aeef202165cf11930ea44a9ad8337aea355d63751a7260552e3e014ad6313bca69c83fa4e3555531d44a1025708183784af0e2002562b7260559ce0e7af262",
		"01ac9d259134ccef987f9f4df3115b0b7a24b379cbebb2aaa91ed811c8cf5e0907",
		[]byte{0x01, 0x85, 0x5e, 0x29, 0x6f, 0x95, 0xd1, 0xea, 0xf3, 0xfe, 0xb7, 0xd4, 0x8c, 0xe0},
		[]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21})
}

func TestAlphabetFileBlob(t *testing.T) {
	genericFileBlobTest(t,
		"4cfb056a184d4377eff9fc3e8364906af4b3b3c9467c2fb8245382bdd535ea17f8a63abc190a92539bd9295152f112d3365d4910737b9f9f3e0eb2f2eef40648",
		"01b11ef5debd728940485629e342c572bcc5b103d7b56de27b07f901b4abcdb5d4",
		[]byte{0x01, 0xf0, 0xea, 0xd9, 0x42, 0x12, 0x73, 0x7b, 0x28, 0x60, 0xea, 0x35, 0xe3, 0x1c, 0x7d, 0xd1, 0x76, 0xb5, 0x62, 0x09, 0x68, 0x2c, 0x3a,
			0x67, 0x92, 0x1d, 0x46, 0x48, 0x23, 0x13, 0xc2, 0x45, 0xd4, 0x55, 0x1c, 0x76, 0x5c, 0x3c, 0xa8, 0x51, 0xd7, 0xf3, 0x75, 0x91, 0x1a, 0x66,
			0xe6, 0xb5, 0x2b, 0x65, 0x0d, 0x51, 0xea, 0xc3},
		[]byte{0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77,
			0x78, 0x79, 0x7a, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54,
			0x55, 0x56, 0x57, 0x58, 0x59, 0x5a})
}

func TestSplitAaaaFile(t *testing.T) {

	storage := NewMemoryBlobStorage()

	writer := FileBlobWriter{Storage: storage}

	// Fill the blob data with 'a' crossing the simple file blob capacity by 1 byte
	b := make([]byte, 1024)
	for i, _ := range b {
		b[i] = 'a'
	}
	for i := 0; i < 16*1024; i++ {
		writer.Write(b)
	}
	writer.Write(b[:1])
	bid, key, err := writer.Finalize()
	if err != nil {
		t.Fatal(err)
	}
	if bid != "f8615f370c23b1bf7b654ed19aadc5e2011ff98d139cd1a05be588a8f4d03af375f3598a10b138e9106702945c7c1642827fa807d70a44454585ec5251d45b8a" ||
		key != "01bffd8d7830029b88367640a067ce1e0220a929fdd20c0a9157f6e1e094b19ff2" {
		t.Fatal("Invalid blob generated for testing")
	}

	rdr := NewFileBlobReader(storage)
	if err = rdr.Open(bid, key); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		t.Fatal(err)
	}

	totalSize := 16*1024*1024 + 1
	if len(data) != totalSize {
		t.Fatalf("Blob size does not match, expected: %v, got %v", totalSize, len(data))
	}

	for i := 0; i < totalSize; i += 1024 {
		blockLen := totalSize - i
		if blockLen > 1024 {
			blockLen = 1024
		}
		if !bytes.Equal(data[i:i+blockLen], b[:blockLen]) {
			t.Error("Invalid blob content")
		}
	}

}
