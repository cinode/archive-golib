package blobstore

import "errors"

var (
	ErrInvalidValidationMethod = errors.New("Invalid blob validation method")

	ErrInvalidFileBlobType              = errors.New("Invalid blob type - not a file blob")
	ErrInvalidSplitFileSize             = errors.New("Invalid size of a split file")
	ErrMalformedSplitFileSizePartsCount = errors.New("Invalid split file blob - number of partial blobs is incorrect")
	ErrMalformedSplitFileExtraData      = errors.New("Invalid split file blob - extra bytes found at the end of the blob")
	ErrMalformedSplitFileExtraDataPart  = errors.New("Invalid split file blob - extra bytes found at the end of the partial blob")
	ErrInvalidFileSubBlobType           = errors.New("Invalid sub blob type - not a file blob")

	ErrMalformedDirInvalidEntriesCount = errors.New("Invalid directory blob - incorrect number of entries found")
	ErrMalformedDirExtraData           = errors.New("Invalid directory blob - extra bytes found at the end")
	ErrNoMoreDirEntries                = errors.New("No more directory entries found")
)
