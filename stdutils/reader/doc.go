/*
This package primarily expands on [io.Reader].


- [ErrorCloser]: It takes an [io.Reader] and returns an [io.ReadCloser] and an error channel

- [WrapReaderInContext]: WrapReaderInContext makes [io.Reader]'s Read receiver context aware. If the context expires, Read fails


Error closer is meant to be used in go routines that need to pipe data to a reader but the process might fail
Caller should call Close on the readcloser to check the error from the channel
Close should be called after it's done reading from the reader
If the error channel is not closed, Close will wait for 10 seconds and return a timeout error

# Using ErrorCloser

Example code

	// create a readcloser
	reader, errCh := ErrorCloser(reader)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// notify the piping process is done
		defer close(errCh)
		defer wg.Done()
		_, err := io.Copy(dst, reader)
		if err != nil {
			// something failed, propagate the error through the `Close()` receiver
			errCh <- err
		}
	}()

	// first wait for the go routine to finish, imitating reader consumption finished

	wg.Wait()

	// either pipe closing fails or piping itself fails
	err := reader.Close()

	if err != nil {
	  log.Fatal(err)
	}

Use [ErrorCloserWithTimeout] to pass a custom timeout (instead of default 10 seconds) before Close gives up on waiting for errCh.

Another helper is [WrapReaderInContext]. It takes an [io.Reader] and a [context.Context] and returns an [io.Reader] that checks the context before every Read.

Once the context is cancelled or its deadline is exceeded, Read returns the context's error instead of delegating to the underlying reader. This is useful to make a reader responsive to a caller's context, for example one built from a request, without the underlying reader supporting it natively.

# Using WrapReaderInContext

[WrapReaderInContext] makes an [io.Reader] context aware, however there is a caveat, if a goroutine is already blocked at the underlying `Read`, cancelling the context means nothing.

[WrapReaderInContext] only mostly helps in one case, where Read is happening in small chunks.

Ideally the underlying struct that implements Read should be context aware, and sharing the same context with [WrapReaderInContext] makes no sense. 

Example code

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := WrapReaderInContext(reader, ctx)

	// Read returns ctx.Err() once ctx is done, instead of blocking on the
	// underlying reader
	_, err := io.Copy(dst, reader)

	if err != nil {
	  log.Fatal(err)
	}

# Using io.Reader super-interfaces

These are composable. For example to create a reader to output both hash and bytecount read from it

	type fullReader struct {
		io.Reader
		ByteCount func() int
		Sum func() string
	}
	var r io.Reader
	f := &fullReader{
	}
	ner := Ner(r)
	f.ByteCount = func() int { return ner.ByteCount() }
	hasher := Sha256Hasher(ner)
	f.Sum = func() string { return hex.EncodeToString(hasher.Sum(nil)) }
	f.Reader = hasher

	// use f as the reader now
	var _ io.Reader = f
*/
package reader
