/*
sets is a convenience package for re-type of a unique element list.

Under the hood it is a map, therefore beyond the package level functions, all [maps] and [iter] packages' functions get carried over.

Only basic set functions have typed functions in this package.

# Iterating

Simply use a range loop:

	for element := range set {
		// use element
	}

# Composability with stdlib

Under the hood, Set is a map, so it works naturally with the
standard library's [maps] and [iter] APIs.

For [iter] package, a [sets.Iter] function exists for ergonomics. But it is important to understant that this is essentially just a call to [maps.Keys].

# To convert a set to a slice

Use [slices.Collect] function.

  import "slices"
  // ...
  elements := slices.Collect(sets.Iter(set))

Or to get just one element off of a set (a bit of an overkill over using range, but for example's sake)

  next, _ := iter.Pull(sets.Iter(set))
  e, _ := next()

# Design

Adding more wrapper in the sets package therefore does not make sense.
Convert any map to a set simply through [sets.Collect]

  sets.Collect(maps.Keys(someMap))
  sets.Collect(maps.Values(someMap))

Clone an existing set using [sets.Clone]
*/
package sets
