package util

import "compress/zlib"

type ZlibReadSeeker struct {
}

func NewZlibReadSeeker(closer ReadSeekCloser) (ZlibReadSeeker, error) {
	zlib.NewReader()
}
