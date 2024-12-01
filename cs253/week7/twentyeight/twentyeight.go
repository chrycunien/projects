package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Freq struct {
	word  string
	count int
}

var stdout io.Writer
var stderr io.Writer
var inputFile string
var stopwordsFile string = "../../stop_words.txt"

// Exercise 28.2
func lines(inputFile string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		f, err := os.Open(inputFile)
		if err != nil {
			os.Exit(1)
		}

		scanner := bufio.NewScanner(f)
		defer f.Close()
		count := 0
		for scanner.Scan() {
			line := scanner.Text()
			if !yield(count, line) {
				return
			}
			count++
		}
	}
}

// Exercise 28.2
func allWords(inputFile string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		count := 0
		pattern := regexp.MustCompile(`[a-z]{2,}`)
		for _, line := range lines(inputFile) {
			words := pattern.FindAllString(strings.ToLower(line), -1)
			for _, word := range words {
				if !yield(count, word) {
					return
				}
				count++
			}
		}
	}
}

func nonStopWords(inputFile string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		b, err := os.ReadFile(stopwordsFile)
		if err != nil {
			os.Exit(1)
		}

		stopWords := strings.Split(string(b), ",")
		stopWordsMap := map[string]struct{}{}
		for _, stopWord := range stopWords {
			stopWordsMap[stopWord] = struct{}{}
		}

		count := 0
		for _, word := range allWords(inputFile) {
			if _, ok := stopWordsMap[word]; !ok {
				if !yield(count, word) {
					return
				}
				count++
			}
		}
	}
}

func countAndSort(inputFile string) iter.Seq2[int, []Freq] {
	return func(yield func(int, []Freq) bool) {
		count, batch := 1, 0
		wordMap := map[string]int{}
		for _, word := range nonStopWords(inputFile) {
			wordMap[word] += 1
			if count%5000 == 0 {
				freqs := []Freq{}
				for k, v := range wordMap {
					freqs = append(freqs, Freq{
						word:  k,
						count: v,
					})
				}
				sort.Slice(freqs, func(i, j int) bool {
					return freqs[i].count >= freqs[j].count
				})
				if !yield(batch, freqs) {
					return
				}
				batch++
			}
			count++
		}
	}
}

func main() {
	stdout = os.Stdout
	stderr = os.Stderr
	inputFile = os.Args[1]
	run()
}

func run() {
	for _, freqs := range countAndSort(inputFile) {
		fmt.Fprintln(stdout, "-----------------------------")
		for _, freq := range freqs[0:25] {
			fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
		}
	}
}
