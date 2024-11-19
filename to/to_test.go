package to_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/empijei/itertools/to"
	"github.com/google/go-cmp/cmp"
)

func TestFirst(t *testing.T) {
	tests := []struct {
		src       []int
		want      int
		wantFound bool
	}{
		{[]int{1, 2, 3, 4}, 2, true},
		{[]int{1, 3, 5}, 0, false},
		{nil, 0, false},
	}

	for _, tt := range tests {
		got, found := to.First(slices.Values(tt.src), func(i int) bool { return i%2 == 0 })
		if got != tt.want {
			t.Errorf("First(%v, isEven): got value %v want %v", tt.src, got, tt.want)
		}
		if found != tt.wantFound {
			t.Errorf("First(%v, isEven): got found %v want %v", tt.src, found, tt.wantFound)
		}
	}
}

func TestContains(t *testing.T) {
	got := to.Contains(slices.Values([]int{1, 2, 3, 4, 5}), func(i int) bool { return i == 3 })
	if !got {
		t.Errorf("Contains(1->5, 3): got false want true")
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		src    []int
		want   int
		wantOk bool
	}{
		{[]int{8, 2, 3, 10}, 2, true},
		{nil, 0, false},
	}

	for _, tt := range tests {
		got, ok := to.Min(slices.Values(tt.src))
		if got != tt.want {
			t.Errorf("Min(%v): got value %v want %v", tt.src, got, tt.want)
		}
		if ok != tt.wantOk {
			t.Errorf("Min(%v): got ok %v want %v", tt.src, got, tt.want)
		}
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		src    []int
		want   int
		wantOk bool
	}{
		{[]int{8, 2, 3, 10, 7}, 10, true},
		{nil, 0, false},
	}

	for _, tt := range tests {
		got, ok := to.Max(slices.Values(tt.src))
		if got != tt.want {
			t.Errorf("Max(%v): got value %v want %v", tt.src, got, tt.want)
		}
		if ok != tt.wantOk {
			t.Errorf("Max(%v): got ok %v want %v", tt.src, got, tt.want)
		}
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		src  []int
		want int
	}{
		{[]int{8, 2, 3, 10, 7}, 5},
		{nil, 0},
	}

	for _, tt := range tests {
		got := to.Len(slices.Values(tt.src))
		if got != tt.want {
			t.Errorf("Len(%v): got %v want %v", tt.src, got, tt.want)
		}
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		src   []int
		start int
		want  int
	}{
		{[]int{1, 2, 3, 4}, 0, 10},
		{[]int{1, 2, 3, 4}, 7, 17},
		{nil, 10, 10},
	}

	for _, tt := range tests {
		got := to.Reduce(slices.Values(tt.src), tt.start, func(accum, cur int) (int, bool) {
			return accum + cur, true
		})
		if got != tt.want {
			t.Errorf("Reduce(%v, start=%v, sumAll): got value %v want %v", tt.src, tt.start, got, tt.want)
		}
	}
}

func TestChan(t *testing.T) {
	t.Run("values are emitted", func(t *testing.T) {
		ctx := context.Background()
		src := []int{1, 2, 3, 4}
		srci := slices.Values(src)
		var got []int
		for v := range to.Chan(ctx, srci, 0) {
			got = append(got, v)
		}
		if diff := cmp.Diff(src, got); diff != "" {
			t.Errorf("Chan(%v): got %v want %v diff:\n%v", src, got, src, diff)
		}
	})
	t.Run("cancellation is handled", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		src := []int{1, 2, 3, 4, 5}
		srci := slices.Values(src)
		var got []int

		c := to.Chan(ctx, srci, 0)
		v, ok := <-c
		got = append(got, v)
		cancel()
		for ok {
			v, ok = <-c
			if ok {
				got = append(got, v)
			}
		}

		if len(got) > 2 {
			t.Errorf("Chan(%v): got %v want less than 2 values", src, got)
		}
	})
}
