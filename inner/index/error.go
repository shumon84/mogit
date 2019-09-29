package index

import "errors"

var (
	ErrNilReader             = errors.New("reader is nil")
	ErrInvalidSignature      = errors.New("invalid signature")
	ErrNilFileInfo           = errors.New("file info is nil")
	ErrForbiddenPermission   = errors.New("forbidden permission")
	ErrNotSupportedVersion   = errors.New("this git index version is not supported")
	ErrIndexOutOfEntryRanges = errors.New("index out of entry ranges")
)
