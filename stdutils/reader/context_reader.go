package reader

import (
	"context"
	"io"
)

type contextReader struct {
	io.Reader
	ctx context.Context
}

// WrapReaderInContext makes [io.Reader]'s Read receiver context aware. If the context expires, Read fails.
func WrapReaderInContext(reader io.Reader, ctx context.Context) io.Reader {
	return &contextReader{
		Reader: reader,
		ctx:    ctx,
	}
}

func (r *contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}

	return r.Reader.Read(p)
}
