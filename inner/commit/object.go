package commit

import (
	"io"
	"os"
	"path/filepath"

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
		return "commit"
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

func ReadObject(digest string) (Object, error) {
	if len(digest) != 20 {
		return nil, ErrInvalidHashLength
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path, err := util.FindGitRoot(wd)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filepath.Join(path, digest[:2], digest[2:]))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// TODO: read git object
	return nil, nil
}
