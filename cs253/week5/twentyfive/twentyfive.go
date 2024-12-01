package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type PlainFunc = func()
type OutFunc = func() any
type InFunc = func(any)
type InOutFunc = func(any) any

type Freq struct {
	word  string
	count int
}

type TFQuarantine struct {
	funcs []InOutFunc
}

func (tfq *TFQuarantine) Bind(f InOutFunc) *TFQuarantine {
	tfq.funcs = append(tfq.funcs, f)
	return tfq
}

func (tfq *TFQuarantine) Execute() {
	guardCallable := func(v any) any {
		if reflect.TypeOf(v).Kind() == reflect.Func {
			return v.(OutFunc)()
		} else {
			return v
		}
	}

	var v any
	v = func() any {
		return nil
	}
	for _, f := range tfq.funcs {
		v = f(guardCallable(v))
	}

	fmt.Fprintf(stdout, "%v", v)
}

func NewTFTFQuarantine() *TFQuarantine {
	return &TFQuarantine{}
}

var stdout io.Writer
var inputFile string
var stopwordsFile string = "../../stop_words.txt"

func main() {
	stdout = os.Stdout
	inputFile = os.Args[1]
	run()
}

func run() {
	ft := NewTFTFQuarantine()
	ft.
		Bind(getInputArgs).
		Bind(extractWords).
		Bind(removeStopwords).
		Bind(countFrequencies).
		Bind(sortFrequencies).
		Bind(top25Frequencies).
		Execute()
}

func getInputArgs(v any) any {
	return func() any {
		return inputFile
	}
}

func extractWords(inputFile any) any {
	return func() any {
		b, err := os.ReadFile(inputFile.(string))
		if err != nil {
			os.Exit(1)
		}
		pattern := regexp.MustCompile(`[a-z]{2,}`)
		words := pattern.FindAllString(strings.ToLower(string(b)), -1)
		return words
	}
}

func removeStopwords(words any) any {
	return func() any {
		b, err := os.ReadFile(stopwordsFile)
		if err != nil {
			os.Exit(1)
		}
		stopwords := strings.Split(string(b), ",")
		stopwordMap := map[string]struct{}{}
		for _, stopword := range stopwords {
			if _, ok := stopwordMap[stopword]; !ok {
				stopwordMap[stopword] = struct{}{}
			}
		}

		validWords := []string{}
		for _, word := range words.([]string) {
			if _, ok := stopwordMap[word]; !ok && len(word) >= 2 {
				validWords = append(validWords, word)
			}
		}

		return validWords
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

func top25Frequencies(v any) any {
	var sb strings.Builder
	freqs := v.([]Freq)
	for i := 0; i < 25; i++ {
		freq := freqs[i]
		fmt.Fprintf(&sb, "%s  -  %d\n", freq.word, freq.count)
	}
	return sb.String()
}
