package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	// "slices"
	// "sort"
	"strings"
)

type Freq struct {
	word  string
	count int
}

const RecursionLimit = 5000

var stdout io.Writer
var inputFile string

func main() {
	inputFile = os.Args[1]
	stdout = os.Stdout
	run()
}

func run() {
	stopwordsBytes, err := os.ReadFile("../../stop_words.txt")
	if err != nil {
		os.Exit(1)
	}
	stopwords := strings.Split(string(stopwordsBytes), ",")
	stopwordMap := map[string]struct{}{}
	for _, stopword := range stopwords {
		if _, ok := stopwordMap[stopword]; !ok {
			stopwordMap[stopword] = struct{}{}
		}
	}

	wordsBytes, err := os.ReadFile(inputFile)
	pattern := regexp.MustCompile(`[a-z]{2,}`)
	words := pattern.FindAllString(strings.ToLower(string(wordsBytes)), -1)

	wordFreqMap := map[string]int{}
	for i := 0; i < len(words); i += RecursionLimit {
		if i+RecursionLimit >= len(words) {
			countWordFreq(words[i:i+RecursionLimit], stopwordMap, wordFreqMap)
		} else {
			countWordFreq(words[i:i+RecursionLimit], stopwordMap, wordFreqMap)
		}
	}

	wordFreqs := []Freq{}
	for k, v := range wordFreqMap {
		wordFreqs = append(wordFreqs, Freq{k, v})
	}

	sortWordFreqs(wordFreqs)

	printWordFreq(stdout, wordFreqs[0:25])
}

func countWordFreq(words []string, stopwordMap map[string]struct{}, wordFreqMap map[string]int) {
	if len(words) == 0 {
		return
	}
	word := words[0]
	if _, ok := stopwordMap[word]; !ok && len(word) >= 2 {
		wordFreqMap[word]++
	}
	countWordFreq(words[1:], stopwordMap, wordFreqMap)
}

func sortWordFreqs(wordFreqs []Freq) {
	if len(wordFreqs) <= 1 {
		return
	}

	sortWordFreqs(wordFreqs[1:])

	// bubble sort
	for i := 0; i < len(wordFreqs)-1; i++ {
		if wordFreqs[i].count >= wordFreqs[i+1].count {
			break
		}
		wordFreqs[i], wordFreqs[i+1] = wordFreqs[i+1], wordFreqs[i]
	}
}

func printWordFreq(w io.Writer, wordFreqs []Freq) {
	if len(wordFreqs) == 0 {
		return
	}
	fmt.Fprintf(w, "%s  -  %d\n", wordFreqs[0].word, wordFreqs[0].count)
	printWordFreq(w, wordFreqs[1:])
}

// func runImperative() {
//   stopwordsBytes, err := os.ReadFile("../../stop_words.txt")
//   if err != nil {
//     os.Exit(1)
//   }
//   stopwords := strings.Split(string(stopwordsBytes), ",")
//   stopwordMap := map[string]struct{}{}
//   for _, stopword := range stopwords {
//     if _, ok := stopwordMap[stopword]; !ok {
//       stopwordMap[stopword] = struct{}{}
//     }
//   }

//   wordsBytes, err := os.ReadFile(inputFile)
//   pattern := regexp.MustCompile(`[a-z]{2,}`)
//   words := pattern.FindAllString(strings.ToLower(string(wordsBytes)), -1)

//   wordFreqMap := map[string]int{}
//   for _, word := range words {
//     if !slices.Contains(stopwords, word) {
//       wordFreqMap[word]++
//     }
//   }

//   wordFreqs := []Freq{}
//   for k, v := range wordFreqMap {
//     wordFreqs = append(wordFreqs, Freq{k, v})
//   }

//   sort.Slice(wordFreqs, func(i, j int) bool {
//     return wordFreqs[i].count >= wordFreqs[j].count
//   })

//   for _, freq := range wordFreqs[0:25] {
//     fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
//   }
// }
