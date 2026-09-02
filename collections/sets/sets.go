package sets

type Set[T comparable] map[T]struct{}

func New[T comparable](elements ...T) Set[T] {
	set := make(Set[T], len(elements))
	for _, e := range elements {
		set[e] = struct{}{}
	}
	return set
}

func Add[T comparable](set Set[T], element T) Set[T] {
	set[element] = struct{}{}
	return set
}

func Has[T comparable](set Set[T], key T) bool {
	_, exists := set[key]
	return exists
}

