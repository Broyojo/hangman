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
	if len(os.Args) > 1 {
		target = os.Args[1]
	} else {
		target = "document"
	}

	words := load()
	wordsByCount := make(map[int][]string)
	letters := make(map[rune]int)
	for _, w := range words {
		wordsByCount[len(w)] = append(wordsByCount[len(w)], w)
		for _, r := range w {
			letters[r]++
		}
	}
	state := NewState(target)
	fmt.Printf("target = %q\n", target)
	uniques := make(map[rune]bool)
	for _, r := range target {
		uniques[r] = true
	}
	words = wordsByCount[len(target)]
	correctGuesses := make(map[rune]bool)
	guessed := make(map[rune]bool)
	var steps []step
	for len(correctGuesses) < len(uniques) {
		if len(words) == 1 {
			break
		}
		var matches []match
		for letter := range letters {
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
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].entropy > matches[j].entropy
		})
		for _, best := range matches {
			if guessed[best.letter] {
				continue
			}
			guessed[best.letter] = true
			if strings.ContainsRune(target, best.letter) {
				correctGuesses[best.letter] = true
				words = best.dict
				state = state.update(target, best.letter)
				var list []string
				for _, word := range words {
					if state.matches(word) {
						list = append(list, word)
					}
				}
				words = list
				steps = append(steps, step{
					letter:  best.letter,
					state:   state,
					correct: true,
				})
				break
			} else {
				steps = append(steps, step{
					letter:  best.letter,
					state:   state,
					correct: false,
				})
			}
		}
	}
	var wrong int
	for _, s := range steps {
		fmt.Printf("%s, guess %c", s.state, s.letter)
		if !s.correct {
			wrong++
			fmt.Printf(" (%d wrong)", wrong)
		}
		fmt.Println()
	}
	if len(words) == 1 {
		fmt.Printf("the word must be %q\n", words[0])
	}
}

type step struct {
	letter  rune
	state   state
	correct bool
}

type match struct {
	letter  rune
	entropy float64
	dict    []string
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

func load() (out []string) {
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
