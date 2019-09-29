//+build !windows

package index

import (
	"os"
	"syscall"
	"time"
)

// NewEntry creates a new git index entry made by file info.
// don't use (*File).Stat() to get os.FileInfo.
// if it is symbolic link, NewEntry will not work correctly.
// must use os.Lstat() to get os.FileInfo.
func NewEntry(fileInfo os.FileInfo, digest []byte) (*Entry, error) {
	if fileInfo == nil {
		return nil, ErrNilFileInfo
	}
	stat := fileInfo.Sys().(*syscall.Stat_t)

	objectType, err := GetObjectType(fileInfo)
	if err != nil {
		return nil, err
	}
	perm := uint16(fileInfo.Mode().Perm())

	// check valid combination of permission and object.
	// the allowed combinations of permission and object are as follows:
	// - RegularFile : 644
	// - RegularFile : 755
	// - Other type  : 000
	if objectType == RegularFile {
		if perm != 0644 && perm != 0755 {
			return nil, ErrForbiddenPermission
		}
	} else {
		if perm != 000 {
			return nil, ErrForbiddenPermission
		}
	}

	entry := &Entry{
		CTime:         time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec),
		MTime:         time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec),
		Dev:           stat.Dev,
		Ino:           stat.Ino,
		ObjectType:    objectType,
		Permission:    perm,
		UserID:        stat.Uid,
		GroupID:       stat.Gid,
		Size:          uint32(fileInfo.Size()),
		Digest:        digest,
		IsAssumeValid: false,
		ConflictFlag:  0,
		Name:          fileInfo.Name(),
	}
	return entry, nil
}
