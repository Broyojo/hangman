package main

type match struct {
	letter  rune     // the letter matched
	entropy float64  // entropy of words that match vs not
	words   []string // matching words
}

type matches []match

// ordering by high entropy and alphabetically by rune
func (m matches) Less(i, j int) bool {
	entropy := func(i int) float64 {
		return m[i].entropy
	}
	letter := func(i int) rune {
		return m[i].letter
	}
	switch {
	case entropy(i) > entropy(j):
		return true
	case entropy(i) < entropy(j):
		return false
	default:
		return letter(i) < letter(j)
	}
}

func (m matches) Len() int {
	return len(m)
}

func (m matches) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
