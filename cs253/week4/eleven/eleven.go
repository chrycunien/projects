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

type DataStorageManager struct {
	data      string
	words     []string
	processed bool
}

func NewDataStorageManager(inputFile string) *DataStorageManager {
	wordsBytes, err := os.ReadFile(inputFile)
	if err != nil {
		os.Exit(1)
	}
	return &DataStorageManager{
		data: strings.ToLower(string(wordsBytes)),
	}
}

func (dsm *DataStorageManager) Words() []string {
	if !dsm.processed {
		pattern := regexp.MustCompile(`[a-z]{2,}`)
		dsm.words = pattern.FindAllString(dsm.data, -1)
		dsm.processed = true
	}
	return dsm.words
}

type StopwordManager struct {
	stopwords map[string]struct{}
}

func NewStopwordManager(stopwordFile string) *StopwordManager {
	stopwordsBytes, err := os.ReadFile(stopwordFile)
	if err != nil {
		os.Exit(1)
	}

	stopwords := strings.Split(string(stopwordsBytes), ",")
	stopwordMap := map[string]struct{}{}
	for _, stopword := range stopwords {
		stopwordMap[stopword] = struct{}{}
	}

	return &StopwordManager{
		stopwords: stopwordMap,
	}
}

func (sm *StopwordManager) IsStopword(word string) bool {
	_, ok := sm.stopwords[word]
	return ok
}

type WordFrequencyManager struct {
	freqs map[string]int
}

func NewWordFrequencyManager() *WordFrequencyManager {
	return &WordFrequencyManager{
		freqs: map[string]int{},
	}
}

func (wfm *WordFrequencyManager) IncrementCount(word string) {
	wfm.freqs[word]++
}

func (wfm *WordFrequencyManager) Sorted() []Freq {
	wordFreqs := []Freq{}
	for k, v := range wfm.freqs {
		wordFreqs = append(wordFreqs, Freq{k, v})
	}

	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count >= wordFreqs[j].count
	})

	return wordFreqs
}

type WordFrequencyController struct {
	dsm *DataStorageManager
	sm  *StopwordManager
	wfm *WordFrequencyManager
}

func NewWordFrequencyController(inputFile, stopwordFile string) *WordFrequencyController {
	return &WordFrequencyController{
		dsm: NewDataStorageManager(inputFile),
		sm:  NewStopwordManager(stopwordsFile),
		wfm: NewWordFrequencyManager(),
	}
}

func (wfc *WordFrequencyController) Run() {
	for _, w := range wfc.dsm.Words() {
		if !wfc.sm.IsStopword(w) {
			wfc.wfm.IncrementCount(w)
		}
	}

	for _, freq := range wfc.wfm.Sorted()[0:25] {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}

func main() {
	stdout = os.Stdout
	inputFile = os.Args[1]
	stopwordsFile = "../../stop_words.txt"
	wfc := NewWordFrequencyController(inputFile, stopwordsFile)
	wfc.Run()
}
