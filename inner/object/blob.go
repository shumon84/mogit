package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/shumon84/binutil"
	"github.com/shumon84/mogit/inner/util"
)

type Blob struct {
	rsc    util.ReadSeekCloser
	size   int64
	sha1   []byte
	encode []byte
	decode []byte
}

func NewBlobFromPath(path string) (*Blob, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return NewBlob(file, stat.Size())
}

func NewBlob(rsc util.ReadSeekCloser, size int64) (*Blob, error) {
	if rsc == nil {
		return nil, ErrNilReadCloser
	}
	if size < 0 {
		return nil, ErrNegativeSize
	}
	return &Blob{
		rsc:  rsc,
		size: size,
	}, nil
}

func (b *Blob) SHA1() ([]byte, error) {
	if b.sha1 != nil {
		digest := make([]byte, len(b.sha1))
		copy(digest, b.sha1)
		return digest, nil
	}
	data, err := b.Decode()
	if err != nil {
		return nil, err
	}
	h := sha1.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}

	digest := h.Sum(nil)
	b.sha1 = make([]byte, len(digest))
	copy(b.sha1, digest)
	return digest, nil
}

func (b *Blob) Type() ObjectType {
	return BlobObject
}

func (b *Blob) Decode() ([]byte, error) {
	if b.decode != nil {
		data := make([]byte, len(b.decode))
		copy(data, b.decode)
		return data, nil
	}

	header := append([]byte(fmt.Sprintf("blob %d", b.size)), 0)

	if _, err := b.rsc.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	r := binutil.NewReader(b.rsc)
	data, err := r.Bytes(int(b.size))
	if err != nil {
		return nil, err
	}

	data = append(header, data...)
	b.decode = make([]byte, len(data))
	copy(b.decode, data)
	return data, nil
}

func (b *Blob) Encode() ([]byte, error) {
	if b.encode != nil {
		data := make([]byte, len(b.encode))
		copy(data, b.encode)
		return data, nil
	}
	rawData, err := b.Decode()
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	zr := zlib.NewWriter(buf)
	defer zr.Close()

	if _, err := zr.Write(rawData); err != nil {
		return nil, err
	}
	//if err := zr.Flush(); err != nil {
	//	return nil, err
	//}

	b.encode = buf.Bytes()
	return buf.Bytes(), nil
}

func (b *Blob) Close() error {
	return b.rsc.Close()
}
