package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
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
	target = os.Args[1]
	all, err := load(0)
	if err != nil {
		return err
	}
	const n = 1000
	var ok int
	last := time.Now()
	var i int
	for {
		i++
		target = all[rand.Intn(len(all))]
		r := Match(target, all)
		if r.wrong <= 8 {
			ok++
		}
		if time.Since(last) > time.Second {
			fmt.Printf("%.2f%% ok (%d tries)\n", 100*float64(ok)/float64(i), i)
			last = time.Now()
		}
	}
	return nil
}

type result struct {
	wrong int
	steps []step
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Match(target string, all []string) result {
	var words []string
	for _, w := range all {
		if len(w) != len(target) {
			continue
		}
		words = append(words, w)
	}
	state := NewState(target)
	guessed, correctGuesses := make(map[rune]bool), make(map[rune]bool)
	var steps []step
	var wrong int
	addStep := func(r rune, correct bool) {
		if !correct {
			wrong++
		}
		steps = append(steps, step{
			state:   state,
			letter:  r,
			correct: correct,
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
		for _, match := range matches {
			if guessed[match.letter] {
				continue
			}
			guessed[match.letter] = true
			if strings.ContainsRune(target, match.letter) {
				correctGuesses[match.letter] = true
				state = state.update(target, match.letter)
				var newWords []string
				for _, word := range match.words {
					if state.matches(word) {
						newWords = append(newWords, word)
					}
				}
				addStep(match.letter, true)
				words = newWords
				break
			} else {
				addStep(match.letter, false)
			}
		}
	}
	return result{
		wrong: wrong,
		steps: steps,
	}
}

type step struct {
	letter  rune
	state   state
	correct bool
	words   int
}

func dump(steps []step, words []string, target string) {
	fmt.Println()
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

func load(min int) ([]string, error) {
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
		if len(line) <= min {
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
