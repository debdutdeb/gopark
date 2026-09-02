/*
This package primarily expands on [io.Reader].

The primary helper is [ErrorCloser]. It takes an [io.Reader] and returns an [io.ReadCloser] and an error channel.

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
*/
package reader
