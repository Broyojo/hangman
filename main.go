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
	dict := LoadDict("words.txt")
	fmt.Println(len(dict))
	matches := FindMatches(dict, "______", "")
	fmt.Println(matches)
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
