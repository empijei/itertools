package from

import (
	"bufio"
	"context"
	"errors"
	"io/fs"
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

// DirStep represents a step in a directory Walk.
type DirStep struct {
	// FullPath represents the path anchored to the root walk directory.
	// This mimics the behavior of "path" in fs.WalkDirFunc.
	FullPath string
	// Entry is the DirEntry that would be passed to fs.WalkDirFunc.
	Entry fs.DirEntry
}

// DirWalk emits all entries for root and its subdirectories.
// Errors are forwarded, and the consumer may decide wether to stop iteration or continue consuming further values.
//
// Use os.DirFS(path) to create fsys from disk.
func DirWalk(ctx context.Context, fsys fs.FS, root string) iter.Seq2[DirStep, error] {
	return func(yield func(DirStep, error) bool) {
		fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
			if !yield(DirStep{FullPath: path, Entry: d}, err) {
				return errors.New("consumer stopped")
			}
			return nil
		})
	}
}
