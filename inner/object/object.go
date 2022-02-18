package object

import (
	"compress/zlib"
	"io"
	"os"
	"path/filepath"

	"github.com/shumon84/binutil"

	"github.com/shumon84/mogit/inner/util"
)

type ObjectType int

const (
	BlobObject ObjectType = iota
	TreeObject
	CommitObject
	TagObject
)

func (o ObjectType) String() string {
	switch o {
	case BlobObject:
		return "blob"
	case TreeObject:
		return "tree"
	case CommitObject:
		return "object"
	case TagObject:
		return "tag"
	default:
		return "undefined"
	}
}

type Object interface {
	io.Closer
	SHA1() ([]byte, error)
	Type() ObjectType
	Decode() ([]byte, error)
	Encode() ([]byte, error)
}

func GetObjectPath(digest string) (string, error) {
	if len(digest) != 20 {
		return "", ErrInvalidHashLength
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err := util.FindGitRoot(wd)
	if err != nil {
		return "", err
	}
	return filepath.Join(path, digest[:2], digest[2:]), nil
}

func ReadObject(digest string) (Object, error) {
	path, err := GetObjectPath(digest)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	zlibr, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	// TODO: read git object
	r := binutil.NewReader(zlibr)
	objectTypeString, err := r.String()
	if err != nil {
		return nil, err
	}
	var object Object
	switch objectTypeString {
	case BlobObject.String():
		object, err = NewBlobFromPath(path)
	case TreeObject.String():
	case CommitObject.String():
	case TagObject.String():
	}
	return nil, nil
}
