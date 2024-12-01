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

var stdout io.Writer = os.Stdout
var inputFile string = os.Args[1]
var stopwordsFile string = "../../stop_words.txt"

func main() {
	run()
}

func run() {
  printWordFreqs(
    sortWordFreqs(
      countWords(
        removeStopwords(
          scanWords(
            filterCharAndNormalize(
              readFile(
                inputFile,
              ),
            ),
          ),
        ),
      ),
    ),
  )
}

func readFile(inputFile string) string {
	b, err := os.ReadFile(inputFile)
	if err != nil {
		os.Exit(1)
	}
	return string(b)
}

func filterCharAndNormalize(wordStr string) string {
	s := make([]byte, len(wordStr))
	for i := 0; i < len(wordStr); i++ {
		c := wordStr[i]
		if !isAlnum(c) {
			s[i] = ' '
		} else {
			s[i] = lower(c)
		}
	}
	return string(s)
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

func scanWords(wordStr string) []string {
	return strings.Fields(wordStr)
}

func removeStopwords(words []string) []string {
	b, err := os.ReadFile(stopwordsFile)
	if err != nil {
		os.Exit(1)
	}
	stopwords := strings.Split(string(b), ",")

	validWords := []string{}
	for _, word := range words {
		if len(word) < 2 {
			continue
		}
		if inStopwords(word, stopwords) {
			continue
		}

		validWords = append(validWords, word)
	}

	return validWords
}

func inStopwords(word string, stopwords []string) bool {
	for _, stopword := range stopwords {
		if word == stopword {
			return true
		}
	}
	return false
}

func countWords(words []string) []Freq {
	var wordFreqs []Freq
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

	return wordFreqs
}

func sortWordFreqs(wordFreqs []Freq) []Freq {
	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count >= wordFreqs[j].count
	})
	return wordFreqs
}

func printWordFreqs(wordFreqs []Freq) {
	for i := 0; i < 25; i++ {
		freq := wordFreqs[i]
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}
