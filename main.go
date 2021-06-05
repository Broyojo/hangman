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
	words := load()
	wordsByCount := make(map[int][]string)
	letters := make(map[rune]int)
	for _, w := range words {
		wordsByCount[len(w)] = append(wordsByCount[len(w)], w)
		for _, r := range w {
			letters[r]++
		}
	}
	const word = "matter"
	uniques := make(map[rune]bool)
	for _, r := range word {
		uniques[r] = true
	}
	words = wordsByCount[len(word)]
	correctGuesses := make(map[rune]bool)
	guessed := make(map[rune]bool)

	for len(correctGuesses) < len(uniques) {
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
			correct := strings.ContainsRune(word, best.letter)
			if correct {
				correctGuesses[best.letter] = true
				words = best.dict
			}
			fmt.Printf("%d / %d. guess %c: %v; %d left\n",
				len(guessed)-len(correctGuesses),
				len(guessed),
				best.letter, correct, len(best.dict),
			)
			if correct {
				break
			}
		}
	}
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
	n := float64(a + b)
	dist := []float64{float64(a) / n, float64(b) / n}
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
