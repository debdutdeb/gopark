package reader

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func buffern(n int) []byte {
	return make([]byte, n)
}

func buffer() []byte {
	return buffern(1)
}

func TestReadAt(t *testing.T) {
	wg := sync.WaitGroup{}

	reader := NopAt(strings.NewReader("abc"))

	mut := sync.Mutex{}

	res := buffern(3)

	// try second first
	wg.Add(1)
	go func() {
		defer wg.Done()
		b := buffer()
		_, err := reader.ReadAt(b, 1)
		if err != nil {
			t.Fatalf("%q at offset 1\n", err)
		}
		mut.Lock()
		res[1] = b[0]
		mut.Unlock()
	}()

	time.Sleep(100 * time.Millisecond)

	// try third next
	wg.Add(1)
	go func() {
		defer wg.Done()
		b := buffer()
		_, err := reader.ReadAt(b, 2)
		if err != nil {
			t.Fatalf("%q at offset 2\n", err)
		}
		mut.Lock()
		res[2] = b[0]
		mut.Unlock()
	}()
	time.Sleep(100 * time.Millisecond)

	// try first last
	wg.Add(1)
	go func() {
		defer wg.Done()
		b := buffer()
		_, err := reader.ReadAt(b, 0)
		if err != nil {
			t.Fatalf("%q at offset 0\n", err)
		}
		mut.Lock()
		res[0] = b[0]
		mut.Unlock()
	}()

	wg.Wait()

	if string(res) != "abc" {
		t.Fatalf("expected \"abc\" got \"%s\"\n", res)
	}
}
