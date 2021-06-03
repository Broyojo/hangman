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
	solver := HangManSolver{}
	solver.LoadDict("words.txt")

	max_wrong := 0
	max_wrong_word := ""

	//bar := pb.StartNew(10000)

	for _, word := range []string{"crwth"} {
		correct := ""
		for i := 0; i < len(word); i++ {
			correct += "_"
		}
		incorrect := ""
		num_wrong := 0
		for strings.Contains(correct, "_") {
			matches := solver.FindMatches(correct, incorrect)

			char, _, _, full_guess := solver.FindBestGuess(correct, matches)

			if full_guess != "" {
				correct = full_guess
			} else {
				if strings.Contains(word, string(char)) {
					s := []rune(correct)
					for k, v := range word {
						if v == char {
							s[k] = v
						}
					}
					correct = string(s)
				} else {
					incorrect += string(char)
					num_wrong++
				}
			}

			fmt.Println(correct, string(char), num_wrong)
		}
		//bar.Increment()
		if num_wrong > max_wrong {
			max_wrong = num_wrong
			max_wrong_word = word
		}
	}
	fmt.Println("max mistakes:", max_wrong)
	fmt.Println("hardest word:", max_wrong_word)
	//bar.Finish()
}

type HangManSolver struct {
	Dict []string
}

func (h *HangManSolver) LoadDict(filepath string) {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		log.Fatal("failed to open dictionary file")
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	isValid := func(s string) bool {
		for _, v := range s {
			if v < 'a' || v > 'z' {
				return false
			}
		}
		return true
	}

	for scanner.Scan() {
		text := strings.ToLower(scanner.Text())
		if isValid(text) {
			h.Dict = append(h.Dict, text)
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(h.Dict), func(i, j int) {
		h.Dict[i], h.Dict[j] = h.Dict[j], h.Dict[i]
	})
}

func (h *HangManSolver) FindMatches(correct, incorrect string) []string {
	var matches []string

	// sees if correct letters are in correct spot in word
	valid := func(s string) bool {
		for k := range s {
			if correct[k] != '_' && s[k] != correct[k] {
				return false
			}
		}
		return true
	}

	for _, word := range h.Dict {
		if len(word) == len(correct) {
			if !strings.ContainsAny(word, incorrect) {
				if valid(word) {
					matches = append(matches, word)
				}
			}
		}
	}

	return matches
}

func (h *HangManSolver) FindBestGuess(correct string, matches []string) (rune, int, int, string) {
	/*

		loop through every match and every character (the ones that are missing from the correct string)
		and find the highest probability letters for each slot
		the next turn is picking the highest probability letter out of all of them

	*/

	if len(matches) == 1 {
		return 0, 0, 0, matches[0]
	}

	var indices []int
	for i := range correct {
		if correct[i] == '_' {
			indices = append(indices, i)
		}
	}

	var maxChar uint8
	var maxCount int
	var maxIndex int
	for _, index := range indices {
		counts := make(map[byte]int)
		for _, word := range matches {
			char := word[index]
			if _, in := counts[char]; !in {
				counts[char] = 1
			} else {
				counts[char] += 1
			}
		}
		for k, v := range counts {
			if v > maxCount && !strings.ContainsRune(correct, rune(k)) {
				maxCount = v
				maxChar = k
				maxIndex = index
			}
		}
	}

	// letter, count, index
	return rune(maxChar), maxCount, maxIndex, ""
}
