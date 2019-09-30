// util is a utility package to handle git repository
package util

import (
	"io"
	"io/ioutil"
	"path/filepath"
)

// FindGitRoot returns path to top level directory of current git repository
func FindGitRoot(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.IsDir() && file.Name() == ".git" {
			return dir, nil
		}
	}
	return FindGitRoot(filepath.Join(dir, ".."))
}

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

type WriteSeekCloser interface {
	io.Writer
	io.Seeker
	io.Closer
}

type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}
