package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Freq struct {
	word  string
	count int
}

const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var wordStr string
var words []string
var stopwords []string
var wordFreqs []Freq
var stdout io.Writer = os.Stdout
var inputFile string = os.Args[1]
var stopwordsFile string = "../../stop_words.txt"

func main() {
	run()
}

func run() {
	readFile()
	filterCharAndNormalize()
	scanWords()
	removeStopwords()
	countWords()
	sortWordFreqs()
	printWordFreqs()
}

func readFile() {
	b, err := os.ReadFile(inputFile)
	if err != nil {
		os.Exit(1)
	}
	wordStr = string(b)
}

func filterCharAndNormalize() {
	s := make([]byte, len(wordStr))
	for i := 0; i < len(wordStr); i++ {
		c := wordStr[i]
		if !isAlnum(c) {
			s[i] = ' '
		} else {
			s[i] = lower(c)
		}
	}
	wordStr = string(s)
}

func isAlnum(c byte) bool {
	for i := 0; i < len(alnum); i++ {
		if c == alnum[i] {
			return true
		}
	}
	return false
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		c = c - 'A' + 'a'
	}
	return c
}

func scanWords() {
	words = strings.Fields(wordStr)
}

func removeStopwords() {
	b, err := os.ReadFile(stopwordsFile)
	if err != nil {
		os.Exit(1)
	}
	stopwords = strings.Split(string(b), ",")

	validWords := []string{}
	for _, word := range words {
		if len(word) < 2 {
			continue
		}
		if inStopwords(word) {
			continue
		}

		validWords = append(validWords, word)
	}

	words = validWords
}

func inStopwords(word string) bool {
	for _, stopword := range stopwords {
		if word == stopword {
			return true
		}
	}
	return false
}

func countWords() {
OUTER:
	for _, word := range words {
		for i, freq := range wordFreqs {
			if freq.word == word {
				wordFreqs[i].count += 1
				continue OUTER
			}
		}

		wordFreqs = append(wordFreqs, Freq{
			word:  word,
			count: 1,
		})
	}
}

func sortWordFreqs() {
	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count >= wordFreqs[j].count
	})
}

func printWordFreqs() {
	for i := 0; i < 25; i++ {
		freq := wordFreqs[i]
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}
