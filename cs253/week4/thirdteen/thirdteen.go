package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

var inputFile string
var stopwordsFile string
var stdout io.Writer

type Freq struct {
	word  string
	count int
}

type PlainFunc = func()
type InitFunc = func(string)
type StringSliceFunc = func() []string
type StringPredicate = func(string) bool
type StringFunc = func(string)
type WordFreqSliceFunc = func() []Freq

type Object map[string]any
type FreqMap map[string]int
type StopwordMap map[string]struct{}

func main() {
	stdout = os.Stdout
	inputFile = os.Args[1]
	stopwordsFile = "../../stop_words.txt"
	run()
}

func run() {
	dataStorageObj := Object{}
	dataStorageObj["data"] = nil
	dataStorageObj["init"] = func(inputFile string) {
		extractWords(dataStorageObj, inputFile)
	}
	// Exercise 13.3
	dataStorageObj["words"] = func(this Object) StringSliceFunc {
		return func() []string {
			return this["data"].([]string)
		}
	}(dataStorageObj)

	stopwordsObj := Object{}
	stopwordsObj["stopwords"] = nil
	stopwordsObj["init"] = func(stopwordFile string) {
		loadStopwords(stopwordsObj, stopwordFile)
	}
	stopwordsObj["is_stop_word"] = func(word string) bool {
		_, ok := stopwordsObj["stopwords"].(StopwordMap)[word]
		return ok
	}

	wordFreqObj := Object{}
	wordFreqObj["freqs"] = FreqMap{}
	wordFreqObj["increment_count"] = func(word string) {
		incrementCount(wordFreqObj, word)
	}
	wordFreqObj["sorted"] = func() []Freq {
		wordFreqs := []Freq{}
		for k, v := range wordFreqObj["freqs"].(FreqMap) {
			wordFreqs = append(wordFreqs, Freq{k, v})
		}

		sort.Slice(wordFreqs, func(i, j int) bool {
			return wordFreqs[i].count >= wordFreqs[j].count
		})

		return wordFreqs
	}

	dataStorageObj["init"].(InitFunc)(inputFile)
	stopwordsObj["init"].(InitFunc)(stopwordsFile)

	for _, w := range dataStorageObj["words"].(StringSliceFunc)() {
		if !stopwordsObj["is_stop_word"].(StringPredicate)(w) {
			wordFreqObj["increment_count"].(StringFunc)(w)
		}
	}

	// Exercise 13.2
	wordFreqObj["top25"] = func() {
		for _, freq := range wordFreqObj["sorted"].(WordFreqSliceFunc)()[0:25] {
			fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
		}
	}

	wordFreqObj["top25"].(PlainFunc)()
}

func extractWords(obj Object, inputFile string) {
	var err error
	obj["data"], err = os.ReadFile(inputFile)
	if err != nil {
		os.Exit(1)
	}

	pattern := regexp.MustCompile(`[a-z]{2,}`)
	obj["data"] = pattern.FindAllString(strings.ToLower(string(obj["data"].([]byte))), -1)
}

func loadStopwords(obj Object, stopwordFile string) {
	var err error
	obj["stopwords"], err = os.ReadFile(stopwordFile)
	if err != nil {
		os.Exit(1)
	}

	stopwords := strings.Split(string(obj["stopwords"].([]byte)), ",")
	obj["stopwords"] = StopwordMap{}
	for _, stopword := range stopwords {
		obj["stopwords"].(StopwordMap)[stopword] = struct{}{}
	}
}

func incrementCount(obj Object, word string) {
	obj["freqs"].(FreqMap)[word]++
}
