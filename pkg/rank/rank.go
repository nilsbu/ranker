package rank

import (
	"encoding/json"
	"sort"
)

const (
	X int = iota
	A
	B
	AA
	BB
	IA
	IB
)

type indices [2]int
type Position [2]string

type Matrix struct {
	Keys  []string `json:"keys"`
	Ranks []int    `json:"ranks"`
}

func InitMatrix(keys []string) *Matrix {
	return &Matrix{
		keys,
		make([]int, len(keys)*len(keys)),
	}
}

func (mtx *Matrix) at(x, y int) int {
	return mtx.Ranks[y*len(mtx.Keys)+x]
}

func (mtx *Matrix) FindFree() (Position, bool) {
	type res struct {
		idx indices
		n   int
	}

	if pos, ok := mtx.findFree(); ok {
		impr := []res{{pos, 0}}
		if best, n, ok := mtx.getOptimalFreeSpace(pos, pos, BB); ok {
			impr = append(impr, res{best, n})
		}
		pos = indices{pos[1], pos[0]}
		if best, n, ok := mtx.getOptimalFreeSpace(pos, pos, AA); ok {
			impr = append(impr, res{best, n})
		}

		idx := 0
		n := 0
		for i := range impr {
			if n < impr[i].n {
				idx = i
				n = impr[i].n
			}
		}

		pos = impr[idx].idx
		return Position{mtx.Keys[pos[0]], mtx.Keys[pos[1]]}, true
	} else {
		return Position{}, false
	}
}

func (mtx *Matrix) findFree() (indices, bool) {
	for i := 1; i < len(mtx.Keys); i++ {
		for j := 0; i+j < len(mtx.Keys); j++ {
			y, x := j, i+j
			if mtx.Ranks[y*len(mtx.Keys)+x] == X {
				return indices{y, x}, true
			}
		}
	}

	return indices{}, false
}

func (mtx *Matrix) getOptimalFreeSpace(
	origin, idcs indices,
	dir int,
) (best indices, n int, ok bool) {
	type res struct {
		idx indices
		n   int
	}

	impr := []res{}

	y := idcs[1]
	for x := 0; x < len(mtx.Keys); x++ {
		if mtx.at(x, y) == dir {
			if best, n, ok = mtx.getOptimalFreeSpace(origin, indices{y, x}, dir); ok {
				impr = append(impr, res{best, n + 1})
			}
		}
	}

	if mtx.at(origin[0], idcs[1]) == X {
		impr = append(impr, res{idcs, 0})
	}

	if len(impr) == 0 {
		return indices{}, 0, false
	}

	idx := 0
	cnt := 0
	for i := range impr {
		if cnt < impr[i].n {
			idx = i
			cnt = impr[i].n
		}
	}

	return indices{origin[0], impr[idx].idx[1]}, impr[idx].n, true
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

	mtx.set(x, y, value)
}

func (mtx *Matrix) set(x, y int, value int) {
	var other int
	if value == X {
		other = x
	} else {
		other = (value-1)/2*2 + value%2 + 1
	}

	mtx.Ranks[y*len(mtx.Keys)+x] = value
	mtx.Ranks[x*len(mtx.Keys)+y] = other
}

func (mtx *Matrix) FindCycle() (cycle []string, ok bool) {
	visited := make([]bool, len(mtx.Keys))
	recStack := make([]bool, len(mtx.Keys))

	for j := range visited {
		if indices, ok := findCycle(mtx.Ranks, visited, recStack, j); ok {
			return transcribe(indices, mtx.Keys), true
		}
	}

	return []string{}, false
}

func findCycle(ranks []int, visited, recStack []bool, i int) (indices []int, ok bool) {
	if visited[i] {
		recStack[i] = false
		return []int{}, false
	}

	visited[i] = true
	recStack[i] = true

	for j := range visited {
		r := ranks[i*len(visited)+j]
		if r == A || r == AA || r == IA {
			if !visited[j] {
				if indices, ok = findCycle(ranks, visited, recStack, j); ok {
					return append(indices, i), true
				}
			}

			if recStack[j] {
				return []int{j, i}, true
			}
		}
	}

	recStack[i] = false

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
			r := mtx.Ranks[i*len(mtx.Keys)+j]
			if r == A || r == AA || r == IA {
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

func (mtx *Matrix) SetImplied() (filled *Matrix, ok bool) {
	filled = &Matrix{
		Keys:  mtx.Keys,
		Ranks: make([]int, len(mtx.Ranks)),
	}

	for i, v := range mtx.Ranks {
		filled.Ranks[i] = v
	}

	for j := range mtx.Keys {
		if !filled.fillImplied([]int{}, j) {
			return filled, false
		}
	}

	return filled, true
}

func (mtx *Matrix) fillImplied(visited []int, i int) (ok bool) {
	for _, idx := range visited {
		r := mtx.Ranks[idx*len(mtx.Keys)+i]
		if r == X {
			mtx.set(i, idx, IA)
		} else if r == B || r == BB || r == IB {
			return false
		}
	}

	visited = append(visited, i)

	for j := range mtx.Keys {
		r := mtx.Ranks[i*len(mtx.Keys)+j]
		if r == AA {
			if !mtx.fillImplied(visited, j) {
				return false
			}
		}
	}

	return true
}

func (mtx *Matrix) CountFree() int {
	count := 0
	for i := 0; i < len(mtx.Keys); i++ {
		for j := i + 1; j < len(mtx.Keys); j++ {
			if mtx.Ranks[i*len(mtx.Keys)+j] == X {
				count++
			}
		}
	}

	return count
}

func (mtx *Matrix) ClearImplied() *Matrix {
	post := &Matrix{
		Keys:  mtx.Keys,
		Ranks: make([]int, len(mtx.Ranks)),
	}

	for i, v := range mtx.Ranks {
		if v == IA || v == IB {
			post.Ranks[i] = X
		} else {
			post.Ranks[i] = v
		}
	}

	return post
}

func (mtx *Matrix) Serialize() []byte {
	b, _ := json.Marshal(mtx)
	return b
}

func Deserialize(bytes []byte) *Matrix {
	mtx := &Matrix{}
	json.Unmarshal(bytes, mtx)
	return mtx
}
