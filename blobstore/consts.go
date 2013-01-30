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

	validationMethodHash = 0x01
)
