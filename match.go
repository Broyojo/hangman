package main

import (
	"fmt"
)

type match struct {
	letter rune    // the letter matched
	value  float64 // number of matching words
}

func (m match) String() string {
	return fmt.Sprintf("%q:%f", m.letter, m.value)
}

type matches []match

const (
	strategy      = "freq" // "freq" vs "entropy"
	lettersByFreq = "etaoinsrhdlucmfywgpbvkxqjz"
)

var letterIndicies map[rune]int

func init() {
	letterIndicies = make(map[rune]int)
	for i, r := range lettersByFreq {
		letterIndicies[r] = i
	}
}

func Less(a, b rune) bool {
	index := func(r rune) int {
		i, ok := letterIndicies[r]
		if !ok {
			return len(letterIndicies)
		}
		return i
	}
	return index(a) < index(b)
}

func (m matches) Best() match {
	best := -1
	for i := range m {
		if best == -1 {
			best = i
		}
		if m.Less(i, best) {
			best = i
		}
	}
	return m[best]
}

func (m matches) Less(i, j int) bool {
	letter := func(i int) rune {
		return m[i].letter
	}
	value := func(i int) float64 {
		return m[i].value
	}
	switch {
	case value(i) > value(j):
		return true
	case value(i) < value(j):
		return false
	default:
		return Less(letter(i), letter(j))
	}
}

func (m matches) Len() int {
	return len(m)
}

func (m matches) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
