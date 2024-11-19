// Package to provides utilities to deconstruct or consume iterators down to other types or values.
package to

import (
	"context"
	"iter"

	"golang.org/x/exp/constraints"
)

func zero[T any]() (zero T) { return }

// First returns the first value for which predicate returns true and stops consuming src.
func First[T any](src iter.Seq[T], predicate func(T) bool) (t T, found bool) {
	for t := range src {
		if predicate(t) {
			return t, true
		}
	}
	return zero[T](), false
}

// Contains reports whether there is at least one value in the source iterator for which
// predicate returns true.
// It stops consuming the source at the first match.
func Contains[T any](src iter.Seq[T], predicate func(T) bool) bool {
	_, found := First(src, predicate)
	return found
}

// Min returns the minimum element emitted by the source and reports whether at
// least one value was consumed.
func Min[T constraints.Ordered](src iter.Seq[T]) (_ T, ok bool) {
	var m T
	var init bool
	for t := range src {
		if !init {
			init = true
			m = t
			continue
		}
		m = min(m, t)
	}
	return m, init
}

// Max returns the maximum element emitted by the source and reports whether at
// least one value was consumed.
func Max[T constraints.Ordered](src iter.Seq[T]) (_ T, ok bool) {
	var m T
	var init bool
	for t := range src {
		if !init {
			init = true
			m = t
			continue
		}
		m = max(m, t)
	}
	return m, init
}

// Len consumes the entire source and reports how many values it consumed.
func Len[T any](src iter.Seq[T]) int {
	var c int
	for range src {
		c++
	}
	return c
}

// Reduce scans the source and applies predicate on all elements it consumes until
// the predicate returns false or the source is exhausted.
//
// It passes subsequent accumulator values to each call of the predicate and returns the last value for it.
//
// The returned value may be the value for which predicate returned false or the
// last one before the source was exhausted.
func Reduce[T any](src iter.Seq[T],
	startAccum T,
	predicate func(accum, current T) (newAccum T, ok bool)) (lastAccum T) {
	accum := startAccum
	for t := range src {
		var ok bool
		accum, ok = predicate(accum, t)
		if !ok {
			return accum
		}
	}
	return accum
}

// Chan spawns a goroutine that consumes values emitted by the source and sends them
// on the returned channel.
// The channel is created with the provided buf size.
//
// Important note: the context cancellation is used to detect when to stop sending data
// on the returned channel and stop consuming the source, but since the iter.Seq
// API doesn't support cancellation, the goroutine spawned by Chan will only return
// once a yield call is performed by the source.
// Users of Chan should make sure that the source iterator stops when the related
// context is done.
func Chan[T any](ctx context.Context, src iter.Seq[T], buf int) <-chan T {
	c := make(chan T, buf)
	go func() {
		defer close(c)
		for t := range src {
			// Make sure we stop as soon as possible.
			select {
			case <-ctx.Done():
			default:
			}
			// Actually try to send the value
			select {
			case <-ctx.Done():
				return
			case c <- t:
			}
		}
	}()
	return c
}
