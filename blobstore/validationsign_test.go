package blobstore

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	//"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestSimpleWriteReadCycle(t *testing.T) {

	// First we need to generate some RSA private key,
	// small one to take short amount of time
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal("Could not generate test RSA key")
	}

	storage := NewMemoryBlobStorage()

	testData := []byte("Hello world!")

	// Generate the blob
	bid, key, err := createSignValidatedBlobFromReaderGenerator(func() io.Reader {
		return bytes.NewReader(testData)
	}, privKey, 832, storage)
	if err != nil {
		t.Fatal("Could not create signed blob:", err)
	}

	reader, err := createReaderForSignedBlob(bid, key, storage)
	if err != nil {
		t.Fatal("Could not create signed blob reader:", err)
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal("Could not read signed blob content:", err)
	}

	if !bytes.Equal(data, testData) {
		t.Fatal("Invalid data read from the blob", data, testData)
	}
}
