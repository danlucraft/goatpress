package goatpress

import (
	"fmt"
	"sort"
	"testing"
)

func TestHashWordSet(t *testing.T) {
	set := newWordSet()
	set.Add("hello")
	set.Add("hi")
	a1 := set.Includes("hello")
	if !a1 {
		t.Errorf("include hello failed", a1, true)
	}
	a2 := set.Includes("hippie")
	if a2 {
		t.Errorf("include hippie failed", a2, false)
	}
	a3 := set.ChooseRandom()
	if a3 != "hi" && a3 != "hello" {
		t.Errorf("ChooseRandom didn't choose one")
	}
}

func TestNewWordSetFromFile(t *testing.T) {
	set := DefaultWordSet
	if !set.Includes("aa") {
		t.Errorf("wordSet doesn't include aa")
	}
	if set.Includes("a") {
		t.Errorf("wordSet includes a")
	}
	if set.Length() != 210661 {
		t.Errorf("wordSet not right length", set.Length(), 210661)
	}
}

func filterInts(ints []int, predicate func(int) bool) []int {
	result := make([]int, 0)
	for _, i := range ints {
		if predicate(i) {
			result = append(result, i)
		}
	}
	return result
}

func TestFuncFilterInts(t *testing.T) {
	ints := []int{1, -3, 0, 4, -7}
	fmt.Println(ints)
	IsPositive := func(i int) bool {
		return i > 0
	}
	result := filterInts(ints, IsPositive)
	fmt.Println(result)
}

func TestFuncFilterInts2(t *testing.T) {
	ints := []int{1, -3, 0, 4, -7}
	fmt.Println(ints)
	result := filterInts(ints, func(i int) bool { return i > 0 })
	fmt.Println(result)
}

func mapInts(ints []int, mapper func(int) int) []int {
	result := make([]int, 0)
	for _, i := range ints {
		result = append(result, mapper(i))
	}
	return result
}

func TestFuncMapInts(t *testing.T) {
	ints := []int{1, -3, 0, 4, -7}
	fmt.Println(ints)
	result := mapInts(ints, func(i int) int { return i * 2 })
	fmt.Println(result)
}

type Thing struct {
	name  string
	score int
}

type IntSorter struct {
	Pairs [][2]int
}

func newSorterByInt() *IntSorter {
	pairs := make([][2]int, 0)
	return &IntSorter{pairs}
}

func (s *IntSorter) Add(index int, value int) {
	s.Pairs = append(s.Pairs, [2]int{index, value})
}

func (s *IntSorter) Sort() {
	sort.Sort(s)
}

func (s *IntSorter) Len() int {
	return len(s.Pairs)
}

func (s *IntSorter) Less(i int, j int) bool {
	return s.Pairs[i][1] < s.Pairs[j][1]
}

func (s *IntSorter) Swap(i int, j int) {
	tmp := s.Pairs[i]
	s.Pairs[i][0] = s.Pairs[j][0]
	s.Pairs[i][1] = s.Pairs[j][1]
	s.Pairs[j][0] = tmp[0]
	s.Pairs[j][1] = tmp[1]
}

// change the type the declaration and the result type
//func sortXByInt(xs []X, mapper func (Thing) int) []X {
//  result := make([]X, len(xs))
//  
//  sorter := newSorterByInt()
//  for i, x := range xs {
//    sorter.Add(i, mapper(x))
//  }
//  sorter.Sort()
//  for i, pair := range sorter.Pairs {
//    result[i] = xs[pair[0]]
//  }
//  return result
//}

func sortThingByInt(xs []Thing, mapper func(Thing) int) []Thing {
	result := make([]Thing, len(xs))

	sorter := newSorterByInt()
	for i, x := range xs {
		sorter.Add(i, mapper(x))
	}
	sorter.Sort()
	for i, pair := range sorter.Pairs {
		result[i] = xs[pair[0]]
	}
	return result
}

func TestFuncSortBy(t *testing.T) {
	things := []Thing{Thing{"dan", 10}, Thing{"alice", 2}}
	fmt.Println(things)
	sortedThings := sortThingByInt(things, func(t Thing) int { return t.score })
	fmt.Println(sortedThings)
}

// This one isn't so bad!!

type SortableThings []Thing

func (a *SortableThings) Len() int {
	return len([]Thing(*a))
}

func (a *SortableThings) Less(i int, j int) bool {
	ts := []Thing(*a)
	return ts[i].score < ts[j].score
}

func (a *SortableThings) Swap(i int, j int) {
	ts := []Thing(*a)
	tmp := ts[i]
	ts[i] = ts[j]
	ts[j] = tmp
}

func TestFuncSortBy2(t *testing.T) {
	things := []Thing{Thing{"dan", 10}, Thing{"alice", 2}}
	fmt.Println(things)

	thingList := SortableThings(things)
	sort.Sort(&thingList)

	fmt.Println(things)
}

func testWordSet() WordSet {
	words := newWordSet()
	words.Add("hello")
	words.Add("state")
	words.Add("jenga")
	words.Add("pages")
	words.Add("valid")
	return words
}
