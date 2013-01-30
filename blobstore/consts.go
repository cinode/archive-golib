package blobstore

const (
	blobTypeSimpleStaticFile = 0x01
	blobTypeSimpleStaticDir  = 0x11

	maxSimpleFileDataSize = 16 * 1024 * 1024
	maxSimpleDirEntries   = 1024

	validationMethodHash = 0x01
)
