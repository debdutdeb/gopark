package reader

import (
	"crypto/sha256"
	"hash"
	"io"
)

// ReadHasher implements both [io.Reader] and [hash.Hash] interfaces.
//
// This allows any ReadHasher to provide the underlying hash of the data that's been read.
type ReadHasher interface {
	io.Reader
	// Sum returns the hash of the read bytes.
	Sum(b []byte) []byte
}

type readHasher struct {
	io.Reader

	hash hash.Hash
}

// Sha256Hasher returns a [ReadHasher] that can be used to get the sha256 hash of
// the underlying bytes that's been read.
//
// Example:
//
//	reader := readerutils.Sha256Hasher(f)
//	_ = io.ReadAll(reader)
//	fmt.Println(hex.EncodeToString(reader.Sum(nil)))
func Sha256Hasher(reader io.Reader) ReadHasher {
	h := sha256.New()
	return &readHasher{
		Reader: io.TeeReader(reader, h),
		hash:   h,
	}
}

// Sum returns the hash of the read bytes.
func (r *readHasher) Sum(b []byte) []byte {
	return r.hash.Sum(b)
}
