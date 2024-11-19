# Itertools

Tools and libraries to deal with iterators in Go.

This aims to be a polished collection of tools to manipulate iterators in Go.

Iterators have been introduced in Go 1.23 and I feel like we'll soon get some helpers
and APIs in the standard library to create and mutate them.

If you are eager to use operations like `Map` or `Filter` or would like a more
idiomatic API for `bufio.Scanner` this module is probably for you.

I'll make sure that, as soon as standard alternatives become available, I'll
deprecate my versions and facilitate the migration to the new ones (potentially
providing tools to automatically rewrite code).

# Notes

It must be highlighted that I am not endorsing a programming style that encourages
mapreduce-like code and that pushes for a higher mental overhead than it's necessary.

This library is intended to help when those operations are actually needed to make code easier to read, but
I would invite the reader of this document to try and avoid creating obscure code bases
that rely too heavily on this library.

# Planned work

- [ ] Improve documentation with examples
- [ ] Improve this document
- [ ] Clearly state how to idiomatically use this package
- [ ] Stabilize API and bump to v1

## Operators (Package `ops`)

Cropping:

- [ ] SkipN
- [ ] SkipUntil
- [ ] TakeUntil

## Extra operators (Package `xops`)

Uniq:

- [ ] Unique
- [ ] Debounce
- [ ] CombineLatest

## Constructors (Package `from`)

- [ ] from.ScannerBytes

## Sinks (Package `to`)

## Harnesses (Package `itertest`)

- [ ] test utils to check for iterators termination
