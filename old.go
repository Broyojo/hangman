package main

import (
	"fmt"
	"strings"
)

func oldcode() {
	dict := LoadDict("words.txt", 2)
	size := 10000
	var failed int
	for i, word := range dict[:size] {
		right := MakeEmptyWord(len(word))
		var wrong string

		matches := FindMatches(dict, right, wrong)

		for strings.ContainsRune(right, '_') {
			guess := MakeGuess(matches, right)

			if strings.Contains(word, guess) {
				// right guess
				if len(guess) == 1 {
					right = FillInWord(word, right, rune(guess[0]))
				} else {
					right = guess
				}
			} else {
				// wrong guess
				wrong += guess
			}

			matches = FindMatches(matches, right, wrong)
			//fmt.Println(right, len(wrong))
			//fmt.Println(len(matches))
		}
		if len(wrong) >= 6 {
			failed++
		}
		fmt.Println(float64(failed) / float64(i) * 100)
	}
}
