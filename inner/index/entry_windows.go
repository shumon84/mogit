//+build windows

package index

import (
	"os"
	"syscall"
	"time"
)

func NewEntry(fileInfo os.FileInfo, digest []byte) (*Entry, error) {
	if fileInfo == nil {
		return nil, ErrNilFileInfo
	}

	stat := fileInfo.Sys().(syscall.Win32FileAttributeData)

	entry := &Entry{
		CTime: time.Unix(
			int64(stat.CreationTime.HighDateTime),
			int64(stat.CreationTime.LowDateTime),
		),
		MTime: time.Unix(
			int64(stat.LastWriteTime.HighDateTime),
			int64(stat.LastWriteTime.LowDateTime),
		),
		Dev:           0,
		Ino:           0,
		ObjectType:    RegularFile,
		Permission:    uint16(fileInfo.Mode().Perm()),
		Size:          uint32(fileInfo.Size()),
		Digest:        digest,
		IsAssumeValid: false,
		ConflictFlag:  0,
		Name:          fileInfo.Name(),
	}

	return entry, nil
}
