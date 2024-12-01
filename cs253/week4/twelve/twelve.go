package main

import (
	"errors"
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

var ErrUnknownMethod = errors.New("unknown method")

type Freq struct {
	word  string
	count int
}

type Message struct {
	action string
	data   any
}

func NewMessage(action string, data any) Message {
	return Message{
		action: action,
		data:   data,
	}
}

type DataStorageManager struct {
	_data     string
	_words    []string
	processed bool
}

func (dsm *DataStorageManager) Dispatch(m Message) (any, error) {
	switch m.action {
	case "init":
		return dsm.init(m.data.(string))
	case "words":
		return dsm.words()
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownMethod, m.action)
	}
}

func (dsm *DataStorageManager) init(inputFile string) (any, error) {
	wordsBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}
	dsm._data = strings.ToLower(string(wordsBytes))
	return nil, nil
}

func (dsm *DataStorageManager) words() (any, error) {
	if !dsm.processed {
		pattern := regexp.MustCompile(`[a-z]{2,}`)
		dsm._words = pattern.FindAllString(dsm._data, -1)
		dsm.processed = true
	}
	return dsm._words, nil
}

type StopwordManager struct {
	stopwords map[string]struct{}
}

func (sm *StopwordManager) Dispatch(m Message) (any, error) {
	switch m.action {
	case "init":
		return sm.init(m.data.(string))
	case "is_stop_word":
		return sm.isStopword(m.data.(string))
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownMethod, m.action)
	}
}

func (sm *StopwordManager) init(stopwordFile string) (any, error) {
	stopwordsBytes, err := os.ReadFile(stopwordFile)
	if err != nil {
		return nil, err
	}

	stopwords := strings.Split(string(stopwordsBytes), ",")
	stopwordMap := map[string]struct{}{}
	for _, stopword := range stopwords {
		stopwordMap[stopword] = struct{}{}
	}

	sm.stopwords = stopwordMap
	return nil, nil
}

func (sm *StopwordManager) isStopword(word string) (any, error) {
	_, ok := sm.stopwords[word]
	return ok, nil
}

type WordFrequencyManager struct {
	freqs map[string]int
}

func (wfm *WordFrequencyManager) Dispatch(m Message) (any, error) {
	switch m.action {
	case "init":
		return wfm.init()
	case "increment_count":
		return wfm.incrementCount(m.data.(string))
	case "sorted":
		return wfm.sorted()
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownMethod, m.action)
	}
}

func (wfm *WordFrequencyManager) init() (any, error) {
	wfm.freqs = map[string]int{}
	return nil, nil
}

func (wfm *WordFrequencyManager) incrementCount(word string) (any, error) {
	wfm.freqs[word]++
	return nil, nil
}

func (wfm *WordFrequencyManager) sorted() (any, error) {
	wordFreqs := []Freq{}
	for k, v := range wfm.freqs {
		wordFreqs = append(wordFreqs, Freq{k, v})
	}

	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count >= wordFreqs[j].count
	})

	return wordFreqs, nil
}

type WordFrequencyController struct {
	dsm *DataStorageManager
	sm  *StopwordManager
	wfm *WordFrequencyManager
}

func (wfc *WordFrequencyController) Dispatch(m Message) (any, error) {
	switch m.action {
	case "init":
		files := m.data.([]string)
		return wfc.init(files[0], files[1])
	case "run":
		return wfc.run()
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownMethod, m.action)
	}
}

func (wfc *WordFrequencyController) init(inputFile, stopwordFile string) (any, error) {
	wfc.dsm = &DataStorageManager{}
	wfc.sm = &StopwordManager{}
	wfc.wfm = &WordFrequencyManager{}

	var err error
	_, err = wfc.dsm.Dispatch(NewMessage("init", inputFile))
	if err != nil {
		return nil, err
	}

	_, err = wfc.sm.Dispatch(NewMessage("init", stopwordFile))
	if err != nil {
		return nil, err
	}

	_, err = wfc.wfm.Dispatch(NewMessage("init", nil))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (wfc *WordFrequencyController) run() (any, error) {
	var err error

	words, err := wfc.dsm.Dispatch(NewMessage("words", nil))
	if err != nil {
		return nil, err
	}
	for _, w := range words.([]string) {
		isStopword, err := wfc.sm.Dispatch(NewMessage("is_stop_word", w))
		if err != nil {
			return nil, err
		}
		if !isStopword.(bool) {
			_, err = wfc.wfm.Dispatch(NewMessage("increment_count", w))
			if err != nil {
				return nil, err
			}
		}
	}

	sorted, err := wfc.wfm.Dispatch(NewMessage("sorted", nil))
	if err != nil {
		return nil, err
	}
	for _, freq := range sorted.([]Freq)[0:25] {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}

  return nil, nil
}

func main() {
	stdout = os.Stdout
	inputFile = os.Args[1]
	stopwordsFile = "../../stop_words.txt"

	var err error
	wfc := &WordFrequencyController{}

	_, err = wfc.Dispatch(NewMessage("init", []string{inputFile, stopwordsFile}))
	if err != nil {
		os.Exit(1)
	}

	_, err = wfc.Dispatch(NewMessage("run", nil))
	if err != nil {
		os.Exit(1)
	}
}
