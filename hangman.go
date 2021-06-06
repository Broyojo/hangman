package main

import (
	"fmt"
	"strings"
)

type hangman struct {
	words []string
}

func NewHangman() (Hangman, error) {
	all, err := load()
	if err != nil {
		return nil, err
	}
	return &hangman{words: all}, nil
}

func (h hangman) Guess(gs GameState) (rune, error) {
	if err := validate(gs); err != nil {
		return 0, err
	}
	s := state(gs.Current)
	var words []string
	for _, w := range h.words {
		if !s.matches(w) {
			continue
		}
		var exclude bool
		for _, r := range gs.Incorrect {
			if strings.ContainsRune(w, r) {
				exclude = true
			}
		}
		if !exclude {
			words = append(words, w)
		}
	}
	guessed := gs.Guessed()
	var matches matches
	for _, letter := range []rune(lettersByFreq) {
		if guessed[letter] {
			continue
		}
		var v float64
		n := float64(len(words))
		for _, word := range words {
			v += float64(count(word, letter)) / n
		}
		matches = append(matches, match{
			letter: letter,
			value:  v,
		})
	}
	if len(matches) == 0 {
		return 0, fmt.Errorf("no guess possible")
	}
	return matches.Best().letter, nil
}

func count(word string, letter rune) (out int) {
	for _, r := range word {
		if r == letter {
			out++
		}
	}
	return
}

func validate(gs GameState) error {
	bad := make(map[rune]bool)
	for _, r := range gs.Incorrect {
		bad[r] = true
	}
	for _, r := range gs.Current {
		if bad[r] {
			return fmt.Errorf("contradiction on %q", r)
		}
	}
	return nil
}
