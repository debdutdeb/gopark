/*
sets is a convinience package for re-type of a unique element list.

Under the hood it is a map, therefore beyond the package level functions, all [maps] and [iter] packages' functions get carried over.

Only basic set functions have typed functions in this package.

# Iterating

Simply use a range loop:

	for element := range set {
		// use element
	}

# Composability with stdlib

All [maps] and [iter] package helpers get carried over. For example

# To convert a set to a slice

Use [slices] package.

  elements := slices.Collect(maps.Keys(set))

Or to get just one element off of a set (a bit of an overkill over using range, but for example's sake)


  next, _ := iter.Pull(maps.Keys(set))
  e, _ := next()


# Extra

Adding more wrapper in sets package therefore does not make sense.
Convert any map to a set simply through [sets.Collect]

  sets.Collect(maps.Keys(someMap))
  sets.Collect(maps.Values(someMap))

Clone an existing set using `sets.Clone`
*/
package sets
