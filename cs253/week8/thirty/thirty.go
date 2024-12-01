package main

import (
	"fmt"
	"io"
	"maps"
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

func processWords(inputFile string, out chan<- string) {
	wordsBytes, err := os.ReadFile(inputFile)
	if err != nil {
		os.Exit(1)

	}
	pattern := regexp.MustCompile(`[a-z]{2,}`)
	words := pattern.FindAllString(strings.ToLower(string(wordsBytes)), -1)
	for _, word := range words {
		out <- word
	}

	close(out)
}

func countFreq(
	stopwordsMap map[string]struct{},
	input <-chan string,
	output chan<- map[string]int,
) {
	wordFreq := map[string]int{}
	for word := range input {
		if _, ok := stopwordsMap[word]; !ok {
			wordFreq[word]++
		}
	}
	output <- wordFreq
}

func main() {
	stdout = os.Stdout
	stderr = os.Stderr
	inputFile = os.Args[1]
	run()
}

func run() {
	wordSpace := make(chan string)
	freqSpace := make(chan map[string]int)

	go processWords(inputFile, wordSpace)

	b, err := os.ReadFile(stopwordsFile)
	if err != nil {
		os.Exit(1)
	}

	stopwords := strings.Split(string(b), ",")
	stopwordsMap := map[string]struct{}{}
	for _, stopWord := range stopwords {
		stopwordsMap[stopWord] = struct{}{}
	}

	workers := 4
	for i := 0; i < workers; i++ {
		go countFreq(maps.Clone(stopwordsMap), wordSpace, freqSpace)
	}

	wordFreq := map[string]int{}
	for i := 0; i < workers; i++ {
		freqMap := <-freqSpace
		for word, count := range freqMap {
			wordFreq[word] += count
		}
	}

	freqs := []Freq{}
	for k, v := range wordFreq {
		freqs = append(freqs, Freq{
			word:  k,
			count: v,
		})
	}
	sort.Slice(freqs, func(i, j int) bool {
		return freqs[i].count >= freqs[j].count
	})

	for _, freq := range freqs[0:25] {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}
