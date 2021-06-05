package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Run() error {
	var target string
	if len(os.Args) != 2 {
		return fmt.Errorf("needs argument")
	}
	fmt.Println()
	target = os.Args[1]
	words, err := load(len(target))
	if err != nil {
		return err
	}
	state := NewState(target)
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
			for _, word := range words {
				if strings.ContainsRune(word, letter) {
					list = append(list, word)
				}
			}
			matches = append(matches, match{
				letter:  letter,
				entropy: entropy2(len(words), len(list)),
				words:   list,
			})
		}
		sort.Sort(matches) // sort by best entropy and alphabetically
		for _, best := range matches {
			if guessed[best.letter] {
				continue
			}
			guessed[best.letter] = true
			if strings.ContainsRune(target, best.letter) {
				correctGuesses[best.letter] = true
				state = state.update(target, best.letter)
				var newWords []string
				for _, word := range best.words {
					if state.matches(word) {
						newWords = append(newWords, word)
					}
				}
				addStep(best.letter, true)
				words = newWords
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
		var msg string
		if words[0] == target {
			msg = "correct"
		} else {
			msg = "wrong"
		}
		fmt.Printf("\nguess %q (%s)\n\n", words[0], msg)
	}
	return nil
}

type step struct {
	letter  rune
	state   state
	correct bool
	words   int
}

// entropy of a set of total size "n" with subdivision of size "x"
func entropy2(n, x int) (out float64) {
	return entropy(n-x, x)
}

// entropy of a set divided into two parts of size "a" and "b"
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

func load(n int) ([]string, error) {
	f, err := os.Open("words.txt.gz")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(gz)
	m := make(map[string]bool)
	var out []string
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
	if err := s.Err(); err != nil {
		return nil, err
	}
	for k := range m {
		out = append(out, k)
	}
	return out, nil
}
