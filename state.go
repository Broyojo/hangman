package main

import (
	"bytes"
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

func (s state) unfinished() bool {
	return strings.Contains(string(s), "_")
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
