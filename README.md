# Itertools [![GoDoc](https://pkg.go.dev/badge/github.com/empijei/itertools)](https://pkg.go.dev/github.com/empijei/itertools) [![Go Report Card](https://goreportcard.com/badge/github.com/empijei/itertools)](https://goreportcard.com/report/github.com/empijei/itertools) [![Go build and tests](https://github.com/empijei/itertools/actions/workflows/go.yml/badge.svg)](https://github.com/empijei/itertools/actions/workflows/go.yml)

Tools and libraries to deal with iterators in Go > 1.23.

This aims to be a polished collection of tools to manipulate iterators in Go.

Iterators [have been introduced in Go 1.23](https://go.dev/blog/range-functions) and we'll likely soon get some helpers
and APIs in the standard library to create and mutate them.

If you are eager to use operations like `Map` or `Filter` or would like a more
idiomatic API for `bufio.Scanner` this module is probably for you.

I'll make sure that, as soon as standard alternatives become available, I'll
deprecate my versions and facilitate the migration to the new ones (e.g. providing
tools to automatically rewrite code).

# Subpackages

If you need to construct or consume iterators please use the [from](https://pkg.go.dev/github.com/empijei/itertools/from) and [to](https://pkg.go.dev/github.com/empijei/itertools/to) subpackages.

# Notes

I am not endorsing a programming style that encourages mapreduce-like code and
that pushes for a higher mental overhead than it's necessary.

This library is intended to help when those operations are actually needed to make
code easier to read, but I would invite the reader of this document to try and
avoid creating obscure code bases that rely too heavily on this library.
