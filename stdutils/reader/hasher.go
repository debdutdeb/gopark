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
	hash.Hash
}

type readHasher struct {
	io.Reader

	hash.Hash
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
		Hash:   h,
	}
}

// Write is a noop, ReadHasher should not be directly written to.
func (r *readHasher) Write(_ []byte) (n int, err error) {
	return 0, nil
}

// Reset is a noop, ReadHasher should not be manually reset.
func (r *readHasher) Reset() {}
