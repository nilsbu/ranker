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
				[]int{X, X, X, X}},
			Position{"a", "b"}, true,
		},
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
				[]int{
					X, A, A, A,
					B, X, A, A,
					B, B, X, A,
					B, B, B, X}},
			Position{}, false,
		},
		{
			Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{
					X, X, X, X,
					X, X, BB, IB,
					X, AA, X, BB,
					X, IA, AA, X}},
			Position{"a", "d"}, true,
		},
		{
			Matrix{
				[]string{"a", "b", "c", "d", "e"},
				[]int{
					X, BB, X, X, X,
					AA, X, X, X, X,
					X, X, X, X, X,
					X, X, X, X, X,
					X, X, X, X, X,
				}},
			Position{"c", "a"}, true,
		},
		{
			Matrix{
				[]string{"a", "b", "c", "d", "e"},
				[]int{
					X, AA, A, X, X,
					BB, X, X, X, X,
					B, X, X, X, X,
					X, X, X, X, X,
					X, X, X, X, X,
				}},
			Position{"b", "c"}, true,
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
				if c.free[1] != free[1] {
					t.Errorf("secind key: %v != %v", c.free[1], free[1])
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
		{
			Matrix{
				[]string{"a", "b", "c"},
				[]int{X, A, X, B, X, X, X, X, X}},
			Position{"a", "b"}, X, X, 1, 4,
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

func TestMatrixSetImplied(t *testing.T) {
	cases := []struct {
		raw    *Matrix
		filled *Matrix
		ok     bool
	}{
		{
			&Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{
					X, AA, X, X,
					BB, X, AA, X,
					X, BB, X, AA,
					X, X, BB, X,
				}},
			&Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{
					X, AA, IA, IA,
					BB, X, AA, IA,
					IB, BB, X, AA,
					IB, IB, BB, X,
				}},
			true,
		},
		{
			&Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{
					X, AA, X, X,
					BB, X, AA, X,
					X, BB, X, A,
					X, X, B, X,
				}},
			&Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{
					X, AA, IA, X,
					BB, X, AA, X,
					IB, BB, X, A,
					X, X, B, X,
				}},
			true,
		},
		{
			&Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{
					X, AA, BB, X,
					BB, X, AA, X,
					AA, BB, X, X,
					X, X, X, X,
				}},
			&Matrix{
				[]string{"a", "b", "c", "d"},
				[]int{}},
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			filled, ok := c.raw.SetImplied()

			if ok != c.ok {
				t.Errorf("%v != %v", c.ok, ok)
			}

			if ok {
				if !reflect.DeepEqual(c.filled.Keys, filled.Keys) {
					t.Errorf("%v != %v", c.filled.Keys, filled.Keys)
				}

				if !reflect.DeepEqual(c.filled.Ranks, filled.Ranks) {
					t.Errorf("%v != %v", c.filled.Ranks, filled.Ranks)
				}
			}
		})
	}
}

func TestMatrixCountFree(t *testing.T) {
	cases := []struct {
		mtx  *Matrix
		free int
	}{
		{
			&Matrix{
				[]string{"a", "b", "c"},
				[]int{X, A, A, B, X, B, B, A, X}},
			0,
		},
		{
			&Matrix{
				[]string{"a", "b", "c"},
				[]int{X, A, X, B, X, X, X, X, X}},
			2,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			free := c.mtx.CountFree()

			if c.free != free {
				t.Fatalf("%v != %v", c.free, free)
			}
		})
	}
}

func TestMatrixClearImplied(t *testing.T) {
	cases := []struct {
		pre  *Matrix
		post *Matrix
	}{
		{
			&Matrix{
				[]string{"a", "b", "c"},
				[]int{X, AA, IA, BB, X, AA, IB, BB, X}},
			&Matrix{
				[]string{"a", "b", "c"},
				[]int{X, AA, X, BB, X, AA, X, BB, X}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			post := c.pre.ClearImplied()

			if !reflect.DeepEqual(c.post.Ranks, post.Ranks) {
				t.Fatalf("%v != %v", c.post.Ranks, post.Ranks)
			}
		})
	}
}

func TestMatrixSerialize(t *testing.T) {
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
