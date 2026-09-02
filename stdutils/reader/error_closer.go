package reader

import (
	"context"
	"fmt"
	"io"
	"time"
)

var defaultCloseTimeout = 10 * time.Second

type errorCloser struct {
	io.Reader
	errCh chan error

	timeout time.Duration
}

// ErrorCloser wraps a reader and returns a ReadCloser and an error channel
// Error closer is meant to be used in go routines that need to pipe data to a reader but the process might fail
// Caller should call Close on the readcloser to check the error from the channel
// Close should be called after it's done reading from the reader
// If the error channel is not closed, Close will wait for 10 seconds and return a timeout error
// Use ErrorCloserWithTimeout if timeout needs to be specified
// Example:
// 
//  reader, errCh := ErrorCloser(reader)
//  wg := sync.WaitGroup{}
//  wg.Add(1)
// 
// 	go func() {
// 		defer close(errCh)
// 		defer wg.Done()
// 		_, err := io.Copy(dst, reader)
// 		if err != nil {
// 			errCh <- err
// 		}
// 	}()
// 
//  // first wait for the go routine to finish, imitating reader consumption finished
// 
// 	wg.Wait()
// 
//  err := reader.Close()
// 
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
func ErrorCloser(reader io.Reader) (io.ReadCloser, chan<- error) {
	errCh := make(chan error, 1)
	return &errorCloser{
		Reader:  reader,
		errCh:   errCh,
		timeout: defaultCloseTimeout,
	}, errCh
}

// ErrorCloserWithTimeout returns a ReadCloser with associated error channel, and accepts a custom timeout for the error channel
// When Close() is called, expectation is the reader consumption has finished. A custom timeout allows the piping goroutine to have enough time to propagate an error or finish piping.
func ErrorCloserWithTimeout(reader io.Reader, timeout time.Duration) (io.ReadCloser, chan<- error) {
	ec, errCh := ErrorCloser(reader)
	ec.(*errorCloser).timeout = timeout
	return ec, errCh
}

func (e *errorCloser) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultCloseTimeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return fmt.Errorf("close timeout: producer failed to signal completion within %v, error channel should be closed to indicate completion: %w", defaultCloseTimeout, ctx.Err())

	case err, ok := <-e.errCh:
		if !ok {
			return nil
		}
		return err
	}
}
