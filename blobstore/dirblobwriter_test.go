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
				"017b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00",
			},
		},
		"01b133787da30e1a7e5bb5f426d4915f0d4ec762b46f197649d27686c01d5eab8e94bba600cd7381d8307c1eadf720b71e62ddb2490933a6535c7c5430f9195a1" +
			"230040665bcfa4e78443d263ae4baf95c719c36b05bd7e0f595e2415033dae5ebb926eb232532753b532c657ea5c80a66315f1d34d80c49dd267a2f977ec1deb9" +
			"0378f5111c499939251beee798e7a25678d8bde5a3c2a5c92aecb33376e77db338a11be869f9b4773d5a549c4f72fb215b7a48293acc4b802284e75990b2b4791" +
			"3a915c445833b5c16308915f8",
		"01ae79b67336f8986b9776bb2931b63ae6d5dcc0c4fd668a7a12f6695aad997d09",
		"2c10cf9253442e8a7f198767490b74a1239d03f022fa3e8c4709f38c90b86b4370f90e1779c94bd975b86dd15d068c05d87ceaf7e1043360337df13552cd02a5",
	},
	{ // Directory with two simple entries, for reverse test
		testDirEntries{
			testDirEntry{
				"a",
				"c9d30a9938ecea16bed58efe5ad5b998927a56da7c8c36c1ee13292dec79aa50c5613fc90d80c37a77a5a422691d1967693a1236892e228ad95ed6fe4b505d85",
				"01504ce2f6de7e33389deb73b21f765570ad2b9f2aa8aaec8328f47b48bc3e841f",
			},
			testDirEntry{
				"empty",
				"b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
				"017b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00",
			},
		},
		"01d3ebd480aac2fb3ddabed6065a2bd1c2a689663db50e27f211add9b26e22fab818e258389444903461225b73e154bd5032b4dad88f8354ebae1eddf551e878eb462120f5fd16d2eeb" +
			"3ce3524a3dd1ea2b056c5f2b9e786654907bf9a43c71f57c56e30556746946eafea7cca639629dfed85bfbe147994151ec306f456b1b68b81ec489e74d478523f0a843b0f936c2622de" +
			"aa0966c4dc32636acc7457a973c746b8362f18582eb5f15195ee2ed552d4ba6fee0fd0a4d1898647f23cdc1a4277a4e00f5365d13c579cbc65bcbab38af6b5056248c0b331aad7f7715" +
			"445bf6d11ec030f0869eee2071fdb6c38e92f9ee72962f8c761e3f93c10c6f9f218e8ea81f65cacfe17b207ebd58e0137ea97eacb8b3bf94e5fd78fc779024095a4035a2a2445d91ea0" +
			"11b58f1173c683392d0b77596f769dfac981177ebc678eb8f48be1a0825646eeae8ec2c7999a5fbbf0fc860921312877402782b5b1a2119c2390aadcdbba70765e8aec3a52280484367" +
			"3bc5b2b57629446964518fcb4aaf7d0e18c9306a4ab434a7dbf92feb2d041f41a6ffee35864",
		"0103772a945fc1d7e62039fcfbe7e8873eff760a883133f9a6d5ebede6504b36ed",
		"9bdc297070a9cc3665305d44fd62f65644f3ccabdb364ee8953fb227cda97612de116dcd1850c50ef1a1fe2b94d802ecb0a99153cfe9e2b6c8ee79e890645f97",
	},
	{ // Directory with two simple entries, reversed data
		testDirEntries{
			testDirEntry{
				"empty",
				"b4f5a7bb878c0cec9cb4bd6ae8bb175a7ea59c1a048c5ab7c119990d0041cb9cfb67c2aa9e6fada8112719777b4b80ffada80205f8ebe6981c0ade97ff3df8e5",
				"017b54b66836c1fbdd13d2441d9e1434dc62ca677fb68f5fe66a464baadecdbd00",
			},
			testDirEntry{
				"a",
				"c9d30a9938ecea16bed58efe5ad5b998927a56da7c8c36c1ee13292dec79aa50c5613fc90d80c37a77a5a422691d1967693a1236892e228ad95ed6fe4b505d85",
				"01504ce2f6de7e33389deb73b21f765570ad2b9f2aa8aaec8328f47b48bc3e841f",
			},
		},
		"01d3ebd480aac2fb3ddabed6065a2bd1c2a689663db50e27f211add9b26e22fab818e258389444903461225b73e154bd5032b4dad88f8354ebae1eddf551e878eb462120f5fd16d2eeb" +
			"3ce3524a3dd1ea2b056c5f2b9e786654907bf9a43c71f57c56e30556746946eafea7cca639629dfed85bfbe147994151ec306f456b1b68b81ec489e74d478523f0a843b0f936c2622de" +
			"aa0966c4dc32636acc7457a973c746b8362f18582eb5f15195ee2ed552d4ba6fee0fd0a4d1898647f23cdc1a4277a4e00f5365d13c579cbc65bcbab38af6b5056248c0b331aad7f7715" +
			"445bf6d11ec030f0869eee2071fdb6c38e92f9ee72962f8c761e3f93c10c6f9f218e8ea81f65cacfe17b207ebd58e0137ea97eacb8b3bf94e5fd78fc779024095a4035a2a2445d91ea0" +
			"11b58f1173c683392d0b77596f769dfac981177ebc678eb8f48be1a0825646eeae8ec2c7999a5fbbf0fc860921312877402782b5b1a2119c2390aadcdbba70765e8aec3a52280484367" +
			"3bc5b2b57629446964518fcb4aaf7d0e18c9306a4ab434a7dbf92feb2d041f41a6ffee35864",
		"0103772a945fc1d7e62039fcfbe7e8873eff760a883133f9a6d5ebede6504b36ed",
		"9bdc297070a9cc3665305d44fd62f65644f3ccabdb364ee8953fb227cda97612de116dcd1850c50ef1a1fe2b94d802ecb0a99153cfe9e2b6c8ee79e890645f97",
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
