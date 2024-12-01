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

type DataStorageManager interface {
  Words() []string
}

type StopwordManager interface {
  IsStopword(word string) bool
}

type WordFrequencyManager interface {
  IncrementCount(word string)
  Sorted() []Freq
}

type Freq struct {
  word  string
  count int
}

type RegexDataStorageManager struct {
  data      string
  words     []string
  processed bool
}

func NewRegexDataStorageManager(inputFile string) *RegexDataStorageManager {
  wordsBytes, err := os.ReadFile(inputFile)
  if err != nil {
    os.Exit(1)
  }
  return &RegexDataStorageManager{
    data: strings.ToLower(string(wordsBytes)),
  }
}

func (rdsm *RegexDataStorageManager) Words() []string {
  if !rdsm.processed {
    pattern := regexp.MustCompile(`[a-z]{2,}`)
    rdsm.words = pattern.FindAllString(rdsm.data, -1)
    rdsm.processed = true
  }
  return rdsm.words
}

type MapStopwordManager struct {
  stopwords map[string]struct{}
}

func NewMapStopwordManager(stopwordFile string) *MapStopwordManager {
  stopwordsBytes, err := os.ReadFile(stopwordFile)
  if err != nil {
    os.Exit(1)
  }

  stopwords := strings.Split(string(stopwordsBytes), ",")
  stopwordMap := map[string]struct{}{}
  for _, stopword := range stopwords {
    stopwordMap[stopword] = struct{}{}
  }

  return &MapStopwordManager{
    stopwords: stopwordMap,
  }
}

func (msm *MapStopwordManager) IsStopword(word string) bool {
  _, ok := msm.stopwords[word]
  return ok
}

type MapWordFrequencyManager struct {
  freqs map[string]int
}

func NewMapWordFrequencyManager() *MapWordFrequencyManager {
  return &MapWordFrequencyManager{
    freqs: map[string]int{},
  }
}

func (mwfm *MapWordFrequencyManager) IncrementCount(word string) {
  mwfm.freqs[word]++
}

func (mwfm *MapWordFrequencyManager) Sorted() []Freq {
  wordFreqs := []Freq{}
  for k, v := range mwfm.freqs {
    wordFreqs = append(wordFreqs, Freq{k, v})
  }

  sort.Slice(wordFreqs, func(i, j int) bool {
    return wordFreqs[i].count >= wordFreqs[j].count
  })

  return wordFreqs
}

type WordFrequencyController struct {
  dsm DataStorageManager
  sm  StopwordManager
  wfm WordFrequencyManager
}

func NewWordFrequencyController(inputFile, stopwordFile string) *WordFrequencyController {
  return &WordFrequencyController{
    dsm: NewRegexDataStorageManager(inputFile),
    sm:  NewMapStopwordManager(stopwordsFile),
    wfm: NewMapWordFrequencyManager(),
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
