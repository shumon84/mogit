// index is a package to handle git index file
//
// Note! Now supported version is 2 only.
//
// The following show git index file format overview.
//
//  -------------------------------------------------------------------------------
//  |  Header  |  4byte - signature('D','I','R','C')
//  |          |===================================================================
//  |  12byte  |  4byte - version number(2 or 3 or 4)
//  |          |===================================================================
//  |          |  4byte - number of index entries
//  -------------------------------------------------------------------------------
//  | Entry  0 |  4byte - ctime second
//  |          |===================================================================
//  |          |  4byte - ctime nano second
//  |          |===================================================================
//  |          |  4byte - mtime second
//  |          |===================================================================
//  |          |  4byte - mtime nano second
//  |          |===================================================================
//  |          |  4byte - device ID
//  |          |===================================================================
//  |          |  4byte - inode number
//  |          |===================================================================
//  |          |  4byte - mode  | 16bit - unused
//  |          |                |==================================================
//  |          |                |  4bit - object type
//  |          |                |         0x8 = regular file
//  |          |                |         0xA = symbolic link
//  |          |                |         0xE = git link
//  |          |                |==================================================
//  |          |                |  3bit - unused
//  |          |                |==================================================
//  |          |                |  9bit - unix like file permission
//  |          |===================================================================
//  |          |  4byte - file owner's user ID
//  |          |===================================================================
//  |          |  4byte - file owner's group ID
//  |          |===================================================================
//  |          |  4byte - file size
//  |          |===================================================================
//  |          | 20byte - SHA-1 digest
//  |          |===================================================================
//  |          |  2byte - flags |  1bit - assume-valid flag
//  |          |                |         see bellow a command more about this flag
//  |          |                |         $ git update-index --assume-unchanged
//  |          |                |==================================================
//  |          |                |  1bit - must be zero
//  |          |                |==================================================
//  |          |                |  2bit - conflict flag
//  |          |                |         00 = no conflict
//  |          |                |         01 = lowest common ancestor object's file
//  |          |                |         10 = current object's file
//  |          |                |         11 = another object's file
//  |          |                |==================================================
//  |          |                | 12bit - file name length
//  |          |===================================================================
//  |          | $(file name length)byte - relative path from top level directory
//  |          |===================================================================
//  |          | Fill in null bytes until file offset of next multiple of 8
//  -------------------------------------------------------------------------------
//  | Entry  1 |
//  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//  * All binary numbers are in network byte order.
//
// If you want to know more about git index file format, please refer to
// https://github.com/git/git/blob/master/Documentation/technical/index-format.txt
package index

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shumon84/mogit/inner/util"

	"github.com/shumon84/binutil"
)

// Index is interface of handle git index file
type Index interface {
	fmt.Stringer
	Header() *Header                    // get git index file header.
	Entries(idx uint32) (*Entry, error) // get idx-th git index entry.
}

type indexImpl struct {
	header  *Header
	entries []*Entry
}

// ReadIndex gets index tree from current repository
func ReadIndex() (Index, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path, err := util.FindGitRoot(currentDir)
	if err != nil {
		return nil, err
	}
	indexFile, err := os.Open(filepath.Join(path, ".git", "index"))
	if err != nil {
		return nil, err
	}
	return ReadIndexFromReader(indexFile)
}

// ReadIndexFromReade reads index tree from byte stream of .git/index
func ReadIndexFromReader(rs io.ReadSeeker) (Index, error) {
	if rs == nil {
		return nil, ErrNilReader
	}
	r := binutil.NewReader(rs)

	header, err := ReadHeader(r)
	if err != nil {
		return nil, err
	}

	entries, err := ReadEntries(r, header.NumOfEntries)
	if err != nil {
		return nil, err
	}

	return &indexImpl{
		header:  header,
		entries: entries,
	}, nil
}

// Header returns *Header in this index tree
func (i *indexImpl) Header() *Header {
	return i.header
}

// Entries returns idx-th entry in this index tree
func (i *indexImpl) Entries(idx uint32) (entry *Entry, err error) {
	if i.Header().NumOfEntries <= idx {
		return nil, ErrIndexOutOfEntryRanges
	}
	return i.entries[idx], nil
}

// String is implementation of fmt.Stringer interface
func (i *indexImpl) String() string {
	str := i.Header().String()
	for i, entry := range i.entries {
		str += fmt.Sprint("\n[Entry ", i, "] ", entry)
	}
	return str
}
