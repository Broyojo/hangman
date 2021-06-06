package main

import "fmt"

type match struct {
	letter  rune     // the letter matched
	entropy float64  // entropy of words that match vs not
	words   []string // matching words
}

func (m match) String() string {
	return fmt.Sprintf("%q (%.4f) %d words", m.letter, m.entropy, len(m.words))
}

type matches []match

const (
	strategy      = "freq" // vs "entropy"
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

// ordering by high entropy and alphabetically by rune
func (m matches) Less(i, j int) bool {
	letter := func(i int) rune {
		return m[i].letter
	}
	switch strategy {
	case "freq":
		freq := func(i int) int {
			return len(m[i].words)
		}
		switch {
		case freq(i) > freq(j):
			return true
		case freq(i) < freq(j):
			return false
		}
	case "entropy":
		entropy := func(i int) float64 {
			return m[i].entropy
		}
		switch {
		case entropy(i) > entropy(j):
			return true
		case entropy(i) < entropy(j):
			return false
		}
	default:
		panic(strategy)
	}
	return Less(letter(i), letter(j))
}

func (m matches) Len() int {
	return len(m)
}

func (m matches) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
