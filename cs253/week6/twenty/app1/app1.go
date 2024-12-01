package main

import (
	"os"
	"regexp"
	"sort"
	"strings"

	"main/week6/twenty/commons"
)

type DataStorageManager struct {
	data      string
	words     []string
	processed bool
}

func NewDataStorageManager(inputFile string) commons.DataStorageManager {
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

func NewStopwordManager(stopwordFile string) commons.StopwordManager {
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

func NewWordFrequencyManager() commons.WordFrequencyManager {
	return &WordFrequencyManager{
		freqs: map[string]int{},
	}
}

func (wfm *WordFrequencyManager) IncrementCount(word string) {
	wfm.freqs[word]++
}

func (wfm *WordFrequencyManager) Sorted() []commons.Freq {
	wordFreqs := []commons.Freq{}
	for k, v := range wfm.freqs {
		wordFreqs = append(wordFreqs, commons.Freq{k, v})
	}

	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].Count >= wordFreqs[j].Count
	})

	return wordFreqs
}
