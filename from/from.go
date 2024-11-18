package from

import (
	"bufio"
	"iter"
)

func ScannerText(s *bufio.Scanner) iter.Seq[string] {
	return func(yield func(string) bool) {
		for s.Scan() {
			if !yield(s.Text()) {
				return
			}
		}
	}
}
