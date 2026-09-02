package sets

import "testing"

type testCase[T comparable] struct {
	want T
	got  T
}

func TestSets(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		set := New(1, 2, 3)

		if len(set) != 3 {
			t.Fatalf("want len 3, got %d", len(set))
		}

		for _, e := range []int{1, 2, 3} {
			if !Has(set, e) {
				t.Errorf("want set to contain %d, it did not", e)
			}
		}
	})

	t.Run("New/empty", func(t *testing.T) {
		set := New[int]()

		if len(set) != 0 {
			t.Fatalf("want len 0, got %d", len(set))
		}
	})

	t.Run("New/duplicates", func(t *testing.T) {
		set := New(1, 1, 2)

		if len(set) != 2 {
			t.Fatalf("want len 2, got %d", len(set))
		}
	})

	t.Run("Add", func(t *testing.T) {
		set := New[string]()

		Add(set, "a")

		if !Has(set, "a") {
			t.Fatalf("want set to contain %q after Add, it did not", "a")
		}

		if len(set) != 1 {
			t.Fatalf("want len 1, got %d", len(set))
		}
	})

	t.Run("Add/existing", func(t *testing.T) {
		set := New("a")

		Add(set, "a")

		if len(set) != 1 {
			t.Fatalf("want len 1 after adding existing element, got %d", len(set))
		}
	})

	t.Run("Add/returns same set", func(t *testing.T) {
		set := New[int]()

		got := Add(set, 1)

		if len(got) != len(set) || !Has(got, 1) {
			t.Fatalf("want Add to return the mutated set")
		}
	})

	t.Run("Has", func(t *testing.T) {
		set := New("x", "y")

		cases := []testCase[bool]{
			{want: true, got: Has(set, "x")},
			{want: true, got: Has(set, "y")},
			{want: false, got: Has(set, "z")},
		}

		for _, c := range cases {
			if c.want != c.got {
				t.Errorf("want %v, got %v", c.want, c.got)
			}
		}
	})

	t.Run("Has/empty set", func(t *testing.T) {
		set := New[int]()

		if Has(set, 0) {
			t.Fatalf("want Has to be false on empty set, got true")
		}
	})

	t.Run("Iter", func(t *testing.T) {
		slice := []int{1, 2, 3, 4}
		set := New(slice...)
		i := 0
		for e := range set {
			if e != slice[i] {
				t.Fatalf("got %d expected %d\n", e, slice[i])
			}
			i++
		}
	})
}
