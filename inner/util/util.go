// util is a utility package to handle git repository
package util

import (
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
