package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"time"
)

type i64 int64

func main() {
	list := randomList(26)
	// list := []i64{3, -2, 4, 7, 0, 3, 1, 9, 66, 2, 7, 5, 13, 11, 10, 9}
	// list := []i64{19, 17, 18, 15, 13, 0, 2, 4, 6, 8}
	// list := []i64{1, 1, 2, 1, 1, 3, 1}
	// list := []i64{1, 3, 2}

	m := NewNonStrictlyMonotonicByList(list)
	m.ShowList()
	m.ShowMaxSubList()
}

func randomList(n int) []i64 {
	list := make([]i64, 0, n)
	for i := 0; i < n; i++ {
		v, _ := rand.Int(rand.Reader, big.NewInt(100))
		list = append(list, i64(v.Int64()))
	}
	return list
}

func NewNonStrictlyMonotonic(value i64) *NonStrictlyMonotonic {
	return &NonStrictlyMonotonic{
		list:      []i64{value},
		max:       value,
		min:       value,
		increase:  []*record{{length: 1, idxes: []int{0}, lastValue: value}},
		decrease:  []*record{{length: 1, idxes: []int{0}, lastValue: value}},
		maxLength: 1,
	}
}

func NewNonStrictlyMonotonicByList(list []i64) *NonStrictlyMonotonic {
	defer TimeConsume(time.Now())
	m := NewNonStrictlyMonotonic(list[0])
	for i := 1; i < len(list); i++ {
		m.append(i, list[i])
	}
	return m
}

type NonStrictlyMonotonic struct {
	list []i64

	max, min  i64
	increase  []*record
	decrease  []*record
	maxLength int
}

func (m *NonStrictlyMonotonic) ShowList() {
	fmt.Printf("input list(length: %d):\n", len(m.list))
	list := make([]string, 0, len(m.list))
	for _, v := range m.list {
		list = append(list, fmt.Sprintf("%2d", v))
	}
	fmt.Println(strings.Join(list, " "))
}

func (m *NonStrictlyMonotonic) ShowMaxSubList() {
	fmt.Printf("max sub monotonic length: %d\n", m.maxLength)
	fmt.Println("increase:")
	for _, item := range m.increase {
		if item.length == m.maxLength {
			m.ShowSubList(item.idxes)
		}
	}
	fmt.Println("decrease:")
	for _, item := range m.decrease {
		if item.length == m.maxLength {
			m.ShowSubList(item.idxes)
		}
	}
}

func (m *NonStrictlyMonotonic) ShowSubList(idxes []int) {
	list := make([]string, 0, len(m.list))
	j := 0
	length := len(idxes)
	for i, v := range m.list {
		if j >= length {
			list = append(list, "  ")
			continue
		}
		idx := idxes[j]
		if i == idx {
			list = append(list, fmt.Sprintf("%2d", v))
			j++
			continue
		}
		list = append(list, "  ")
	}
	fmt.Println(strings.Join(list, " "))
}

func (m *NonStrictlyMonotonic) append(i int, value i64) {
	m.list = append(m.list, value)
	m.appendIncrease(i, value)
	m.appendDecrease(i, value)
}

func (m *NonStrictlyMonotonic) appendIncrease(i int, value i64) {
	if value < m.min {
		m.increase = append(m.increase, &record{length: 1, idxes: []int{i}, lastValue: value})
		m.min = value
		return
	}

	increaseMaxLength := 0
	recheckIdxes := make([]int, 0, len(m.increase))
	for idx, item := range m.increase {
		if item.length < increaseMaxLength-1 {
			continue
		}
		if value >= item.lastValue {
			m.increase[idx].idxes = append(m.increase[idx].idxes, i)
			m.increase[idx].lastValue = value
			m.increase[idx].length++
			if m.increase[idx].length > m.maxLength {
				m.maxLength = m.increase[idx].length
			}
			// increaseMaxLength <= a.increase[idx].length
			increaseMaxLength = m.increase[idx].length
		} else if item.length >= increaseMaxLength {
			recheckIdxes = append(recheckIdxes, idx)
		}
	}

	tempMap := make(map[int][][]int)
	for _, idx := range recheckIdxes {
		for j := m.increase[idx].length - 2; j >= 0; j-- {
			if j+2 < increaseMaxLength {
				// 长度不够 剪枝
				break
			}
			checkIdx := m.increase[idx].idxes[j]
			if m.list[checkIdx] <= value {
				newIdxes := make([]int, j+1)
				copy(newIdxes, m.increase[idx].idxes)

				_, ok := tempMap[checkIdx]
				if !ok {
					tempMap[checkIdx] = [][]int{}
				}
				repeat := false
				for _, item := range tempMap[checkIdx] {
					if sliceEq(item, newIdxes) {
						repeat = true
						break
					}
				}
				if repeat {
					break
				}
				tempList := make([]int, j+1)
				copy(tempList, newIdxes)
				tempMap[checkIdx] = append(tempMap[checkIdx], tempList)

				newIdxes = append(newIdxes, i)
				m.increase = append(m.increase, &record{length: j + 2, idxes: newIdxes, lastValue: value})
				increaseMaxLength = j + 2
				break
			}
			// a.list[checkIdx] > value
		}
	}
}

func (m *NonStrictlyMonotonic) appendDecrease(i int, value i64) {
	if value > m.max {
		m.decrease = append(m.decrease, &record{length: 1, idxes: []int{i}, lastValue: value})
		m.max = value
		return
	}

	decreaseMaxLength := 0
	recheckIdxes := make([]int, 0, len(m.decrease))
	for idx, item := range m.decrease {
		if value <= item.lastValue {
			m.decrease[idx].idxes = append(m.decrease[idx].idxes, i)
			m.decrease[idx].lastValue = value
			m.decrease[idx].length++
			if m.decrease[idx].length > m.maxLength {
				m.maxLength = m.decrease[idx].length
			}
			// decreaseMaxLength <= a.decrease[idx].length
			decreaseMaxLength = m.decrease[idx].length
		} else if item.length >= decreaseMaxLength {
			recheckIdxes = append(recheckIdxes, idx)
		}
	}

	tempMap := make(map[int][][]int)
	for _, idx := range recheckIdxes {
		for j := m.decrease[idx].length - 2; j >= 0; j-- {
			if j+2 < decreaseMaxLength {
				// 长度不够 剪枝
				break
			}
			checkIdx := m.decrease[idx].idxes[j]
			if m.list[checkIdx] >= value {
				newIdxes := make([]int, j+1)
				copy(newIdxes, m.decrease[idx].idxes)

				_, ok := tempMap[checkIdx]
				if !ok {
					tempMap[checkIdx] = [][]int{}
				}
				repeat := false
				for _, item := range tempMap[checkIdx] {
					if sliceEq(item, newIdxes) {
						repeat = true
						break
					}
				}
				if repeat {
					break
				}
				tempList := make([]int, j+1)
				copy(tempList, newIdxes)
				tempMap[checkIdx] = append(tempMap[checkIdx], tempList)

				newIdxes = append(newIdxes, i)
				m.decrease = append(m.decrease, &record{length: j + 2, idxes: newIdxes, lastValue: value})
				break
			}
			//  a.list[checkIdx] < value
		}
	}
}

type record struct {
	length    int
	idxes     []int
	lastValue i64
}

func sliceEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TimeConsume(start time.Time) {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return
	}

	// get Fun object from pc
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Printf("%s: %s\n", funcName, time.Since(start).String())
}
