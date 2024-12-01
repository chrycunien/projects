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
	stdout = os.Stdout
	if len(os.Args) < 2 {
		fmt.Fprintf(stderr, "%v\b", ErrNotEnoughArgument)
		os.Exit(1)
	}
	inputFile = os.Args[1]
	run()
}

func run() {
	words, err := extractWords(inputFile)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		os.Exit(1)
	}

	validWords, err := removeStopwords(words)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		os.Exit(1)
	}

	freqs := countFrequencies(validWords)
	sortedFreqs := sortFrequencies(freqs)
	if len(sortedFreqs) < 25 {
		fmt.Fprintln(stderr, "the file contains less than 25 different words")
		os.Exit(1)
	}

	printFrequencies(sortedFreqs[0:25])
}

func extractWords(inputFile string) ([]string, error) {
	b, err := os.ReadFile(inputFile)
	if err != nil {
		return []string{}, err
	}

	pattern := regexp.MustCompile(`[a-z]{2,}`)
	words := pattern.FindAllString(strings.ToLower(string(b)), -1)
	return words, nil
}

func removeStopwords(words []string) ([]string, error) {
	b, err := os.ReadFile(stopwordsFile)
	if err != nil {
		return []string{}, err
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

	return validWords, nil
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
