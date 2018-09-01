package rank

import "sort"

const (
	X int = iota
	A
	B
)

type Position [2]string

type Matrix struct {
	Keys  []string
	Ranks []int
}

func InitMatrix(keys []string) *Matrix {
	return &Matrix{
		keys,
		make([]int, len(keys)*len(keys)),
	}
}

func (mtx *Matrix) FindFree() (Position, bool) {
	for i := 1; i < len(mtx.Keys); i++ {
		for j := 0; i+j < len(mtx.Keys); j++ {
			y, x := j, i+j
			if mtx.Ranks[y*len(mtx.Keys)+x] == X {
				return Position{mtx.Keys[y], mtx.Keys[x]}, true
			}
		}
	}

	return Position{}, false
}

func (mtx *Matrix) Set(pos Position, value int) {
	y, x := 0, 0

	for i, key := range mtx.Keys {
		if key == pos[0] {
			y = i
			break
		}
	}
	for i, key := range mtx.Keys {
		if key == pos[1] {
			x = i
			break
		}
	}

	mtx.Ranks[y*len(mtx.Keys)+x] = value
	mtx.Ranks[x*len(mtx.Keys)+y] = 3 - value
}

func (mtx *Matrix) FindCycle() (cycle []string, ok bool) {
	visited := make([]bool, len(mtx.Keys))

	for j := range visited {
		if indices, ok := findCycle(mtx.Ranks, visited, j); ok {
			return transcribe(indices, mtx.Keys), true
		}
	}

	return []string{}, false
}

func findCycle(ranks []int, visited []bool, i int) (indices []int, ok bool) {
	if visited[i] {
		return []int{i}, true
	}

	visited[i] = true

	for j := range visited {
		if ranks[i*len(visited)+j] == A {
			if indices, ok = findCycle(ranks, visited, j); ok {
				return append(indices, i), true
			}
		}
	}

	visited[i] = false

	return []int{}, false
}

func transcribe(indices []int, keys []string) (cycle []string) {
	enter := false
	for i := len(indices) - 1; i >= 0; i-- {
		if enter {
			cycle = append(cycle, keys[indices[i]])
		}
		if indices[i] == indices[0] {
			enter = true
		}
	}

	return
}

type ks struct {
	k string
	s int
}

func (mtx *Matrix) Rank() (keys []string) {
	scores := make([]ks, len(mtx.Keys))
	for i := range mtx.Keys {
		scores[i].k = mtx.Keys[i]
		for j := range mtx.Keys {
			if mtx.Ranks[i*len(mtx.Keys)+j] == A {
				scores[i].s++
			}
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].s > scores[j].s
	})

	for _, key := range scores {
		keys = append(keys, key.k)
	}

	return keys
}
