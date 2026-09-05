package reader

import (
	"fmt"
	"io"
	"runtime"
	"sync/atomic"
)

type noopReaderAt struct {
	reader  io.Reader
	offset  int
	control *atomic.Uint32
}

var ErrUnexpectedRegress = fmt.Errorf("unexpected regress of reader cursor")

// NopAt returns a reader that implements the ReadAt receiver.
//
// A no-op [io.ReaderAt] keeps track of its own position.
// A read attempt is blocked indefinitely unless the passed
// matches the internally tracked offset exactly.
// If a goroutine requests data starting at 1025th byte, the reader
// must have already exhausted the first 1024 bytes.
//
// If reader has exhausted more than 1024 bytes:
// ReadAt returns [ErrUnexpectedRegress].
//
// If requested offset does not match tracked offset (reader has not been read to the requested offset yet)
// the goroutine will be in a sleeping state.
func NopAt(reader io.Reader) io.ReaderAt {
	return &noopReaderAt{
		reader:  reader,
		control: &atomic.Uint32{},
	}
}

func (r *noopReaderAt) trylock() bool {
	return r.control.CompareAndSwap(0, 1)
}

func (r *noopReaderAt) unlock() {
	r.control.Store(0)
}

func (r *noopReaderAt) ReadAt(buffer []byte, offset int64) (n int, err error) {
	for {
		if !r.trylock() {
			runtime.Gosched()
			continue
		}
		trackedOffset := int64(r.offset) // tired of casting per line
		if offset == trackedOffset {
			// fine to advance
			n, err = r.reader.Read(buffer)
			r.offset += n
			r.unlock()
			return
		}

		if offset > trackedOffset {
			r.unlock()
			runtime.Gosched()
			continue
		}

		// if offset < trackedOffset
		r.unlock()
		return 0, ErrUnexpectedRegress
	}
}
