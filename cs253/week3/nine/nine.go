package main

import (
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

type ReadFileFunc func(string, FilterAndNormalizeFunc)
type FilterAndNormalizeFunc func(string, RemoveStopwordsFunc)
type RemoveStopwordsFunc func([]string, CountFunc)
type CountFunc func([]string, FlattenFunc)
type FlattenFunc func(map[string]int, SortFunc)
type SortFunc func([]Freq, PrintFreqFunc)
type PrintFreqFunc func([]Freq, NoOpFunc)
type NoOpFunc func()

var stdout io.Writer
var inputFile string
var stopwordsFile string = "../../stop_words.txt"

func main() {
	stdout = os.Stdout
	inputFile = os.Args[1]
	run()
}

func run() {
	readFile(inputFile, filterCharAndNormalize)
}

func readFile(inputFile string, f FilterAndNormalizeFunc) {
	b, err := os.ReadFile(inputFile)
	if err != nil {
		os.Exit(1)
	}
	f(string(b), removeStopwords)
}

func filterCharAndNormalize(wordStr string, f RemoveStopwordsFunc) {
	pattern := regexp.MustCompile(`[a-z]{2,}`)
	words := pattern.FindAllString(strings.ToLower(wordStr), -1)
	f(words, countWords)
}

func removeStopwords(words []string, f CountFunc) {
	b, err := os.ReadFile(stopwordsFile)
	if err != nil {
		os.Exit(1)
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

	f(validWords, flattenWordFreq)
}

func countWords(words []string, f FlattenFunc) {
	wordFreqMap := map[string]int{}
	for _, word := range words {
		wordFreqMap[word]++
	}

	f(wordFreqMap, sortWordFreqs)
}

func flattenWordFreq(wordFreqMap map[string]int, f SortFunc) {
	var wordFreqs []Freq
	for word, count := range wordFreqMap {
		wordFreqs = append(wordFreqs, Freq{
			word:  word,
			count: count,
		})
	}

	f(wordFreqs, printWordFreqs)
}

func sortWordFreqs(wordFreqs []Freq, f PrintFreqFunc) {
	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count >= wordFreqs[j].count
	})

	f(wordFreqs, func() {})
}

func printWordFreqs(wordFreqs []Freq, f NoOpFunc) {
	for i := 0; i < 25; i++ {
		freq := wordFreqs[i]
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}

	f()
}
