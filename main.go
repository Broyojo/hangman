package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type state string

func NewState(word string) state {
	b := new(bytes.Buffer)
	for range word {
		b.WriteRune('_')
	}
	return state(b.String())
}

func (s state) matches(w string) bool {
	n := len(s)
	if n != len(w) {
		return false
	}
	stateRunes := []rune(s)
	wordRunes := []rune(w)
	for i := 0; i < n; i++ {
		if r := stateRunes[i]; r == '_' {
			continue
		} else if r != wordRunes[i] {
			return false
		}
	}
	return true
}

func (s state) unfinished() bool {
	return strings.Contains(string(s), "_")
}

func (s state) update(w string, letter rune) state {
	n := len(s)
	if n != len(w) {
		panic("length mismatch")
	}
	stateRunes := []rune(s)
	wordRunes := []rune(w)
	b := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		if letter == wordRunes[i] {
			b.WriteRune(letter)
		} else {
			b.WriteRune(stateRunes[i])
		}
	}
	return state(b.String())
}

func main() {
	var target string
	if len(os.Args) != 2 {
		fmt.Println("needs argument")
		os.Exit(1)
	}
	target = os.Args[1]
	words := load(len(target))
	state := NewState(target)
	uniques := make(map[rune]bool)
	for _, r := range target {
		uniques[r] = true
	}
	correctGuesses := make(map[rune]bool)
	guessed := make(map[rune]bool)
	var steps []step
	addStep := func(r rune, c bool) {
		steps = append(steps, step{
			state:   state,
			letter:  r,
			correct: c,
			words:   len(words),
		})
	}
	for state.unfinished() && len(words) != 1 {
		var matches matches
		for _, letter := range []rune("abcdefghijklmnopqrstuvwxyz") {
			var list []string
			for _, w := range words {
				if strings.ContainsRune(w, letter) {
					list = append(list, w)
				}
			}
			has := len(list)
			matches = append(matches, match{
				letter:  letter,
				entropy: entropy(has, len(words)-has),
				dict:    list,
			})
		}
		sort.Sort(matches)
		for _, best := range matches {
			if guessed[best.letter] {
				continue
			}
			guessed[best.letter] = true
			if strings.ContainsRune(target, best.letter) {
				correctGuesses[best.letter] = true
				state = state.update(target, best.letter)
				var list []string
				for _, word := range best.dict {
					if state.matches(word) {
						list = append(list, word)
					}
				}
				addStep(best.letter, true)
				words = list
				break
			} else {
				addStep(best.letter, false)
			}
		}
	}
	var wrong int
	for _, s := range steps {
		if !s.correct {
			wrong++
		}
		fmt.Printf("%6d words, %2d wrong; guess %q", s.words, wrong, s.letter)
		if s.correct {
			fmt.Printf(" %s", s.state)
		}
		fmt.Println()
	}
	if len(words) == 1 {
		fmt.Printf("\nguess %q\n", words[0])
	}
}

type step struct {
	letter  rune
	state   state
	correct bool
	words   int
}

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

func entropy(a, b int) (out float64) {
	if a == 0 || b == 0 {
		return 0
	}
	n := float64(a + b + 2)
	dist := []float64{float64(a+1) / n, float64(b+1) / n}
	for _, p := range dist {
		out -= p * math.Log2(p)
	}
	return
}

func load(n int) (out []string) {
	f, err := os.Open("words.txt.gz")
	check(err)
	defer f.Close()
	gz, err := gzip.NewReader(f)
	check(err)
	s := bufio.NewScanner(gz)
	m := make(map[string]bool)
	for s.Scan() {
		line := strings.ToLower(strings.TrimSpace(s.Text()))
		if len(line) == 0 {
			continue
		}
		if len(line) != n {
			continue
		}
		m[line] = true
	}
	check(s.Err())
	for k := range m {
		out = append(out, k)
	}
	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
