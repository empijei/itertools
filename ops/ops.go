// Package ops provides operators for iterators.
package ops

import "iter"

/*
List of planned operators:
* TakeUntil
* Scan
* Contains
* Reduce
* Min
* Max
* Len
* Concat
* SkipN
* SkipUntil
* Unique
* Deduplicate
* Tap

List of planned constructors:
* from.ScannerText
* from.ScannerBytes
* from.Chan
*/

// Take forwards the first n items of the source iterator.
func TakeN[T any](src iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(src)
		defer stop()
		for i := 0; i < n; i++ {
			v, ok := next()
			if !ok {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Keys returns the keys, or first items of every couple emitted by the source iterator.
func Keys[K, V any](src iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range src {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns the values, or second items of every couple emitted by the source iterator.
func Values[K, V any](src iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range src {
			if !yield(v) {
				return
			}
		}
	}
}

// Map applies the predicate to the source iterator until either source is exhausted
// or the consumer stops the iteration.
func Map[T, V any](src iter.Seq[T], predicate func(T) V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for t := range src {
			if v := predicate(t); !yield(v) {
				return
			}
		}
	}
}

// Filter forwards the item that the predicate returns true for.
func Filter[T any](src iter.Seq[T], predicate func(T) (ok bool)) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range src {
			if ok := predicate(t); !ok {
				continue
			}
			if !yield(t) {
				return
			}
		}
	}
}

// PairWise emits all values with the value that preceded them.
// This means all values will be emitted twice except for the first and last one.
// Values are emitted once as the second value, then as the first, in this order.
// Pairs can be imagined as a sliding window on the source iterator.
func PairWise[T any](src iter.Seq[T]) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		next, stop := iter.Pull(src)
		defer stop()
		prev, ok := next()
		if !ok {
			return
		}
		for {
			cur, ok := next()
			if !ok {
				return
			}
			if !yield(prev, cur) {
				return
			}
			prev = cur
		}
	}
}

// Zip returns an iterator that emits every time both source iterators have emitted
// a value, thus generating couples of values where no source value is used more than
// once and no one is discarded except for the last ones.
func Zip[T, V any](src1 iter.Seq[T], src2 iter.Seq[V]) iter.Seq2[T, V] {
	return func(yield func(T, V) bool) {
		next1, stop1 := iter.Pull(src1)
		defer stop1()
		next2, stop2 := iter.Pull(src2)
		defer stop2()
		for {
			v1, ok1 := next1()
			if !ok1 {
				return
			}
			v2, ok2 := next2()
			if !ok2 {
				return
			}
			if !yield(v1, v2) {
				return
			}

		}

	}
}
