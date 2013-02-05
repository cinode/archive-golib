package blobstore

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"strings"
	"testing"
)

type testDirEntry struct{ name, bid, key string }
type testDirEntries []testDirEntry

var simpleDirTests = []struct {
	entries testDirEntries
	blobHex string
	key     string
	bid     string
}{
	{ // Empty Directory
		testDirEntries{},
		"01c659",
		"0129d7159641f64847d66fc4091d1320ff201147e2ca7e221080ce08933f1e1fd3",
		"cc347605074b230f9ca42f53c0f16475e3560df75c9378c0e9f7608781a6a04127f178bd428a10c1442b608e239148283a9e52f3bf0efdf514dfd7e1f9326372",
	},
	{ // Directory with one empty file
		testDirEntries{
			testDirEntry{
				"empty",
				"b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
				"017b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00"},
		},
		"01b133787da30e1a7e5bb5f426d4915f0d4ec762b46f197649d27686c01d5eab8e94bba600cd7381d8307c1eadf720b71e62ddb2490933a6535c7c5430f9195a1" +
			"230040665bcfa4e78443d263ae4baf95c719c36b05bd7e0f595e2415033dae5ebb926eb232532753b532c657ea5c80a66315f1d34d80c49dd267a2f977ec1deb9" +
			"0378f5111c499939251beee798e7a25678d8bde5a3c2a5c92aecb33376e77db338a11be869f9b4773d5a549c4f72fb215b7a48293acc4b802284e75990b2b4791" +
			"3a915c445833b5c16308915f8",
		"01ae79b67336f8986b9776bb2931b63ae6d5dcc0c4fd668a7a12f6695aad997d09",
		"2c10cf9253442e8a7f198767490b74a1239d03f022fa3e8c4709f38c90b86b4370f90e1779c94bd975b86dd15d068c05d87ceaf7e1043360337df13552cd02a5",
	},
}

func TestSimpleDirs(t *testing.T) {
	for _, test := range simpleDirTests {
		key := strings.Replace(test.key, " ", "", -1)
		bid := strings.Replace(test.bid, " ", "", -1)
		blob, _ := hex.DecodeString(strings.Replace(test.blobHex, " ", "", -1))

		m := NewMemoryBlobStorage()
		dw := DirBlobWriter{Storage: m}
		for _, entry := range test.entries {
			dw.AddEntry(entry.name, entry.bid, entry.key)
		}

		rbid, rkey, err := dw.Finalize()
		if err != nil {
			t.Error(err)
		}

		if rbid != bid {
			t.Errorf("Invalid blob id generated, got: %v..., expected: %v...", rbid[:16], bid[:16])
		}

		if rkey != key {
			t.Errorf("Invalid key generated, got: %v..., expected: %v...", rkey[:16], key[:16])
		}

		reader, err := m.NewBlobReader(rbid)
		if err != nil {
			t.Errorf("Couldn't open the blob with id: %v... for reading: %v", rbid[:16], err)
		} else {

			readBytes, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Errorf("Couldn't read the blob with id: %v..., error: %v", rbid[:16], err)
			} else if !bytes.Equal(readBytes, blob) {

				readBytesHex, blobHex := "", ""

				if len(readBytes) > 10 {
					readBytesHex = hex.EncodeToString(readBytes[:10]) + "..."
				} else {
					readBytesHex = hex.EncodeToString(readBytes)
				}

				if len(blob) > 10 {
					blobHex = hex.EncodeToString(blob[:10]) + "..."
				} else {
					blobHex = hex.EncodeToString(blob)
				}

				t.Errorf("The blob with id: %v... has invalid content, got: %v, expected %v",
					rbid[:16],
					readBytesHex,
					blobHex)
			}
		}
	}
}
