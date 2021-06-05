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
	right := "gaze__"
	wrong := "sytlrw"
	dict := LoadDict("words.txt")
	fmt.Println(len(dict))
	matches := FindMatches(dict, right, wrong)
	fmt.Println(matches)
	guess := MakeGuess(matches, right)
	fmt.Println(guess)
}

func LoadDict(path string) []string {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal("failed to open dictionary file")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var dict []string

	for scanner.Scan() {
		text := scanner.Text()
		if len(text) != 1 {
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
