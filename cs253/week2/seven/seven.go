package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
)

type Freq struct {
	word  string
	count int
}

var stdout io.Writer = os.Stdout
var inputFile string = os.Args[1]

func must[T any](x T, _ error) T {
	return x
}

func main() {
	run()
}

func run() {
	stopwords, words, wordFreqMap, wordFreqs := strings.Split(string(must(os.ReadFile("../../stop_words.txt"))), ","), regexp.MustCompile(`[a-z]{2,}`).FindAllString(strings.ToLower(string(must(os.ReadFile(inputFile)))), -1), map[string]int{}, []Freq{}
	for _, word := range words {
		if len(word) >= 2 && !slices.Contains(stopwords, word) {
			wordFreqMap[word]++
		}
	}
	for k, v := range wordFreqMap {
		wordFreqs = append(wordFreqs, Freq{k, v})
	}
	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count >= wordFreqs[j].count
	})
	for _, freq := range wordFreqs[0:25] {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}
