// Package ops provides operators for iterators.
//
// All operators are guaranteed to:
//   - run in linear time
//   - allocate constant memory
//   - depend only on the iter and constraints packages
//   - not spawn additional goroutines
//
// Operators that cannot be implemented within these constraint will be added to
// a separate packages.
package ops

import "iter"

/***********
* Cropping *
************/

// Take emits the first n items of the source iterator.
func TakeN[T any](src iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(src)
		defer stop()
		for i := 0; i < n; i++ {
			t, ok := next()
			if !ok {
				return
			}
			if !yield(t) {
				return
			}
		}
	}
}

/***********************
* Plucking and packing *
************************/

// Keys emits the keys, or first items of every couple emitted by the source iterator.
func Keys[K, V any](src iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range src {
			if !yield(k) {
				return
			}
		}
	}
}

// Values emits the values, or second items of every couple emitted by the source iterator.
func Values[K, V any](src iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range src {
			if !yield(v) {
				return
			}
		}
	}
}

// Entries emits couples of values that represent the key-value pairs from the source iterator.
func Entries[K, V any](src iter.Seq2[K, V]) iter.Seq[struct {
	K K
	V V
}] {
	return func(yield func(struct {
		K K
		V V
	}) bool) {
		for k, v := range src {
			if !yield(struct {
				K K
				V V
			}{k, v}) {
				return
			}
		}
	}
}

/***************
* Transforming *
****************/

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

// Map2 is like [Map] for iter.Seq2.
func Map2[K1, V1, K2, V2 any](src iter.Seq2[K1, V1], predicate func(K1, V1) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range src {
			if k2, v2 := predicate(k1, v1); !yield(k2, v2) {
				return
			}
		}
	}
}

// Map12 is like [Map] but it transforms the iterator from Seq to Seq2.
func Map12[T, K, V any](src iter.Seq[T], predicate func(T) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for t := range src {
			if k, v := predicate(t); !yield(k, v) {
				return
			}
		}
	}
}

// Map21 is like [Map] but it transforms the iterator from Seq2 to Seq.
func Map21[K, V, T any](src iter.Seq2[K, V], predicate func(K, V) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k, v := range src {
			if t := predicate(k, v); !yield(t) {
				return
			}
		}
	}
}

// Filter emits the item that the predicate returns true for.
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

// Filter2 is like [Filter] for Seq2.
func Filter2[K, V any](src iter.Seq2[K, V], predicate func(K, V) (ok bool)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range src {
			if ok := predicate(k, v); !ok {
				continue
			}
			if !yield(k, v) {
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

// Zip emits every time both source iterators have emitted
// a value, thus generating couples of values where no source value is used more than
// once and no one is discarded except for the trailing ones after one of the sources
// has stopped generating values.
func Zip[T, V any](src1 iter.Seq[T], src2 iter.Seq[V]) iter.Seq2[T, V] {
	return func(yield func(T, V) bool) {
		next1, stop1 := iter.Pull(src1)
		defer stop1()
		next2, stop2 := iter.Pull(src2)
		defer stop2()
		for {
			t, ok1 := next1()
			if !ok1 {
				return
			}
			v, ok2 := next2()
			if !ok2 {
				return
			}
			if !yield(t, v) {
				return
			}
		}
	}
}

// Tap calls peek for all values emitted by src and consumed by the returned Seq.
//
// Peek must not modify or keep a reference to the values it observes.
func Tap[T any](src iter.Seq[T], peek func(T)) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range src {
			peek(t)
			if !yield(t) {
				return
			}
		}
	}
}

// Deduplicate removes duplicates emitted by src. It doesn't check that
// the entire iterator never emits two identical values, it just removes consecutive
// identical values.
func Deduplicate[T comparable](src iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(src)
		defer stop()
		prev, ok := next()
		if !ok || !yield(prev) {
			return
		}
		for {
			t, ok := next()
			if !ok {
				return
			}
			if t == prev {
				continue
			}
			prev = t
			if !yield(t) {
				return
			}
		}
	}
}

/***************
* Higher order *
****************/

// Flatten emits all values emitted by the inner iterators, flattening the source iterator
// structure to be one layer.
func Flatten[T any](src iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range src {
			for t := range i {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// FlattenSlice is like [Flatten] for iterators of slices.
func FlattenSlice[T any](src iter.Seq[[]T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range src {
			for _, t := range i {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// Flatten2 emits all values emitted by the inner iterators, flattening the source iterator
// structure to be one layer. Keys for inner iterators are repeated for every inner emission.
func Flatten2[K, V any](src iter.Seq2[K, iter.Seq[V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, i := range src {
			for v := range i {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Concat emits all values from the provided sources, in order.
func Concat[T any](srcs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, src := range srcs {
			for t := range src {
				if !yield(t) {
					return
				}
			}
		}
	}
}
