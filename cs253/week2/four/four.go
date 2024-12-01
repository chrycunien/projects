package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Freq struct {
	word  string
	count int
}

const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var stdout io.Writer = os.Stdout
var inputFile string = os.Args[1]
var stopwordsFile string = "../../stop_words.txt"

func main() {
	run()
}

func run() {
	var wordFreqs []Freq
	var stopwords []string
	var f *os.File
	var b []byte
	var err error

	b, err = os.ReadFile(stopwordsFile)
	if err != nil {
		os.Exit(1)
	}
	stopwords = strings.Split(string(b), ",")

	f, err = os.Open(inputFile)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()

	// iterate through the file one line at a time
	word := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word = ""
		line := scanner.Text() + "\n"
		for i := 0; i < len(line); i++ {
			c := line[i]

			isAlnum := false
			for k := 0; k < len(alnum); k++ {
				if c == alnum[k] {
					isAlnum = true
					break
				}
			}

			if isAlnum {
				// lowercase the char and append it to the word
				if c >= 'A' && c <= 'Z' {
					c = c - 'A' + 'a'
				}
				word += string(c)
				continue
			}

			// not reaching the next word
			if word == "" {
				continue
			}
			// ignore single character
			if len(word) < 2 {
				word = ""
				continue
			}

			// ingore stopwords
			inStopwords := false
			for _, stopword := range stopwords {
				if word == stopword {
					inStopwords = true
					break
				}
			}

			if inStopwords {
				word = ""
				continue
			}

			// fmt.Println(word)

			inFreq := false
			freqPos := 0
			for pos, freq := range wordFreqs {
				if freq.word == word {
					inFreq = true
					freqPos = pos
					wordFreqs[freqPos].count += 1
					break
				}
			}

			// not in freq array, just append
			if !inFreq {
				wordFreqs = append(wordFreqs, Freq{
					word:  word,
					count: 1,
				})
				word = ""
				continue
			}

			// sort the freq array
			for j := freqPos - 1; j >= 0; j-- {
				if wordFreqs[j].count >= wordFreqs[freqPos].count {
					break
				}
				// swap
				wordFreqs[j], wordFreqs[freqPos] = wordFreqs[freqPos], wordFreqs[j]
				// fmt.Println(wordFreqs[j], wordFreqs[freqPos])
				freqPos = j
			}
			word = ""
		}
	}

	for i := 0; i < 25; i++ {
		freq := wordFreqs[i]
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}
