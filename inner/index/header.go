package index

import (
	"fmt"
	"io"

	"github.com/shumon84/binutil"
)

// HeaderSignature is correct signature for git index file header.
// it is just 4byte and []byte{'D','I','R','C'}
const HeaderSignature = "DIRC"

// Header is struct of git index file header.
type Header struct {
	Signature    string
	Version      uint32
	NumOfEntries uint32
}

// String is implementation of fmt.Stringer interface
func (h *Header) String() string {
	return fmt.Sprintf(`Signature    : %s
Version      : %d
len(entries) : %d`, h.Signature, h.Version, h.NumOfEntries)
}

// ReadHeader reads git index file header.
// r of parameters must be byte stream of .git/index
func ReadHeader(r binutil.Reader) (*Header, error) {
	if r == nil {
		return nil, ErrNilReader
	}

	// initialize reading position
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	signature, err := readSignature(r)
	if err != nil {
		return nil, err
	}
	version, err := readVersion(r)
	if err != nil {
		return nil, err
	}
	numOfEntries, err := readNumOfEntries(r)
	if err != nil {
		return nil, err
	}

	return &Header{
		Signature:    signature,
		Version:      version,
		NumOfEntries: numOfEntries,
	}, nil
}

func readSignature(r binutil.Reader) (string, error) {
	signatureByte, err := r.Bytes(4)
	if err != nil {
		return "", err
	}
	signature := string(signatureByte)
	if signature != HeaderSignature {
		return "", ErrInvalidSignature
	}
	return signature, nil
}

func readVersion(r binutil.Reader) (uint32, error) {
	supportedVersions := map[uint32]struct{}{
		2: {},
	}
	version, err := r.UInt32()
	if err != nil {
		return 0, err
	}
	if _, ok := supportedVersions[version]; !ok {
		return 0, ErrNotSupportedVersion
	}
	return version, nil
}

func readNumOfEntries(r binutil.Reader) (uint32, error) {
	return r.UInt32()
}
