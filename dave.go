package main

import (
	"bytes"
	"fmt"
)

type dave struct {
	dict []string
}

func (d dave) Guess(gs GameState) (rune, error) {
	wrong := new(bytes.Buffer)
	for _, r := range gs.Incorrect {
		wrong.WriteRune(r)
	}
	matches := FindMatches(d.dict, gs.Current, wrong.String())
	guess := MakeGuess(matches, gs.Current)
	runes := []rune(guess)
	if n := len(runes); n != 1 {
		return 0, fmt.Errorf("bad response from code: %d runes (%q)", n, guess)
	}
	return runes[0], nil
}

func NewDave() (Hangman, error) {
	dict := LoadDict("words.txt", 2)
	return &dave{
		dict: dict,
	}, nil
}
