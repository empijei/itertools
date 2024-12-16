package from_test

import (
	"bufio"
	"context"
	"io/fs"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/empijei/itertools/from"
	"github.com/google/go-cmp/cmp"
)

func TestScannerText(t *testing.T) {
	src := `Foo
Bar

After Empty
stop
last line`
	s := bufio.NewScanner(strings.NewReader(src))
	var got []string
	for s := range from.ScannerText(s) {
		if s == "stop" {
			break
		}
		got = append(got, s)
	}
	want := []string{"Foo", "Bar", "", "After Empty"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ScannerText(%q): got %v want %v diff:\n%v", src, got, want, diff)
	}
}

func TestChan(t *testing.T) {
	t.Run("values are emitted", func(t *testing.T) {
		src := []int{1, 2, 3, 4}
		srcc := make(chan int)
		go func() {
			for _, v := range src {
				srcc <- v
			}
			close(srcc)
		}()
		got := slices.Collect(from.Chan(context.Background(), srcc))
		if diff := cmp.Diff(src, got); diff != "" {
			t.Errorf("Chan(%v): got %v want %v diff:\n%v", src, got, src, diff)
		}
	})
	t.Run("cancellation is handled", func(t *testing.T) {
		srcc := make(chan int)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			srcc <- 1
			srcc <- 2
			cancel()
		}()
		got := slices.Collect(from.Chan(ctx, srcc))
		want := []int{1, 2}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Chan(1 2 CANCELLED): got %v want %v diff:\n%v", got, want, diff)
		}
	})
}

func TestDirWalk(t *testing.T) {
	fsys := fstest.MapFS(map[string]*fstest.MapFile{
		"root/foo/bar.txt": {Data: []byte("hello bar")},
		"root/empty":       {Mode: fs.ModeDir},
		"root/cat.txt":     {Data: []byte("hello cat")},
	})

	var errs []error
	var dirs []string
	from.DirWalk(context.Background(), fsys, "root")(func(ds from.DirStep, err error) bool {
		if err != nil {
			errs = append(errs, err)
		}
		dirs = append(dirs, ds.Path)
		return true
	})
	want := []string{"root", "root/cat.txt", "root/empty", "root/foo", "root/foo/bar.txt"}
	if diff := cmp.Diff(want, dirs); diff != "" {
		t.Errorf("DirWalk: got %v want %v diff:\n%s", dirs, want, diff)
	}
	if len(errs) != 0 {
		t.Errorf("DirWalk errors: got %v want none", errs)
	}
}
