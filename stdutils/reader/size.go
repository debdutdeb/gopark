package reader

import "io"

// ReadNer is a super-interface that returns the total byte count of what was read from it.
//
// NOTE(self): This is a bit weirder name. But I'm gonna keep it. Inspiration is from `n, err := r.Read()`, `n` is the convention, so `Read`+`N`+`er`.
type ReadNer interface {
	io.Reader
	ByteCount() int
}

type readNer struct {
	io.Reader
	n int
}

// Ner returns a ReadNer interface from an [io.Reader]
func Ner(reader io.Reader) ReadNer {
	return &readNer{Reader: reader, n: 0}
}

func (r *readNer) Read(buffer []byte) (n int, err error) {
	n, err = r.Reader.Read(buffer)
	r.n += n
	return n, err
}

func (r *readNer) ByteCount() int {
	return r.n
}

