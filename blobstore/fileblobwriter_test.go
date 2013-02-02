package blobstore

import (
	"encoding/hex"
	"strings"
	"testing"
)

var simpleTests = []struct {
	content string
	blobHex string
	key     string
	bid     string
}{
	{
		"",
		"01 eb",
		"01 7b54b668 36c1fbdd 13d2441d 9e1434dc 62ca677f b68f5fe6 6a464baa decdbd00",
		"b4f5a7bb 878c0cec 9cb4bd6a e8bb175a 7ea59c1a 048c5ab7 c119990d 0041cb9c fb67c2aa 9e6fada8 11271977 7b4b80ff ada80205 f8ebe698 1c0ade97 ff3df8e5"},
}

func TestEmptyFile(t *testing.T) {

	for _, test := range simpleTests {

		key := strings.Replace(test.key, " ", "", -1)
		bid := strings.Replace(test.bid, " ", "", -1)
		blob, _ := hex.DecodeString(strings.Replace(test.blobHex, " ", "", -1))

		m := NewMemoryBlobStorage()
		bw := FileBlobWriter{Storage: m}
		bw.Write([]byte(test.content))

		rbid, rkey, err := bw.Finalize()
		if err != nil {
			t.Error(err)
		}

		if rbid != bid {
			t.Errorf("Invalid blob id generated, got: %v, expected: %v", rbid, bid)
		}

		if rkey != key {
			t.Error("Invalid key generated, got: %v, expected: %v", rkey, key)
		}
		// TODO: Test the blob content
		blob = blob
	}
}
