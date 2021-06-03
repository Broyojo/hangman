package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	pb "github.com/cheggaaa/pb"
)

func main() {
	solver := HangManSolver{}
	solver.LoadDict("words.txt")
	size := 5000
	bar := pb.StartNew(size)
	var failed int
	for _, word := range solver.Dict[:size] {
		solver.Word = word
		solver.CurrentWord = ""
		for i := 0; i < len(solver.Word); i++ {
			solver.CurrentWord += "_"
		}
		solver.WrongLetters = ""
		for !solver.Finished {
			solver.NextMove()
		}
		solver.Finished = false
		solver.Matches = make([]string, 0)
		if len(solver.WrongLetters) >= 6 {
			failed++
		}
		bar.Increment()
	}
	bar.Finish()
	fmt.Println(float64(failed) / float64(size) * 100)
}

type HangManSolver struct {
	Dict         []string
	Word         string
	CurrentWord  string
	WrongLetters string
	Matches      []string
	Finished     bool
}

func (h *HangManSolver) LoadDict(filepath string) {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		log.Fatal("failed to open dictionary file")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	isLetter := func(s string) bool {
		for _, v := range s {
			if v < 'a' || v > 'z' {
				return false
			}
		}
		return true
	}
	for scanner.Scan() {
		text := strings.ToLower(scanner.Text())
		if isLetter(text) {
			h.Dict = append(h.Dict, text)
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(h.Dict), func(i, j int) {
		h.Dict[i], h.Dict[j] = h.Dict[j], h.Dict[i]
	})
}

func (h *HangManSolver) FindMatches() {
	matchesNonUnderscore := func(s string) bool {
		for i := range h.CurrentWord {
			if h.CurrentWord[i] != '_' && h.CurrentWord[i] != s[i] {
				return false
			}
		}
		return true
	}

	if len(h.Matches) == 0 {
		for _, word := range h.Dict {
			if len(word) == len(h.CurrentWord) {
				if !strings.ContainsAny(word, h.WrongLetters) {
					if matchesNonUnderscore(word) {
						h.Matches = append(h.Matches, word)
					}
				}
			}
		}
		if len(h.Matches) == 0 {
			log.Fatal("word not in dictionary")
		}
	} else {
		var list []string
		for _, word := range h.Matches {
			if !strings.ContainsAny(word, h.WrongLetters) {
				if matchesNonUnderscore(word) {
					list = append(list, word)
				}
			}
		}
		h.Matches = list
	}
}

func (h *HangManSolver) NextMove() {
	h.FindMatches()
	var guess byte
	if h.CurrentWord != h.Word {
		if len(h.Matches) == 1 {
			h.CurrentWord = h.Matches[0]
		} else {
			guess = h.FindNextLetter()
			if strings.Contains(h.Word, string(guess)) {
				h.UpdateCurrentWord(guess)
			} else {
				h.WrongLetters += string(guess)
			}
		}
		//fmt.Println(string(guess), h.CurrentWord, h.WrongLetters)
	} else {
		h.Finished = true
	}
}

func (h *HangManSolver) FindNextLetter() byte {
	var guess byte
	var maxCount int
	for i, c := range h.CurrentWord {
		if c == '_' {
			counts := make(map[byte]int)
			for _, match := range h.Matches {
				char := match[i]
				if _, in := counts[char]; !in {
					counts[char] = 1
				} else {
					counts[char]++
				}
			}
			for char, count := range counts {
				if count > maxCount && !strings.Contains(h.CurrentWord, string(char)) {
					maxCount = count
					guess = char
				}
			}
		}
	}
	return guess
}

func (h *HangManSolver) UpdateCurrentWord(guess byte) {
	s := []byte(h.CurrentWord)
	for i := range h.Word {
		if guess == h.Word[i] {
			s[i] = h.Word[i]
		}
	}
	h.CurrentWord = string(s)
}
