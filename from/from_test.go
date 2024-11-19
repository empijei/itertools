package from_test

import (
	"bufio"
	"context"
	"slices"
	"strings"
	"testing"

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
