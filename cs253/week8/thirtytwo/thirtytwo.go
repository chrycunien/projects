package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
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

func partitions(inputFile string, nlines int) <-chan string {
	ch := make(chan string)
	f, err := os.Open(inputFile)
	if err != nil {
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	go func() {
		defer f.Close()
		lines := []string{}
		for scanner.Scan() {
			lines = append(lines, strings.ToLower(scanner.Text()))
			if len(lines) == nlines {
				ch <- strings.Join(lines, "\n")
				lines = []string{}
			}
		}
		ch <- strings.Join(lines, "\n")
		close(ch)
	}()

	return ch
}

func splitWords(
	stopwordMap map[string]struct{},
	partitions <-chan string,
) []<-chan string {
	groups := []chan string{}
	for i := 0; i < 5; i++ {
		groups = append(groups, make(chan string))
	}

	go func() {
		pattern := regexp.MustCompile(`[a-z]{2,}`)
		for partition := range partitions {
			words := pattern.FindAllString(partition, -1)
			for _, word := range words {
				if _, ok := stopwordMap[word]; !ok {
					c := word[0]
          // Exercise 32.3
					switch {
					case c <= 'e':
						groups[0] <- word
					case c <= 'j':
						groups[1] <- word
					case c <= 'o':
						groups[2] <- word
					case c <= 't':
						groups[3] <- word
					default:
						groups[4] <- word
					}
				}
			}
		}

		for _, group := range groups {
			close(group)
		}
	}()

	outGroups := []<-chan string{}
	for _, group := range groups {
		outGroups = append(outGroups, group)
	}

	return outGroups
}

func groupFreq(groups []<-chan string) <-chan map[string]int {
	out := make(chan map[string]int)
	for i := 0; i < 5; i++ {
		go func() {
			wordMap := map[string]int{}
			for word := range groups[i] {
				wordMap[word] += 1
			}
			out <- wordMap
		}()
	}
	return out
}

func wordFreqs(wordMaps <-chan map[string]int) []Freq {
	freqMap := map[string]int{}
	for i := 0; i < 5; i++ {
		wordMap := <-wordMaps
		for word, count := range wordMap {
			freqMap[word] += count
		}
	}

	freqs := []Freq{}
	for k, v := range freqMap {
		freqs = append(freqs, Freq{
			word:  k,
			count: v,
		})
	}
	sort.Slice(freqs, func(i, j int) bool {
		return freqs[i].count >= freqs[j].count
	})

	return freqs
}

func main() {
	stdout = os.Stdout
	stderr = os.Stderr
	inputFile = os.Args[1]
	run()
}

func run() {
	b, err := os.ReadFile(stopwordsFile)
	if err != nil {
		os.Exit(1)
	}

	stopwords := strings.Split(string(b), ",")
	stopwordsMap := map[string]struct{}{}
	for _, stopWord := range stopwords {
		stopwordsMap[stopWord] = struct{}{}
	}

	splits := splitWords(maps.Clone(stopwordsMap), partitions(inputFile, 200))
	freqs := wordFreqs(groupFreq(splits))
	for _, freq := range freqs[0:25] {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}
