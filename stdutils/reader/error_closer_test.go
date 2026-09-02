package reader

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

func TestErrorCloser_PipeErrored_ClosePropagatesError(t *testing.T) {
	r := strings.NewReader("data")
	rc, errCh := ErrorCloser(r)

	sentErr := errors.New("pipe failed")
	go func() {
		errCh <- sentErr
		close(errCh)
	}()

	// Consume reader so producer can finish
	_, _ = io.ReadAll(rc)

	err := rc.Close()
	if err == nil {
		t.Fatal("Close() expected to return the propagated error, got nil")
	}
	if !errors.Is(err, sentErr) && err.Error() != sentErr.Error() {
		t.Errorf("Close() = %v, want %v", err, sentErr)
	}
}

func TestErrorCloser_ErrorChannelNeverClosed_ReturnsTimeoutError(t *testing.T) {
	// Use short timeout so test doesn't block 10s
	orig := defaultCloseTimeout
	defaultCloseTimeout = 100 * time.Millisecond
	t.Cleanup(func() { defaultCloseTimeout = orig })

	r := strings.NewReader("data")
	rc, errCh := ErrorCloser(r)

	// Never close errCh and never send; just leak the channel
	_ = errCh

	_, _ = io.ReadAll(rc)

	err := rc.Close()
	if err == nil {
		t.Fatal("Close() expected timeout error when channel is never closed, got nil")
	}
	if !strings.Contains(err.Error(), "close timeout") || !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Close() = %v, want error containing 'close timeout'", err)
	}
}

func TestErrorCloser_NoErrors_CloseReturnsNil(t *testing.T) {
	r := strings.NewReader("data")
	rc, errCh := ErrorCloser(r)

	go func() {
		close(errCh) // signal completion with no error
	}()

	_, _ = io.ReadAll(rc)

	err := rc.Close()
	if err != nil {
		t.Errorf("Close() = %v, want nil (no errors)", err)
	}
}
