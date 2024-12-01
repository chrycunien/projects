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
	ErrTypeMismatched    = errors.New("mismatched type")
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

	freqs, err := countFrequencies(validWords)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		os.Exit(1)
	}

	sortedFreqs, err := sortFrequencies(freqs)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		os.Exit(1)
	}

	if len(sortedFreqs) < 25 {
		fmt.Fprintln(stderr, "the file contains less than 25 different words")
		os.Exit(1)
	}

	err = printFrequencies(sortedFreqs[0:25])
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		os.Exit(1)
	}
}

func extractWords(inputFile any) ([]string, error) {
	inputFileStr, ok := inputFile.(string)
	if !ok {
		return []string{}, ErrTypeMismatched
	}
	b, err := os.ReadFile(inputFileStr)
	if err != nil {
		return []string{}, err
	}

	pattern := regexp.MustCompile(`[a-z]{2,}`)
	words := pattern.FindAllString(strings.ToLower(string(b)), -1)
	return words, nil
}

func removeStopwords(words any) ([]string, error) {
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

	wordList, ok := words.([]string)
	if !ok {
		return []string{}, ErrTypeMismatched
	}
	validWords := []string{}
	for _, word := range wordList {
		if _, ok := stopwordMap[word]; !ok && len(word) >= 2 {
			validWords = append(validWords, word)
		}
	}

	return validWords, nil
}

func countFrequencies(words any) ([]Freq, error) {
	wordList, ok := words.([]string)
	if !ok {
		return []Freq{}, ErrTypeMismatched
	}
	wordFreqMap := map[string]int{}
	for _, word := range wordList {
		wordFreqMap[word]++
	}

	var wordFreqs []Freq
	for word, count := range wordFreqMap {
		wordFreqs = append(wordFreqs, Freq{
			word:  word,
			count: count,
		})
	}

	return wordFreqs, nil
}

func sortFrequencies(freqs any) ([]Freq, error) {
	freqList, ok := freqs.([]Freq)
	if !ok {
		return []Freq{}, ErrTypeMismatched
	}
	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].count >= freqList[j].count
	})
	return freqList, nil
}

func printFrequencies(freqs any) error {
	freqList, ok := freqs.([]Freq)
	if !ok {
		return ErrTypeMismatched
	}

	for _, freq := range freqList {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
	return nil
}
