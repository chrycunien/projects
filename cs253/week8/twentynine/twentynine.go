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

type Message struct {
	msg  string
	data []any
}

func NewMessage(msg string, data ...any) Message {
	return Message{
		msg:  msg,
		data: data,
	}
}

var stdout io.Writer
var stderr io.Writer
var inputFile string
var stopwordsFile string = "../../stop_words.txt"

type Actor interface {
	Send(msg Message)
	Run()
}

var _ Actor = (*DataStorageManager)(nil)
var _ Actor = (*StopwordManager)(nil)
var _ Actor = (*WordFreqManager)(nil)
var _ Actor = (*WordFreqController)(nil)

type DataStorageManager struct {
	inputFile       string
	data            string
	words           []string
	stopWordManager Actor
	queue           chan Message
	done            <-chan struct{}
}

func NewDataStorageManager(done <-chan struct{}) *DataStorageManager {
	return &DataStorageManager{
		queue: make(chan Message),
		done:  done,
	}
}

func (dsm *DataStorageManager) Send(msg Message) {
	dsm.queue <- msg
}

func (dsm *DataStorageManager) Run() {
	go func() {
		for {
			select {
			case msg := <-dsm.queue:
				switch msg.msg {
				case "init":
					dsm.inputFile = msg.data[0].(string)
					dsm.stopWordManager = msg.data[1].(Actor)
					dsm.readInput()
				case "send_word_freqs":
					recipient := msg.data[0].(Actor)
					dsm.process()
					dsm.forward(recipient)
				default:
					dsm.stopWordManager.Send(msg)
				}
			case <-dsm.done:
				return
			}
		}
	}()
}

func (dsm *DataStorageManager) readInput() {
	wordsBytes, err := os.ReadFile(dsm.inputFile)
	if err != nil {
		os.Exit(1)

	}
	dsm.data = strings.ToLower(string(wordsBytes))
}

func (dsm *DataStorageManager) process() {
	pattern := regexp.MustCompile(`[a-z]{2,}`)
	dsm.words = pattern.FindAllString(dsm.data, -1)
}

func (dsm *DataStorageManager) forward(recipient Actor) {
	for _, word := range dsm.words {
		dsm.stopWordManager.Send(NewMessage("filter", word))
	}
	dsm.stopWordManager.Send(NewMessage("top25", recipient))
}

type StopwordManager struct {
	stopwordFile    string
	stopwords       map[string]struct{}
	wordFreqManager Actor
	queue           chan Message
	done            <-chan struct{}
}

func NewStopwordManager(done <-chan struct{}) *StopwordManager {
	return &StopwordManager{
		queue:     make(chan Message),
		done:      done,
		stopwords: map[string]struct{}{},
	}
}

func (sm *StopwordManager) Send(msg Message) {
	sm.queue <- msg
}

func (sm *StopwordManager) Run() {
	go func() {
		for {
			select {
			case msg := <-sm.queue:
				switch msg.msg {
				case "init":
					sm.stopwordFile = msg.data[0].(string)
					sm.wordFreqManager = msg.data[1].(Actor)
					sm.readStopwords()
				case "filter":
					word := msg.data[0].(string)
					if sm.isStopword(word) {
						sm.wordFreqManager.Send(NewMessage("word", word))
					}
				default:
					sm.wordFreqManager.Send(msg)
				}
			case <-sm.done:
				return
			}
		}
	}()
}

func (sm *StopwordManager) readStopwords() {
	b, err := os.ReadFile(sm.stopwordFile)
	if err != nil {
		os.Exit(1)
	}

	stopWords := strings.Split(string(b), ",")
	stopWordsMap := map[string]struct{}{}
	for _, stopWord := range stopWords {
		stopWordsMap[stopWord] = struct{}{}
	}

	sm.stopwords = stopWordsMap
}

func (sm *StopwordManager) isStopword(word string) bool {
	_, ok := sm.stopwords[word]
	return !ok
}

type WordFreqManager struct {
	wordFreq map[string]int
	queue    chan Message
	done     <-chan struct{}
}

func NewWordFreqManager(done <-chan struct{}) *WordFreqManager {
	return &WordFreqManager{
		queue:    make(chan Message),
		done:     done,
		wordFreq: map[string]int{},
	}
}

func (wfm *WordFreqManager) Send(msg Message) {
	wfm.queue <- msg
}

func (wfm *WordFreqManager) Run() {
	go func() {
		for {
			select {
			case msg := <-wfm.queue:
				switch msg.msg {
				case "word":
					word := msg.data[0].(string)
					wfm.count(word)
				case "top25":
					recipient := msg.data[0].(Actor)
					wfm.top25(recipient)
				}
			case <-wfm.done:
				return
			}
		}
	}()
}

func (wfm *WordFreqManager) count(word string) {
	wfm.wordFreq[word] += 1
}

func (wfm *WordFreqManager) top25(recipient Actor) {
	freqs := []Freq{}
	for k, v := range wfm.wordFreq {
		freqs = append(freqs, Freq{
			word:  k,
			count: v,
		})
	}
	sort.Slice(freqs, func(i, j int) bool {
		return freqs[i].count >= freqs[j].count
	})

	recipient.Send(NewMessage("top25", freqs[0:25]))
}

type WordFreqController struct {
	queue chan Message
	stop  chan<- struct{}
	done  chan<- struct{}
}

func NewWordFreqController(done, stop chan<- struct{}) *WordFreqController {
	return &WordFreqController{
		queue: make(chan Message),
		stop:  stop,
		done:  done,
	}
}

func (wfc *WordFreqController) Send(msg Message) {
	wfc.queue <- msg
}

func (wfc *WordFreqController) Run() {
	go func() {
		for {
			for msg := range wfc.queue {
				switch msg.msg {
				case "run":
					storageManager := msg.data[0].(Actor)
					storageManager.Send(NewMessage("send_word_freqs", wfc))
				case "top25":
					freqs := msg.data[0].([]Freq)
					wfc.display(freqs)
					// TODO: maybe count how many children to close
					close(wfc.done)
					close(wfc.stop)
				}
			}
		}
	}()
}

func (wfc *WordFreqController) display(freqs []Freq) {
	for _, freq := range freqs {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}

func main() {
	stdout = os.Stdout
	stderr = os.Stderr
	inputFile = os.Args[1]
	run()
}

func run() {
	stop := make(chan struct{})
	done := make(chan struct{})

	dataStorageManager := NewDataStorageManager(done)
	stopwordManager := NewStopwordManager(done)
	wordFreqManager := NewWordFreqManager(done)
	wordFreqController := NewWordFreqController(done, stop)

	dataStorageManager.Run()
	stopwordManager.Run()
	wordFreqManager.Run()
	wordFreqController.Run()

	dataStorageManager.Send(NewMessage("init", inputFile, stopwordManager))
	stopwordManager.Send(NewMessage("init", stopwordsFile, wordFreqManager))
	wordFreqController.Send(NewMessage("run", dataStorageManager))
	<-stop
}
