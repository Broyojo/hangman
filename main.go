package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Run() error {
	var mode string
	flag.StringVar(&mode, "m", "dave", "who's algo to run")
	flag.Parse()
	var f func() (Hangman, error)
	switch mode {
	case "mike":
		f = NewHangman
	case "dave":
		f = NewDave
	default:
		return fmt.Errorf("unknown mode %q", mode)
	}
	h, err := f()
	if err != nil {
		return err
	}
	return RunTests(h)
}

func Dave() error {
	dict := LoadDict("words.txt", 2)
	right := "woe"
	wrong := "antsdbmlfpyrhgvc"
	matches := FindMatches(dict, right, wrong)
	fmt.Println(matches)
	guess := MakeGuess(matches, right)
	fmt.Println(guess)
	return nil
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

func MakeEmptyWord(length int) string {
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
