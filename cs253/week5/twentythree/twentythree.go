package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Freq struct {
	word  string
	count int
}

var stdout io.Writer
var stderr io.Writer
var inputFile string
var stopwordsFile string = "../../stop_words.txt"

var (
	ErrNotEnoughArgument = errors.New("not enough argument")
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(stderr, "recover from: %v\n", r)
			os.Exit(1)
		}
	}()

	stdout = os.Stdout
	if len(os.Args) < 2 {
		panic(ErrNotEnoughArgument)
	}

	inputFile = os.Args[1]
	run()
}

func run() {
	words := extractWords(inputFile)
	validWords := removeStopwords(words)
	freqs := countFrequencies(validWords)
	sortedFreqs := sortFrequencies(freqs)
	printFrequencies(sortedFreqs[0:25])
}

func extractWords(inputFile string) []string {
	b, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	pattern := regexp.MustCompile(`[a-z]{2,}`)
	words := pattern.FindAllString(strings.ToLower(string(b)), -1)
	return words
}

func removeStopwords(words []string) []string {
	b, err := os.ReadFile(stopwordsFile)
	if err != nil {
		panic(err)
	}
	stopwords := strings.Split(string(b), ",")
	stopwordMap := map[string]struct{}{}
	for _, stopword := range stopwords {
		if _, ok := stopwordMap[stopword]; !ok {
			stopwordMap[stopword] = struct{}{}
		}
	}

	validWords := []string{}
	for _, word := range words {
		if _, ok := stopwordMap[word]; !ok && len(word) >= 2 {
			validWords = append(validWords, word)
		}
	}

	return validWords
}

func countFrequencies(words []string) []Freq {
	wordFreqMap := map[string]int{}
	for _, word := range words {
		wordFreqMap[word]++
	}

	var wordFreqs []Freq
	for word, count := range wordFreqMap {
		wordFreqs = append(wordFreqs, Freq{
			word:  word,
			count: count,
		})
	}

	return wordFreqs
}

func sortFrequencies(freqs []Freq) []Freq {
	sort.Slice(freqs, func(i, j int) bool {
		return freqs[i].count >= freqs[j].count
	})
	return freqs
}

func printFrequencies(freqs []Freq) {
	for _, freq := range freqs {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}
