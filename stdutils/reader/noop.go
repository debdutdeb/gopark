package reader

import "io"

type noopReader struct {
	io.Reader
}

var (
	_ io.Reader  = &noopReader{}
	_ ReadHasher = &noopReader{}
)

// NopHasher returns a ReadHasher with No-op Sum implementation.
func NopHasher(reader io.Reader) ReadHasher {
	return &noopReader{Reader: reader}
}

func (r *noopReader) Sum([]byte) []byte {
	return nil
}
