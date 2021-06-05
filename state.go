package main

import (
	"bytes"
	"strings"
)

// a hangman-style string with '_' rune ("blank") for "unknown letter"
type state string

func NewState(word string) state {
	b := new(bytes.Buffer)
	for range word {
		b.WriteRune('_')
	}
	return state(b.String())
}

// does state match the given word?
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

// does state still have blanks?
func (s state) unfinished() bool {
	return strings.Contains(string(s), "_")
}

// update the state by filling in the given letter according to target word
func (s state) update(target string, letter rune) state {
	n := len(s)
	if n != len(target) {
		panic("length mismatch")
	}
	stateRunes := []rune(s)
	wordRunes := []rune(target)
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
