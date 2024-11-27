// Package meta is experimental and tries to create an iteration API that allows
// for piping and composition by offering transformation constructors.
package meta

import (
	"iter"

	"github.com/empijei/itertools"
)

// Map returns a function that applies [itertools.Map] to the source iterator.
func Map[T, V any](predicate func(T) V) func(iter.Seq[T]) iter.Seq[V] {
	return func(src iter.Seq[T]) iter.Seq[V] {
		return itertools.Map(src, predicate)
	}
}

// Filter returns a function that applies [itertools.Filter] to the source iterator.
func Filter[T any](predicate func(T) bool) func(iter.Seq[T]) iter.Seq[T] {
	return func(src iter.Seq[T]) iter.Seq[T] {
		return itertools.Filter(src, predicate)
	}
}

// I almost had a stroke writing the signature for this function. I don't think
// this is very Go-like and the benefits composition provides are dwarfed by the
// added complexity.

// Combine combines two iterators transformations into one.
func Combine[T, I, V any](
	a func(iter.Seq[T]) iter.Seq[I],
	b func(iter.Seq[I]) iter.Seq[V],
) func(iter.Seq[T]) iter.Seq[V] {

	return func(s iter.Seq[T]) iter.Seq[V] {
		return b(a(s))
	}
}
