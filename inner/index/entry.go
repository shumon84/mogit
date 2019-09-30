package index

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/shumon84/binutil"
)

// Entry is a type representing one of the files what is included git index tree.
type Entry struct {
	CTime         time.Time    // last change time of the file what this entry specifies
	MTime         time.Time    // last modification time of the file what this entry specifies
	Dev           int32        // Device ID of device existing the file what this entry specifies
	Ino           uint64       // inode of the file what this entry specifies
	ObjectType    ObjectType   // object type of the file what this entry specifies
	Permission    uint16       // permission of the file what this entry specifies
	UserID        uint32       // file owner's user ID
	GroupID       uint32       // file owner's group ID
	Size          uint32       // size of the file what this entry specifies
	Digest        []byte       // SHA1 digest of the file what this entry specifies
	IsAssumeValid bool         // bool value of whether this entry specifying file is assume valid
	ConflictFlag  ConflictFlag // flag to handle a conflicting file
	Name          string       // name of the file what this entry specifies
}

// String is implementation of fmt.Stringer interface
func (e *Entry) String() string {
	return fmt.Sprintf(`%s
  CTime       : %s
  MTime       : %s
  DeviceID    : %d
  Inode       : %d
  ObjectType  : %s
  Permission  : %o
  UserID      : %d
  GroupID     : %d
  FileSize    : %d
  SHA1        : %s
  AssumeValid : %v
  Conflict    : %d`,
		e.Name,
		e.CTime,
		e.MTime,
		e.Dev,
		e.Ino,
		e.ObjectType,
		e.Permission,
		e.UserID,
		e.GroupID,
		e.Size,
		hex.EncodeToString(e.Digest),
		e.IsAssumeValid,
		e.ConflictFlag)
}

// ReadEntries reads git index file entries.
// r of parameters must be byte stream of .git/index
func ReadEntries(r binutil.Reader, numOfEntries uint32) ([]*Entry, error) {
	if r == nil {
		return nil, ErrNilReader
	}

	// skip header section and jump to entries section
	if _, err := r.Seek(12, io.SeekStart); err != nil {
		return nil, err
	}

	entries := make([]*Entry, numOfEntries)
	for i := range entries {
		entry, err := readEntry(r)
		if err != nil {
			return nil, err
		}
		entries[i] = entry

		if err := seekToNextEntry(r); err != nil {
			return nil, err
		}
	}
	return entries, nil
}

func readEntry(r binutil.Reader) (*Entry, error) {
	ctime, err := readTime(r)
	if err != nil {
		return nil, err
	}
	mtime, err := readTime(r)
	if err != nil {
		return nil, err
	}
	dev, err := readDev(r)
	if err != nil {
		return nil, err
	}
	ino, err := readIno(r)
	if err != nil {
		return nil, err
	}
	permission, objectType, err := readMode(r)
	if err != nil {
		return nil, err
	}
	userID, err := readUserID(r)
	if err != nil {
		return nil, err
	}
	groupID, err := readGroupID(r)
	if err != nil {
		return nil, err
	}
	size, err := readSize(r)
	if err != nil {
		return nil, err
	}
	digest, err := readDigest(r)
	if err != nil {
		return nil, err
	}
	isAssumeValid, conflictFlag, fileNameLength, err := readFlags(r)
	if err != nil {
		return nil, err
	}
	name, err := readName(r, fileNameLength)
	if err != nil {
		return nil, err
	}

	return &Entry{
		CTime:         ctime,
		MTime:         mtime,
		Dev:           dev,
		Ino:           ino,
		ObjectType:    objectType,
		Permission:    permission,
		UserID:        userID,
		GroupID:       groupID,
		Size:          size,
		Digest:        digest,
		IsAssumeValid: isAssumeValid,
		ConflictFlag:  conflictFlag,
		Name:          name,
	}, nil
}

func readTime(r binutil.Reader) (time.Time, error) {
	sec, err := r.UInt32()
	if err != nil {
		return time.Time{}, err
	}
	nano, err := r.UInt32()
	if err != nil {
		return time.Time{}, err
	}
	log.Println(sec, ":", nano)
	return time.Unix(int64(sec), int64(nano)), nil
}

func readDev(r binutil.Reader) (int32, error) {
	dev, err := r.UInt32()
	if err != nil {
		return 0, err
	}
	return int32(dev), nil
}

func readIno(r binutil.Reader) (uint64, error) {
	dev, err := r.UInt32()
	if err != nil {
		return 0, err
	}
	return uint64(dev), nil
}

func readMode(r binutil.Reader) (uint16, ObjectType, error) {
	mode, err := r.UInt32()
	if err != nil {
		return 0, 0, err
	}

	// split mode(32 bit) to permission(0 ~ 9 bit) and object-type(12 ~ 16 bit).
	// -----------------------------------------------------------------
	// |                             mode                              |
	// |===============================================================|
	// | 00 01 02 03 04 05 06 07 08 | 09 0A 0B | 0C 0D 0E 0F | 10 - 1F |
	// |    unix like permission    |  unused  | object type | unused  |
	// -----------------------------------------------------------------
	permission := uint16(mode & 0x1FF)
	objectType := ObjectType((mode >> 12) & 0xF)

	return permission, objectType, nil
}

func readUserID(r binutil.Reader) (uint32, error) {
	return r.UInt32()
}

func readGroupID(r binutil.Reader) (uint32, error) {
	return r.UInt32()
}

func readSize(r binutil.Reader) (uint32, error) {
	return r.UInt32()
}

func readDigest(r binutil.Reader) ([]byte, error) {
	return r.Bytes(20)
}

func readFlags(r binutil.Reader) (bool, ConflictFlag, int, error) {
	flags, err := r.UInt16()
	if err != nil {
		return false, 0, 0, err
	}

	// split flags(16 bit) to assume-valid flag(0 bit), conflict flag(3 ~ 4 bit) and file name length(5 ~ 16 bit).
	// ------------------------------------------------------------
	// |                          flags                           |
	// |==========================================================|
	// |   0    ~    B    |    C     D    | E |         F         |
	// | file name length | conflict flag | 0 | assume valid flag |
	// ------------------------------------------------------------
	fileNameLength := int(flags & 0xFFF)
	conflictFlag := ConflictFlag((flags >> 12) & 0x3)
	isAssumeValid := ((flags >> 15) & 0x1) == 1

	return isAssumeValid, conflictFlag, fileNameLength, nil
}

func readName(r binutil.Reader, fileNameLength int) (string, error) {
	fileNameByte, err := r.Bytes(fileNameLength)
	if err != nil {
		return "", err
	}
	return string(fileNameByte), nil
}

func seekToNextEntry(r binutil.Reader) error {
	currentPosition, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	nextMultipleOf8Position := 8 - (currentPosition-12)&0x7
	if _, err := r.Seek(nextMultipleOf8Position, io.SeekCurrent); err != nil {
		if err == io.ErrUnexpectedEOF {
			return nil
		}
		return err
	}
	return nil
}
