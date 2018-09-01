package rank

import (
	"reflect"
	"testing"
)

func TestInitMatrix(t *testing.T) {
	keys := []string{"a", "b", "c"}
	r := InitMatrix(keys)

	if !reflect.DeepEqual(r.Keys, keys) {
		t.Errorf("%v != %v", r.Keys, keys)
	}

	if len(r.Ranks) != 3*3 {
		t.Errorf("size of Ranks must be n^2 = 9, was %v", len(r.Ranks))
	}

	for i, x := range r.Ranks {
		if x != 0 {
			t.Fatalf("at position %v, Ranks was not 0 but %v", i, x)
		}
	}
}

func TestMatrixFindFree(t *testing.T) {
	cases := []struct {
		mtx     Matrix
		free    Position
		hasFree bool
	}{
		{
			Matrix{
				[]string{"a", "b"},
				[]int{X, A, B, X}},
			Position{}, false,
		},
		{
			Matrix{
				[]string{"a", "b", "c"},
				[]int{X, A, A, B, X, X, B, X, X}},
			Position{"b", "c"}, true,
		},
		{
			Matrix{
				[]string{"a", "b", "c"},
				[]int{X, A, X, B, X, X, X, X, X}},
			Position{"b", "c"}, true,
		},
		{
			Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{X, A, A, A, B, X, A, A, B, B, X, A, B, B, B, X}},
			Position{}, false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			free, hasFree := (&c.mtx).FindFree()

			if c.hasFree != hasFree {
				t.Fatalf("hasFree: %v != %v", c.hasFree, hasFree)
			}

			if hasFree {
				if c.free[0] != free[0] {
					t.Errorf("first key: %v != %v", c.free[0], free[0])
				}
			}
		})
	}
}

func TestMatrixSet(t *testing.T) {
	cases := []struct {
		mtx        Matrix
		pos        Position
		val1, val2 int
		idx1, idx2 int
	}{
		{
			Matrix{
				[]string{"a", "b", "c"},
				[]int{X, A, X, B, X, X, X, X, X}},
			Position{"b", "c"}, A, B, 5, 7,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			c.mtx.Set(c.pos, c.val1)

			if (&c.mtx).Ranks[c.idx1] != c.val1 {
				t.Fatalf("%v != %v", c.mtx.Ranks[c.idx1], c.val1)
			}
			if (&c.mtx).Ranks[c.idx2] != c.val2 {
				t.Fatalf("%v != %v", c.mtx.Ranks[c.idx2], c.val2)
			}
		})
	}
}

func TestMatrixFindCycle(t *testing.T) {
	cases := []struct {
		mtx   *Matrix
		cycle []string
		ok    bool
	}{
		{
			&Matrix{
				[]string{"a", "b", "c"},
				[]int{X, B, A, A, X, B, B, A, X}},
			[]string{"c", "b", "a"}, true,
		},
		{
			&Matrix{
				[]string{"a", "b", "c"},
				[]int{X, A, A, B, X, A, B, B, X}},
			[]string{}, false,
		},
		{
			&Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{
					X, A, A, A,
					B, X, A, B,
					B, B, X, A,
					B, A, B, X,
				}},
			[]string{"c", "d", "b"}, true,
		},
		{
			&Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{
					X, B, B, B,
					A, X, A, B,
					A, B, X, A,
					A, A, B, X,
				}},
			[]string{"c", "d", "b"}, true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			cycle, ok := c.mtx.FindCycle()

			if c.ok != ok {
				t.Fatalf("%v != %v", c.ok, ok)
			}

			if ok {
				if !reflect.DeepEqual(c.cycle, cycle) {
					t.Fatalf("%v != %v", c.cycle, cycle)
				}
			}
		})
	}
}

func TestMatrixRank(t *testing.T) {
	cases := []struct {
		mtx  *Matrix
		rank []string
	}{
		{
			&Matrix{
				[]string{"a", "b", "c"},
				[]int{X, A, A, B, X, B, B, A, X}},
			[]string{"a", "c", "b"},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			rank := c.mtx.Rank()

			if !reflect.DeepEqual(c.rank, rank) {
				t.Fatalf("%v != %v", c.rank, rank)
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	cases := []struct {
		mtx Matrix
	}{
		{
			Matrix{
				[]string{"a", "b"},
				[]int{X, A, B, X}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			s := c.mtx.Serialize()
			d := Deserialize(s)

			if !reflect.DeepEqual(c.mtx.Keys, d.Keys) {
				t.Errorf("%v != %v", d.Keys, c.mtx.Keys)
			}

			if !reflect.DeepEqual(c.mtx.Ranks, d.Ranks) {
				t.Errorf("%v != %v", d.Ranks, c.mtx.Ranks)
			}

		})
	}
}
