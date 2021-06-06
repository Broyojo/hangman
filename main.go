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
	h, err := NewHangman()
	check(err)
	r, err := h.Guess(GameState{
		BadGuesses: nil,
		Current:    "m___h",
	})
	check(err)
	fmt.Printf("guess: %q\n", r)

	return

	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type GameState struct {
	BadGuesses []rune // guesses thus far in the game
	Current    string // a string with correct guesses and underscores
}

func (gs GameState) Guessed() map[rune]bool {
	guessed := make(map[rune]bool)
	for _, r := range gs.BadGuesses {
		guessed[r] = true
	}
	for _, r := range gs.Current {
		if r != '_' {
			guessed[r] = true
		}
	}
	return guessed
}

type Hangman interface {
	Guess(GameState) (rune, error)
}

func Run() error {
	all, err := load()
	if err != nil {
		return err
	}
	var target string
	if len(os.Args) > 1 {
		target := os.Args[1]
		r := Match(target, all)
		dump(r.steps, r.words, target)
		return nil
	}
	const n = 1000
	var ok int
	last := time.Now()
	var i int
	for {
		i++
		target = all[rand.Intn(len(all))]
		r := Match(target, all)
		if r.wrong < 6 {
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
	words []string
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
				addStep(match.letter, true)
				words = filter(words, func(word string) bool {
					return state.matches(word)
				})
				break
			} else {
				addStep(match.letter, false)
				words = filter(words, func(word string) bool {
					return !strings.Contains(word, string(match.letter))
				})
			}
		}
	}
	return result{
		wrong: wrong,
		steps: steps,
		words: words,
	}
}

func filter(words []string, f func(string) bool) (out []string) {
	for _, w := range words {
		if f(w) {
			out = append(out, w)
		}
	}
	return
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

func load() ([]string, error) {
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
