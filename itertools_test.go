package itertools_test

import (
	"fmt"
	"iter"
	"slices"
	"strconv"
	"testing"

	. "github.com/empijei/itertools"
	"github.com/google/go-cmp/cmp"
)

// TODO(empijei): improve termination tests and harnesses.

func TestTermination11(t *testing.T) {
	t.Parallel()
	const (
		target = 10
		margin = 10
	)

	tests := []struct {
		name string
		it   func(iter.Seq[int]) iter.Seq[int]
	}{
		{"TakeN", func(src iter.Seq[int]) iter.Seq[int] {
			return TakeN(src, target+margin)
		}},
		{"Map", func(src iter.Seq[int]) iter.Seq[int] {
			return Map(src, func(i int) int { return i })
		}},
		{"Filter", func(src iter.Seq[int]) iter.Seq[int] {
			return Filter(src, func(_ int) bool { return true })
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reads := 0
			tapSource := func(yield func(int) bool) {
				for i := range target + margin {
					if !yield(i) {
						return
					}
					reads++
				}
			}

			writes := 0
			countYield := func(int) bool {
				writes++
				return writes < target+1
			}

			tt.it(tapSource)(countYield)

			if reads != target {
				t.Errorf("%v reads: got %v want %v", tt.name, reads, target)
			}
			if writes != target+1 {
				t.Errorf("%v writes: got %v want %v", tt.name, writes, target)
			}
		})
	}
}

func TestTermination12(t *testing.T) {
	t.Parallel()
	const (
		target = 10
		margin = 10
	)

	tests := []struct {
		name       string
		it         func(iter.Seq[int]) iter.Seq2[int, int]
		adjr, adjw int
	}{
		{name: "PairWise",
			it: func(src iter.Seq[int]) iter.Seq2[int, int] {
				return PairWise(src)
			},
			adjr: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reads := 0
			tapSource := func(yield func(int) bool) {
				for i := range target + margin {
					if !yield(i) {
						return
					}
					reads++
				}
			}

			writes := 0
			countYield := func(int, int) bool {
				writes++
				return writes < target+1
			}

			tt.it(tapSource)(countYield)

			if got := reads + tt.adjr; got != target {
				t.Errorf("%v reads: got %v want %v", tt.name, got, target)
			}
			if got := writes + tt.adjw; got != target+1 {
				t.Errorf("%v writes: got %v want %v", tt.name, got, target)
			}
		})
	}
}

func TestTakeN(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  []int
		n    int
		want []int
	}{
		{[]int{}, 3, nil},
		{nil, 0, nil},
		{[]int{1, 2, 3, 4}, 2, []int{1, 2}},
		{[]int{3, 7, 11}, 10, []int{3, 7, 11}},
		{[]int{3, 7, 11}, 0, nil},
	}
	for _, tt := range tests {
		projection := TakeN(slices.Values(tt.src), tt.n)
		got := slices.Collect(projection)
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("TakeN(%v, %v): got %v want %v diff:\n%v", tt.src, tt.n, got, tt.want, diff)
		}
	}
}

func TestTakeWhile(t *testing.T) {
	t.Parallel()
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	pred := func(i int) bool { return i < 3 }
	got := slices.Collect(TakeWhile(slices.Values(src), pred))
	want := []int{1, 2}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TakeWhile(1->9, <3): got %v want %v", got, want)
	}
}

func TestSkipN(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  []int
		n    int
		want []int
	}{
		{[]int{}, 3, nil},
		{nil, 0, nil},
		{[]int{1, 2, 3, 4}, 2, []int{3, 4}},
		{[]int{3, 7, 11}, 10, nil},
		{[]int{3, 7, 11}, 0, []int{3, 7, 11}},
	}
	for _, tt := range tests {
		projection := SkipN(slices.Values(tt.src), tt.n)
		got := slices.Collect(projection)
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("SkipN(%v, %v): got %v want %v diff:\n%v", tt.src, tt.n, got, tt.want, diff)
		}
	}
}

func TestSkipUntil(t *testing.T) {
	t.Parallel()
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var callCount int
	pred := func(i int) bool {
		callCount++
		return i > 7
	}
	got := slices.Collect(SkipUntil(slices.Values(src), pred))
	want := []int{8, 9}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("SkipUntil(1->9, >7): got %v want %v", got, want)
	}
	if wantCalls := 8; callCount != wantCalls {
		t.Errorf("SkipUntil(1->9, >7): got %v calls to predicate, want %v", callCount, wantCalls)
	}
}

func TestKeys(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  []int
		want []int
	}{
		{[]int{1, 2, 3, 4}, []int{0, 1, 2, 3}},
		{nil, nil},
	}
	for _, tt := range tests {
		projection := Keys(slices.All(tt.src))
		got := slices.Collect(projection)
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Keys(%v): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

func TestValues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  []int
		want []int
	}{
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{nil, nil},
	}
	for _, tt := range tests {
		projection := Values(slices.All(tt.src))
		got := slices.Collect(projection)
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Values(%v): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

type intPair = struct {
	K, V int
}

func TestEntries(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		src  iter.Seq2[int, int]
		want []intPair
	}{
		{
			slices.All([]int{2, 3, 4, 5}),
			[]intPair{{0, 2}, {1, 3}, {2, 4}, {3, 5}},
		},
	}
	for _, tt := range tests {
		got := slices.Collect(Entries(tt.src))
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Entries(slices.All(%v)): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

func TestMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  []int
		want []int
	}{
		{[]int{}, nil},
		{nil, nil},
		{[]int{1, 2}, []int{2, 4}},
		{[]int{3, 7, 11}, []int{6, 14, 22}},
	}
	times2 := func(i int) int { return i * 2 }

	for _, tt := range tests {
		projection := Map(slices.Values(tt.src), times2)
		got := slices.Collect(projection)
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Map(%v, times2): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  []int
		want []int
	}{
		{[]int{}, nil},
		{nil, nil},
		{[]int{1, 2, 3, 4}, []int{2, 4}},
		{[]int{3, 7, 11}, nil},
		{[]int{2, 4, 6}, []int{2, 4, 6}},
	}
	isEven := func(i int) bool { return i%2 == 0 }

	for _, tt := range tests {
		projection := Filter(slices.Values(tt.src), isEven)
		got := slices.Collect(projection)
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Filter(%v, isEven): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

func TestMapFilterHigherArity(t *testing.T) {
	t.Parallel()
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	withStrKeys := Map12(slices.Values(src), func(i int) (string, int) {
		return strconv.Itoa(i), i
	})
	doubled := Map2(withStrKeys, func(k string, v int) (string, int) {
		return k + "*2", v * 2
	})
	onlyMultiple3 := Filter2(doubled, func(_ string, v int) (ok bool) {
		return v%3 == 0
	})
	toString := Map21(onlyMultiple3, func(k string, v int) string {
		return fmt.Sprintf("%v = %v", k, v)
	})
	got := slices.Collect(toString)
	want := []string{"3*2 = 6", "6*2 = 12", "9*2 = 18"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("1->10 | *2 | %%3 == 0 | toString: got %v want %v diff:\n%v", got, want, diff)
	}
}

func TestPairWise(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  []int
		want [][]int
	}{
		{
			[]int{1, 2, 3, 4},
			[][]int{{1, 2}, {2, 3}, {3, 4}},
		},
		{
			[]int{1},
			nil,
		},
		{nil, nil},
	}

	for _, tt := range tests {
		srci := slices.Values(tt.src)
		var got [][]int
		for a, b := range PairWise(srci) {
			got = append(got, []int{a, b})
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Pairwise(%v): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

func TestZip(t *testing.T) {
	t.Parallel()
	tests := []struct {
		srca, srcb []int
		want       [][]int
	}{
		{
			[]int{1, 3, 5, 7},
			[]int{2, 4, 6, 8},
			[][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}},
		},
		{
			[]int{1, 3, 5, 7},
			[]int{2},
			[][]int{{1, 2}},
		},
		{
			[]int{1},
			[]int{2, 4, 6, 8},
			[][]int{{1, 2}},
		},
		{
			nil,
			[]int{2, 4, 6, 8},
			nil,
		},
	}

	for _, tt := range tests {
		srcai := slices.Values(tt.srca)
		srcbi := slices.Values(tt.srcb)
		var got [][]int
		for a, b := range Zip(srcai, srcbi) {
			got = append(got, []int{a, b})
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Zip(%v,%v): got %v want %v diff:\n%v", tt.srca, tt.srcb, got, tt.want, diff)
		}
	}
}

func TestTap(t *testing.T) {
	t.Parallel()
	t.Run("stops consuming", func(t *testing.T) {
		t.Parallel()
		var sum int
		it := Tap(slices.Values([]int{1, 2, 3, 4, 5}), func(i int) {
			sum += i
		})
		it(func(i int) bool {
			return i <= 3
		})
		if sum != 10 {
			t.Errorf("Tap([1 2 3 4 Cancelled], sumAll): got %v want 10", sum)
		}
	})
	t.Run("consumes all", func(t *testing.T) {
		t.Parallel()
		var sum int
		it := Tap(slices.Values([]int{1, 2, 3, 4, 5}), func(i int) {
			sum += i
		})
		it(func(_ int) bool {
			return true
		})
		if sum != 15 {
			t.Errorf("Tap([1 2 3 4 5], sumAll): got %v want 15", sum)
		}
	})
}

func TestDeduplicate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  []int
		want []int
	}{
		{
			[]int{1, 2, 3, 4, 4, 4, 5, 5, 6, 7},
			[]int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			[]int{1},
			[]int{1},
		},
		{
			[]int{1, 2, 3, 4, 4, 4, 5, 5, 1, 6, 7},
			[]int{1, 2, 3, 4, 5, 1, 6, 7},
		},
		{},
	}

	for _, tt := range tests {
		got := slices.Collect(Deduplicate(slices.Values(tt.src)))
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Deduplicate(%v): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

func TestFlatten(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  [][]int
		want []int
	}{
		{
			[][]int{{1, 2, 3}, {4, 5, 6}},
			[]int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		var in []iter.Seq[int]
		for _, i := range tt.src {
			in = append(in, slices.Values(i))
		}
		got := slices.Collect(Flatten(slices.Values(in)))
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Flatten(%v): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

func TestFlattenSlice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  [][]int
		want []int
	}{
		{
			[][]int{{1, 2, 3}, {4, 5, 6}},
			[]int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		got := slices.Collect(FlattenSlice(slices.Values(tt.src)))
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("FlattenSlice(%v): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

func TestFlatten2(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  [][]int
		want [][2]int
	}{
		{
			[][]int{{1, 2}, {3, 4}},
			[][2]int{{0, 1}, {0, 2}, {1, 3}, {1, 4}},
		},
	}
	for _, tt := range tests {
		it := Flatten2(func(yield func(int, iter.Seq[int]) bool) {
			for k, v := range tt.src {
				if !yield(k, slices.Values(v)) {
					return
				}
			}
		})
		var got [][2]int
		it(func(k, v int) bool {
			got = append(got, [2]int{k, v})
			return true
		})
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Flatten2(%v): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}

func TestConcat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		src  [][]int
		want []int
	}{
		{
			[][]int{{1, 2, 3}, {4, 5, 6}},
			[]int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		var in []iter.Seq[int]
		for _, i := range tt.src {
			in = append(in, slices.Values(i))
		}
		got := slices.Collect(Concat(in...))
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Concat(%v): got %v want %v diff:\n%v", tt.src, got, tt.want, diff)
		}
	}
}
