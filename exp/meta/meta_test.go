package meta_test

import (
	"slices"
	"testing"

	"github.com/empijei/itertools/exp/meta"
	"github.com/google/go-cmp/cmp"
)

func TestCombineMapFilter(t *testing.T) {
	cmb := meta.Combine(
		meta.Map(func(i int) int {
			return i * 2
		}),
		meta.Filter(func(i int) bool {
			return i%3 == 0
		}),
	)

	got := slices.Collect(cmb(slices.Values([]int{1, 2, 3, 4, 5, 6})))
	want := []int{6, 12}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Combine(Map(*2), Filter(%%3==0))(1->6): got %v want %v diff:\n%v", got, want, diff)
	}
}
