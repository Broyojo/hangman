package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Hangman interface {
	// one step in the game:
	Guess(GameState) (rune, error)
}

type GameState struct {
	Current   string // a string with correct guesses and underscores
	Incorrect []rune // incorrect guesses thus far in the game
}

// tests a hangman algo and returns number of wrong guesses
func Test(h Hangman, word string) (int, error) {
	s := NewState(len(word))
	gs := GameState{
		Current: string(s),
	}
	guessed := make(map[rune]bool)
	for s.unfinished() {
		r, err := h.Guess(gs)
		if err != nil {
			return 0, err
		}
		if _, ok := letterIndicies[r]; !ok {
			return 0, fmt.Errorf("illegal guess %q", r)
		}
		if guessed[r] {
			return 0, fmt.Errorf("already guessed %q", r)
		}
		guessed[r] = true
		if strings.ContainsRune(word, r) {
			s = s.update(word, r)
			gs.Current = string(s)
		} else {
			gs.Incorrect = append(gs.Incorrect, r)
		}
	}
	return len(gs.Incorrect), nil
}

func Mike() error {
	h, err := NewHangman()
	if err != nil {
		return err
	}
	return RunTests(h)
}

func RunTests(h Hangman) error {
	const word = "comfortable"
	n, err := Test(h, word)
	if err != nil {
		return err
	}
	fmt.Printf("%d bad guesses for %q\n", n, word)
	if err := Stats(h); err != nil {
		return err
	}
	return nil
}

func (gs GameState) Guessed() map[rune]bool {
	guessed := make(map[rune]bool)
	for _, r := range gs.Incorrect {
		guessed[r] = true
	}
	for _, r := range gs.Current {
		if r != '_' {
			guessed[r] = true
		}
	}
	return guessed
}

func Stats(h Hangman) error {
	all, err := load()
	if err != nil {
		return err
	}

	var ok int
	last := time.Now()
	var i int
	for {
		i++
		target := all[rand.Intn(len(all))]
		n, err := Test(h, target)
		if err != nil {
			return err
		}
		if n < 6 {
			ok++
		}
		if time.Since(last) > time.Second {
			fmt.Printf("%.2f%% ok (%d tries)\n", 100*float64(ok)/float64(i), i)
			last = time.Now()
		}
	}
	return nil
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

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
