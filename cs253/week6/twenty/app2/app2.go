package main

import (
	"os"
	"slices"
	"strings"

	"main/week6/twenty/commons"
)

const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
		data: string(wordsBytes),
	}
}

func (dsm *DataStorageManager) Words() []string {
	if !dsm.processed {
		dsm.words = scanWords(filterCharAndNormalize(dsm.data))
		dsm.processed = true
	}
	return dsm.words
}

func filterCharAndNormalize(wordStr string) string {
	s := make([]byte, len(wordStr))
	for i := 0; i < len(wordStr); i++ {
		c := wordStr[i]
		if !isAlnum(c) {
			s[i] = ' '
		} else {
			s[i] = lower(c)
		}
	}
	return string(s)
}

func isAlnum(c byte) bool {
	for i := 0; i < len(alnum); i++ {
		if c == alnum[i] {
			return true
		}
	}
	return false
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		c = c - 'A' + 'a'
	}
	return c
}

func scanWords(wordStr string) []string {
	return strings.Fields(wordStr)
}

type StopwordManager struct {
	stopwords []string
}

func NewStopwordManager(stopwordFile string) commons.StopwordManager {
	stopwordsBytes, err := os.ReadFile(stopwordFile)
	if err != nil {
		os.Exit(1)
	}

	stopwords := strings.Split(string(stopwordsBytes), ",")
	return &StopwordManager{
		stopwords: stopwords,
	}
}

func (sm *StopwordManager) IsStopword(word string) bool {
	return len(word) < 2 || slices.Contains(sm.stopwords, word)
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
		wordFreqs = append(wordFreqs, commons.Freq{
			Word:  k,
			Count: v,
		})
	}

	sortWordFreqs(wordFreqs)
	return wordFreqs
}

func sortWordFreqs(wordFreqs []commons.Freq) {
	if len(wordFreqs) <= 1 {
		return
	}

	sortWordFreqs(wordFreqs[1:])

	// bubble sort
	for i := 0; i < len(wordFreqs)-1; i++ {
		if wordFreqs[i].Count >= wordFreqs[i+1].Count {
			break
		}
		wordFreqs[i], wordFreqs[i+1] = wordFreqs[i+1], wordFreqs[i]
	}
}
