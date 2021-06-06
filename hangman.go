package main

import (
	"fmt"
	"sort"
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
		for _, r := range gs.BadGuesses {
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
	if len(matches) == 0 {
		return 0, fmt.Errorf("no guess possible")
	}
	sort.Sort(matches)
	return matches[0].letter, nil
}

func validate(gs GameState) error {
	bad := make(map[rune]bool)
	for _, r := range gs.BadGuesses {
		bad[r] = true
	}
	for _, r := range gs.Current {
		if bad[r] {
			return fmt.Errorf("contradiction on %q", r)
		}
	}
	return nil
}
