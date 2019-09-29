package index

import "os"

// ObjectType is a object type of git index entry.
type ObjectType uint

// constants of git index entry object type.
const (
	RegularFile  ObjectType = 0x8
	SymbolicLink ObjectType = 0xA
	GitLink      ObjectType = 0xE
)

// String is an implementation of fmt.Stringer interface.
func (objectType ObjectType) String() string {
	switch objectType {
	case RegularFile:
		return "regular"
	case SymbolicLink:
		return "symbolic link"
	case GitLink:
		return "git link"
	default:
		return "undefined"
	}
}

// GetObjectType ges object type from os.FileInfo
// Note! This function has not supported GitLink type yet.
func GetObjectType(fileInfo os.FileInfo) (ObjectType, error) {
	if fileInfo == nil {
		return 0, ErrNilFileInfo
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return SymbolicLink, nil
	}

	return RegularFile, nil
}
