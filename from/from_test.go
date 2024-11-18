package from_test

import (
	"bufio"
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
