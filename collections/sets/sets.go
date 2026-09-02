package sets

import (
	"iter"
	"maps"
)

// Set is a convinience re-type for a unique element list.
// Under the hood it is a map, therefore beyond the package level functions, all `maps` and `iter` package functions get carried over.
type Set[T comparable] map[T]struct{}

// New creates a new Set, with optional list of elements. It is a variadic function.
func New[T comparable](elements ...T) Set[T] {
	set := make(Set[T], len(elements))
	for _, e := range elements {
		set[e] = struct{}{}
	}
	return set
}

// Add adds eleements to the Set, mutates the Set and returns the same Set if want to chain
func Add[T comparable](set Set[T], elements ...T) Set[T] {
	for _, e := range elements {
		set[e] = struct{}{}
	}
	return set
}

// Has checks if an element exists in the Set or not.
func Has[T comparable](set Set[T], key T) bool {
	_, exists := set[key]
	return exists
}

// Collect creates a new Set from an iterator
func Collect[T comparable](seq iter.Seq[T]) Set[T] {
	set := New[T]()
	for e := range seq {
		Add(set, e)
	}
	return set
}

// Clone copies the Set passed to a new Set
func Clone[T comparable](set Set[T]) Set[T] {
	return Collect(maps.Keys(set))
}
