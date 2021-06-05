package main

type match struct {
	letter  rune
	entropy float64
	dict    []string
}

type matches []match

func (m matches) Len() int {
	return len(m)
}

func (m matches) Less(i, j int) bool {
	e := func(i int) float64 {
		return m[i].entropy
	}
	l := func(i int) rune {
		return m[i].letter
	}
	switch {
	case e(i) > e(j):
		return true
	case e(i) < e(j):
		return false
	default:
		return l(i) < l(j)
	}
}

func (m matches) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
