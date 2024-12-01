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
var stopwordFile string
var stdout io.Writer
var stderr io.Writer

type Freq struct {
	word  string
	count int
}
type PlainFunc = func()
type StringFunc = func(string)

type WordFrequencyFramework struct {
	loadEventHandlers   []StringFunc
	doWorkEventHandlers []PlainFunc
	endEventHandlers    []PlainFunc
}

func (wff *WordFrequencyFramework) RegisterForLoadEvent(f StringFunc) {
	wff.loadEventHandlers = append(wff.loadEventHandlers, f)
}

func (wff *WordFrequencyFramework) RegisterForDoWorkEvent(f PlainFunc) {
	wff.doWorkEventHandlers = append(wff.doWorkEventHandlers, f)
}

func (wff *WordFrequencyFramework) RegisterForEndEvent(f PlainFunc) {
	wff.endEventHandlers = append(wff.endEventHandlers, f)
}

func (wff *WordFrequencyFramework) Run(input string) {
	for _, h := range wff.loadEventHandlers {
		h(input)
	}
	for _, h := range wff.doWorkEventHandlers {
		h()
	}
	for _, h := range wff.endEventHandlers {
		h()
	}
}

type DataStorage struct {
	filter        *StopwordFilter
	data          []string
	eventHandlers []StringFunc
}

func NewDataStorage(wff *WordFrequencyFramework, filter *StopwordFilter) *DataStorage {
	ds := &DataStorage{
		filter: filter,
	}
	wff.RegisterForLoadEvent(func(inputFile string) {
		ds.load(inputFile)
	})
	wff.RegisterForDoWorkEvent(func() {
		ds.produceWords()
	})
	return ds
}

func (ds *DataStorage) load(inputFile string) {
	wordsBytes, err := os.ReadFile(inputFile)
	if err != nil {
		os.Exit(1)
	}

	pattern := regexp.MustCompile(`[a-z]{2,}`)
	ds.data = pattern.FindAllString(strings.ToLower(string(wordsBytes)), -1)
}

func (ds *DataStorage) produceWords() {
	for _, w := range ds.data {
		if !ds.filter.IsStopword(w) {
			for _, h := range ds.eventHandlers {
				h(w)
			}
		}
	}
}

func (ds *DataStorage) RegisterEventHandler(f StringFunc) {
	ds.eventHandlers = append(ds.eventHandlers, f)
}

type StopwordFilter struct {
	stopwords map[string]struct{}
}

func NewStopwordFilter(wff *WordFrequencyFramework, stopwordFile string) *StopwordFilter {
	sf := &StopwordFilter{
		stopwords: map[string]struct{}{},
	}
	wff.RegisterForLoadEvent(func(_ string) {
		sf.load(stopwordFile)
	})
	return sf
}

func (sf *StopwordFilter) load(stopwordFile string) {
	stopwordsBytes, err := os.ReadFile(stopwordFile)
	if err != nil {
		os.Exit(1)
	}

	stopwords := strings.Split(string(stopwordsBytes), ",")
	stopwordMap := map[string]struct{}{}
	for _, stopword := range stopwords {
		stopwordMap[stopword] = struct{}{}
	}

	sf.stopwords = stopwordMap
}

func (sf *StopwordFilter) IsStopword(word string) bool {
	_, ok := sf.stopwords[word]
	return ok
}

type WordFrequencyCounter struct {
	freqs map[string]int
}

func NewWordFrequencyCounter(wff *WordFrequencyFramework, ds *DataStorage) *WordFrequencyCounter {
	wfc := &WordFrequencyCounter{
		freqs: map[string]int{},
	}
	ds.RegisterEventHandler(func(word string) {
		wfc.incrementCount(word)
	})
	wff.RegisterForEndEvent(func() {
		wfc.printFreqs()
	})
	return wfc
}

func (wfc *WordFrequencyCounter) incrementCount(word string) {
	wfc.freqs[word]++
}

func (wfc *WordFrequencyCounter) printFreqs() {
	wordFreqs := []Freq{}
	for k, v := range wfc.freqs {
		wordFreqs = append(wordFreqs, Freq{k, v})
	}

	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count >= wordFreqs[j].count
	})

	for _, freq := range wordFreqs[0:25] {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}

type WordWithZPrinter struct {
	matched []string
}

func NewWordWithZPrinter(wff *WordFrequencyFramework, ds *DataStorage) *WordWithZPrinter {
	wzp := &WordWithZPrinter{}
	ds.RegisterEventHandler(func(word string) {
		wzp.countWord(word)
	})
	wff.RegisterForEndEvent(func() {
		wzp.printZWords()
	})
	return wzp
}

func (wzp *WordWithZPrinter) countWord(word string) {
	pattern := regexp.MustCompile(`z|Z`)
	if pattern.MatchString(word) {
		wzp.matched = append(wzp.matched, word)
	}
}

func (wzp *WordWithZPrinter) printZWords() {
	switch len(wzp.matched) {
	case 0:
		fmt.Fprintln(stderr, "File contains no word with z.")
	case 1:
		fmt.Fprintln(stderr, "File contains: 1 word with z.")
	default:
		fmt.Fprintf(stderr, "File contains: %d words with z.\n", len(wzp.matched))
	}
}

func main() {
	stdout = os.Stdout
	stderr = os.Stderr
	inputFile = os.Args[1]
	stopwordFile = "../../stop_words.txt"
	run()
}

func run() {
	wff := &WordFrequencyFramework{}
	stopwordFilter := NewStopwordFilter(wff, stopwordFile)
	ds := NewDataStorage(wff, stopwordFilter)
	_ = NewWordFrequencyCounter(wff, ds)
	// Exercise 15.2
	_ = NewWordWithZPrinter(wff, ds)
	wff.Run(inputFile)
}
