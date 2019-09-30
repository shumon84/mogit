package commit

import "errors"

var (
	ErrNilReadCloser     = errors.New("io.ReadCloser is nil")
	ErrNegativeSize      = errors.New("size is negative number")
	ErrInvalidHashLength = errors.New("invalid hash length")
)
