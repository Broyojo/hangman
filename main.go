package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	/*
		dict := LoadDict("words.txt", 2)
		size := 10000
		var failed int
		for i, word := range dict[:size] {
			right := MakeRightWord(len(word))
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
			if i%10 == 0 {
				fmt.Println(float64(failed) / float64(i) * 100)
			}
		}
	*/

	dict := LoadDict("words.txt", 2)
	right := "_oom"
	wrong := "ankl"
	matches := FindMatches(dict, right, wrong)
	fmt.Println(matches)
	guess := MakeGuess(matches, right)
	fmt.Println(guess)
}

func LoadDict(path string, minWordLength int) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("failed to open dictionary file")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var dict []string

	for scanner.Scan() {
		text := scanner.Text()
		if len(text) >= minWordLength {
			dict = append(dict, scanner.Text())
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(dict), func(i, j int) {
		dict[i], dict[j] = dict[j], dict[i]
	})

	return dict
}

func FindMatches(dict []string, right, wrong string) []string {
	var matches []string

	for _, word := range dict {
		// check if length is the same
		if len(word) == len(right) {
			// check if any wrong letters are not in the word
			if !strings.ContainsAny(word, wrong) {
				if func() bool {
					for i := range right {
						if right[i] == '_' {
							// check that word doesn't have more
							// correct letters in missing spots
							for _, char := range right {
								if rune(word[i]) == char {
									return false
								}
							}
							// check if non-underscore characters match
						} else if right[i] != word[i] {
							return false
						}
					}
					return true
				}() {
					matches = append(matches, word)
				}
			}
		}
	}

	return matches
}

func MakeGuess(matches []string, right string) string {
	var maxCount int
	var maxLetter rune

	if len(matches) == 1 {
		return matches[0]
	}

	for i := range right {
		if right[i] == '_' {
			counts := make(map[rune]int)
			for _, match := range matches {
				char := rune(match[i])
				if _, ok := counts[char]; !ok {
					counts[char] = 1
				} else {
					counts[char]++
				}
			}
			fmt.Println(counts)
			for char, count := range counts {
				if count > maxCount {
					maxCount = count
					maxLetter = char
				}
			}
		}
	}

	return string(maxLetter)
}

func MakeRightWord(length int) string {
	var s string
	for i := 0; i < length; i++ {
		s += "_"
	}
	return s
}

func FillInWord(word, right string, letter rune) string {
	s := []byte(right)
	for i := range right {
		if rune(word[i]) == letter {
			s[i] = word[i]
		}
	}
	return string(s)
}
