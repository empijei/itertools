package from

import (
	"bufio"
	"context"
	"iter"
)

// ScannerText emits all text emitted by s.
//
// Cancellation must be handled by closing the source reader.
//
// After ScannerText returns the caller is responsible for checking s.Err().
func ScannerText(s *bufio.Scanner) iter.Seq[string] {
	return func(yield func(string) bool) {
		for s.Scan() {
			if !yield(s.Text()) {
				return
			}
		}
	}
}

// Chan emits all values received on src and stops whenever src is closed or the context is cancelled.
func Chan[T any](ctx context.Context, src <-chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-src:
				if !ok {
					return
				}
				if !yield(t) {
					return
				}
			}
		}
	}
}
