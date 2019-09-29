package index

// ConflictFlag is a type representing git index entry conflict flag
type ConflictFlag uint

// constants of git index entry conflict flag
const (
	NoConflict                 ConflictFlag = iota // 00
	LowestCommonAncestorCommit                     // 01
	CurrentCommit                                  // 10
	AnotherCommit                                  // 11
)

// String is an implementation of fmt.Stringer interface
func (conflictFlag ConflictFlag) String() string {
	switch conflictFlag {
	case NoConflict:
		return "no conflict"
	case LowestCommonAncestorCommit:
		return "lowest common ancestor commit"
	case CurrentCommit:
		return "current commit"
	case AnotherCommit:
		return "another commit"
	default:
		return "undefined"
	}
}
