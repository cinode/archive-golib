package blobstore

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"strings"
	"testing"
)

type testDirEntry struct{ name, mimeType, bid, key string }
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
				"text/plain",
				"b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
				"017b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00",
			},
		},
		"018d2d9990d03c595e64bb0c9fbe60f13e837882e5a9c70c669f007bb8b656047229ffd9d29f0dc06201d1ac3412078e98d854b3230a320f32765022cf9326239" +
			"1f3a99e926b50794e90a9dcaf888dfac4d482546e8cd5e1ebbf881884c9b105706c0f654ac4f350398a41fbfec3bdbc782b8b49e55b2e92e9e377b6bcf4fd849f" +
			"928e5b022278fa3f5df1b5cf24433cf810d056d2942165e81010b05393697d76c330a51aed9df4933638d2380fb80220a1af7c7ec88dcfa7cf0586676cb452f4b" +
			"8b730806430fafb2843ca9dfb946c31942dea54dfd186b3",
		"01240ed7e29e427ec80f0b5ea26bdb622b61831a46716d434f9b072889c8c40ad2",
		"1d3e607118f19a66288a98c1f83e7e56490bf7c7e8707e6eda2a2e277bb0adc0606b9d511c559a028613b01dd8a41fd22e83f21a333767ee3aa25f18a7bdfe17",
	},
	{ // Directory with two simple entries, for reverse test
		testDirEntries{
			testDirEntry{
				"a",
				"text/plain",
				"c9d30a9938ecea16bed58efe5ad5b998927a56da7c8c36c1ee13292dec79aa50c5613fc90d80c37a77a5a422691d1967693a1236892e228ad95ed6fe4b505d85",
				"01504ce2f6de7e33389deb73b21f765570ad2b9f2aa8aaec8328f47b48bc3e841f",
			},
			testDirEntry{
				"empty",
				"text/plain",
				"b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
				"017b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00",
			},
		},
		"012dd708e76cc4414c4c05f8ed8f3217dbe26975b2388bccd4b56dc86fc2f93c5829b3d52a0d2c93683d9476fc6c6606f036c2b3360ac037802fab37c925027119ff20b47b81766bbb0" +
			"68d9d64ac4a69a7249f4628a61469bb05b0530a57c68ef13f150bd71995dd3da1cbf6fd021389db67a15d4d52b39bec7b51f79176e82c4d163f64028bde320d1d007c728db2ec8fb782" +
			"0266af9e029456c9343aa33b591e35c385269a6b21c2daef224d9e5797d00080432e7a4491e591b25c5fcf86b13ad0cd71188924812dc287b71f66e52ac497f5898476b1eeb77837705" +
			"b82b81abf5663ad0a942b6a955ef7ebef63ebcfb137e50ae707456d6e94c447062cbf92f32a77a1ee7b7ad4352c80fc362765892f4d25afe4b31bc901e3116b37b626c7b6a8f0a490ae" +
			"3006d0a9bb45a457712b9231700907207e474815cce6b563243ed184d2eded03c9330e280ac9e443dad2422c417f54bccedb2b02c092127d95bf79499f2db1f0d0d3d15349f17ef86f4" +
			"622e6cb2253327d30219ddb30d09bdce00eb85b68ce0b3441db5fef4a342ddc57f59927ac9171a6e670cc81888bea8a7de5e806ff603a00b4ff5e58",
		"01cab7ba0d3a81ce75fbfab9f0da673702958d777f5fed6ecd2284ae25c0a94cc6",
		"e573cdff9334c7c6e5266bf83615f4f5074117d5c471c1bafed88d8d5a721ce03926b832dc712d2590d6a887a53be83b72c884dcfcc2346fec655d5036717182",
	},
	{ // Directory with two simple entries, reversed data
		testDirEntries{
			testDirEntry{
				"empty",
				"text/plain",
				"b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
				"017b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00",
			},
			testDirEntry{
				"a",
				"text/plain",
				"c9d30a9938ecea16bed58efe5ad5b998927a56da7c8c36c1ee13292dec79aa50c5613fc90d80c37a77a5a422691d1967693a1236892e228ad95ed6fe4b505d85",
				"01504ce2f6de7e33389deb73b21f765570ad2b9f2aa8aaec8328f47b48bc3e841f",
			},
		},
		"012dd708e76cc4414c4c05f8ed8f3217dbe26975b2388bccd4b56dc86fc2f93c5829b3d52a0d2c93683d9476fc6c6606f036c2b3360ac037802fab37c925027119ff20b47b81766bbb0" +
			"68d9d64ac4a69a7249f4628a61469bb05b0530a57c68ef13f150bd71995dd3da1cbf6fd021389db67a15d4d52b39bec7b51f79176e82c4d163f64028bde320d1d007c728db2ec8fb782" +
			"0266af9e029456c9343aa33b591e35c385269a6b21c2daef224d9e5797d00080432e7a4491e591b25c5fcf86b13ad0cd71188924812dc287b71f66e52ac497f5898476b1eeb77837705" +
			"b82b81abf5663ad0a942b6a955ef7ebef63ebcfb137e50ae707456d6e94c447062cbf92f32a77a1ee7b7ad4352c80fc362765892f4d25afe4b31bc901e3116b37b626c7b6a8f0a490ae" +
			"3006d0a9bb45a457712b9231700907207e474815cce6b563243ed184d2eded03c9330e280ac9e443dad2422c417f54bccedb2b02c092127d95bf79499f2db1f0d0d3d15349f17ef86f4" +
			"622e6cb2253327d30219ddb30d09bdce00eb85b68ce0b3441db5fef4a342ddc57f59927ac9171a6e670cc81888bea8a7de5e806ff603a00b4ff5e58",
		"01cab7ba0d3a81ce75fbfab9f0da673702958d777f5fed6ecd2284ae25c0a94cc6",
		"e573cdff9334c7c6e5266bf83615f4f5074117d5c471c1bafed88d8d5a721ce03926b832dc712d2590d6a887a53be83b72c884dcfcc2346fec655d5036717182",
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
			dw.AddEntry(DirEntry{
				Name:     entry.name,
				MimeType: entry.mimeType,
				Bid:      entry.bid,
				Key:      entry.key,
			})
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
