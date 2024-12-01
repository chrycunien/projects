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

type EventType string

const (
	EventTypeLoad      EventType = "load"
	EventTypeRun       EventType = "run"
	EventTypeWord      EventType = "word"
	EventTypeValidWord EventType = "valid_word"
	EventTypePrint     EventType = "print"
	EventTypeEof       EventType = "eof"
	EventTypeStart     EventType = "start"
	EventTypeStop      EventType = "stop"
)

type Event struct {
	eventType EventType
	data      any
}

func NewEvent(eventType EventType, data any) Event {
	return Event{
		eventType: eventType,
		data:      data,
	}
}

type EventHandler func(Event)

type EventManager struct {
	subscriptions map[EventType][]EventHandler
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscriptions: map[EventType][]EventHandler{},
	}
}

func (em *EventManager) Subscribe(etype EventType, handler EventHandler) {
	em.subscriptions[etype] = append(em.subscriptions[etype], handler)
}

func (em *EventManager) Publish(e Event) {
	if handlers, ok := em.subscriptions[e.eventType]; ok {
		for _, h := range handlers {
			h(e)
		}
	}
}

type DataStorage struct {
	eventManager *EventManager
	data         []string
}

func NewDataStorage(em *EventManager) *DataStorage {
	ds := &DataStorage{
		eventManager: em,
	}
	ds.eventManager.Subscribe(EventTypeLoad, func(e Event) {
		ds.load(e.data.(string))
	})
	ds.eventManager.Subscribe(EventTypeRun, func(e Event) {
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
		ds.eventManager.Publish(NewEvent(EventTypeWord, w))
	}
	ds.eventManager.Publish(NewEvent(EventTypeEof, nil))
}

type StopwordFilter struct {
	eventManager *EventManager
	stopwords    map[string]struct{}
}

func NewStopwordFilter(em *EventManager, stopwordFile string) *StopwordFilter {
	sf := &StopwordFilter{
		eventManager: em,
		stopwords:    map[string]struct{}{},
	}
	em.Subscribe(EventTypeLoad, func(e Event) {
		sf.load(stopwordFile)
	})
	em.Subscribe(EventTypeWord, func(e Event) {
		sf.isStopword(e.data.(string))
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

func (sf *StopwordFilter) isStopword(word string) {
	if _, ok := sf.stopwords[word]; !ok {
		sf.eventManager.Publish(NewEvent(EventTypeValidWord, word))
	}
}

type WordFrequencyCounter struct {
	freqs map[string]int
}

func NewWordFrequencyCounter(em *EventManager) *WordFrequencyCounter {
	wfc := &WordFrequencyCounter{
		freqs: map[string]int{},
	}
	em.Subscribe(EventTypeValidWord, func(e Event) {
		wfc.incrementCount(e.data.(string))
	})
	em.Subscribe(EventTypePrint, func(e Event) {
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

type WordFrequencyApplication struct {
	eventManager *EventManager
}

func NewWordFrequencyApplication(em *EventManager) *WordFrequencyApplication {
	wfa := &WordFrequencyApplication{
		eventManager: em,
	}
	em.Subscribe(EventTypeStart, func(e Event) {
		wfa.start(e.data.(string))
	})
	em.Subscribe(EventTypeEof, func(e Event) {
		wfa.stop()
	})
	return wfa
}

func (wfa *WordFrequencyApplication) start(inputFile string) {
	wfa.eventManager.Publish(NewEvent(EventTypeLoad, inputFile))
	wfa.eventManager.Publish(NewEvent(EventTypeRun, nil))
}

func (wfa *WordFrequencyApplication) stop() {
	wfa.eventManager.Publish(NewEvent(EventTypePrint, nil))
}

type WordWithZPrinter struct {
	eventManager *EventManager
	matched      []string
}

func NewWordWithZPrinter(em *EventManager) *WordWithZPrinter {
	wzp := &WordWithZPrinter{}
	em.Subscribe(EventTypeValidWord, func(e Event) {
		wzp.countWord(e.data.(string))
	})
	em.Subscribe(EventTypePrint, func(e Event) {
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
	em := NewEventManager()
	_ = NewDataStorage(em)
	_ = NewStopwordFilter(em, stopwordFile)
	_ = NewWordFrequencyCounter(em)
	_ = NewWordFrequencyApplication(em)
	// Exercise 16.2
	_ = NewWordWithZPrinter(em)
	em.Publish(NewEvent(EventTypeStart, inputFile))
}
