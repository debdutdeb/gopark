package reader

import "io"

type noopReader struct {
	io.Reader
}

// NopHasher returns a ReadHasher with No-op Sum implementation.
func NopHasher(reader io.Reader) ReadHasher {
	return &noopReader{Reader: reader}
}

func NopNer(reader io.Reader) ReadNer {
	return &noopReader{Reader: reader}
}

func (r *noopReader) Sum([]byte) []byte {
	return nil
}

func (r *noopReader) ByteCount() int { return -1 }
