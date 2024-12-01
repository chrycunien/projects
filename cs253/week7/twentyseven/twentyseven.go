package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
)

type Column interface {
	Update()
}

type Freq struct {
	word  string
	count int
}

type AllWords struct {
	Words []string
}

func NewAllWords() *AllWords {
	return &AllWords{}
}

func (aw *AllWords) Update() {
	// no-op
}

type StopWords struct {
	Words []string
}

func NewStopWords() *StopWords {
	return &StopWords{}
}

func (sw *StopWords) Update() {
	// no-op
}

type NonStopWords struct {
	Words []string
}

func NewNonStopWords() *NonStopWords {
	return &NonStopWords{}
}

func (nsw *NonStopWords) Update() {
	stopWordMap := map[string]struct{}{}
	for _, stopWord := range stopWords.Words {
		stopWordMap[stopWord] = struct{}{}
	}

	words := []string{}
	for _, word := range allWords.Words {
		if _, ok := stopWordMap[word]; !ok {
			words = append(words, word)
		} else {
			words = append(words, "")
		}
	}

	nsw.Words = words
}

type UniqueWords struct {
	WordMap map[string]struct{}
}

func NewUniqueWords() *UniqueWords {
	return &UniqueWords{
		WordMap: map[string]struct{}{},
	}
}

func (uw *UniqueWords) Update() {
	wordMap := map[string]struct{}{}
	for _, word := range nonStopWords.Words {
		if word != "" {
			wordMap[word] = struct{}{}
		}
	}

	uw.WordMap = wordMap
}

type CountWords struct {
	Freqs []Freq
}

func NewCountWords() *CountWords {
	return &CountWords{}
}

func (cw *CountWords) Update() {
	counts := map[string]int{}
	for _, word := range nonStopWords.Words {
		if _, ok := uniqueWords.WordMap[word]; ok {
			counts[word] += 1
		}
	}

	freqs := []Freq{}
	for k, v := range counts {
		freqs = append(freqs, Freq{
			word:  k,
			count: v,
		})
	}

	cw.Freqs = freqs
}

type SortedWords struct {
	Freqs []Freq
}

func NewSortedWords() *SortedWords {
	return &SortedWords{}
}

func (sw *SortedWords) Update() {
	freqs := slices.Clone(countWords.Freqs)
	sort.Slice(freqs, func(i, j int) bool {
		return freqs[i].count >= freqs[j].count
	})
	sw.Freqs = freqs
}

var stdout io.Writer
var stderr io.Writer
var inputFile string
var stopwordsFile string = "../../stop_words.txt"

var allWords *AllWords = NewAllWords()
var stopWords *StopWords = NewStopWords()
var nonStopWords *NonStopWords = NewNonStopWords()
var uniqueWords *UniqueWords = NewUniqueWords()
var countWords *CountWords = NewCountWords()
var sortedWords *SortedWords = NewSortedWords()

var columns []Column = []Column{
	allWords,
	stopWords,
	nonStopWords,
	uniqueWords,
	countWords,
	sortedWords,
}

func update() {
	for _, col := range columns {
		col.Update()
	}
}

func main() {
	stdout = os.Stdout
	stderr = os.Stderr
	// inputFile = os.Args[1]

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter file name:")
		if !scanner.Scan() {
			break
		}
		
		inputFile = scanner.Text()
		if inputFile == "" {
			fmt.Fprintln(stderr, "empty file name: quit")
			break
		}
		run(inputFile)
	}
}

func run(inputFile string) {
	b, err := os.ReadFile(inputFile)
	if err != nil {
		os.Exit(1)
	}
	pattern := regexp.MustCompile(`[a-z]{2,}`)
	allWords.Words = pattern.FindAllString(strings.ToLower(string(b)), -1)

	b, err = os.ReadFile(stopwordsFile)
	if err != nil {
		os.Exit(1)
	}
	stopWords.Words = strings.Split(string(b), ",")

	update()

	for _, freq := range sortedWords.Freqs[0:25] {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}

func countFrequencies(words any) any {
	wordFreqMap := map[string]int{}
	wordList := words.([]string)
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
	return wordFreqs
}

func sortFrequencies(v any) any {
	freqs := v.([]Freq)
	sort.Slice(freqs, func(i, j int) bool {
		return freqs[i].count >= freqs[j].count
	})
	return freqs
}
