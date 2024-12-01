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

type TFTheOne struct {
	v any
}

func (ft *TFTheOne) Bind(f func(any) any) *TFTheOne {
	ft.v = f(ft.v)
	return ft
}

func NewTFTheOne(v any) *TFTheOne {
	return &TFTheOne{
		v: v,
	}
}

var stdout io.Writer
var inputFile string
var stopwordsFile string = "../../stop_words.txt"

func main() {
	stdout = os.Stdout
	inputFile = os.Args[1]
	run()
}

func run() {
	ft := NewTFTheOne(inputFile)
	ft.
    Bind(readFile).
    Bind(filterCharAndNormalize).
    Bind(removeStopwords).
    Bind(countWords).
    Bind(flattenWordFreqs).
    Bind(sortWordFreqs).
    Bind(printWordFreqs)
}

func readFile(inputFile any) any {
	b, err := os.ReadFile(inputFile.(string))
	if err != nil {
		os.Exit(1)
	}
	return string(b)
}

func filterCharAndNormalize(wordStr any) any {
	pattern := regexp.MustCompile(`[a-z]{2,}`)
	words := pattern.FindAllString(strings.ToLower(wordStr.(string)), -1)
	return words
}

func removeStopwords(words any) any {
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
	for _, word := range words.([]string) {
		if _, ok := stopwordMap[word]; !ok && len(word) >= 2 {
			validWords = append(validWords, word)
		}
	}

	return validWords
}

func countWords(words any) any {
	wordFreqMap := map[string]int{}
	for _, word := range words.([]string) {
		wordFreqMap[word]++
	}
	return wordFreqMap
}

func flattenWordFreqs(wordFreqMap any) any {
	var wordFreqs []Freq
	for word, count := range wordFreqMap.(map[string]int) {
		wordFreqs = append(wordFreqs, Freq{
			word:  word,
			count: count,
		})
	}
	return wordFreqs
}

func sortWordFreqs(v any) any {
	wordFreqs := v.([]Freq)
	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count >= wordFreqs[j].count
	})
	return wordFreqs
}

func printWordFreqs(v any) any {
	wordFreqs := v.([]Freq)
	for i := 0; i < 25; i++ {
		freq := wordFreqs[i]
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
	return nil
}
