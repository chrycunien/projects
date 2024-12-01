package main

import (
	"flag"
	"fmt"
	"main/wordcount"
	"os"
)

var (
	inputFile     string
	stopwordsFile string
	commonN       int
	resultSep     string
)

func main() {
	parseArguments()

	var err error
	var wc wordcount.WordCounter = wordcount.NewImperativeCounter()
	if err = wc.Init(inputFile, stopwordsFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = wc.Process(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var freqs []wordcount.Frequency
	if freqs, err = wc.CountCommon(commonN); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wordcount.PrintFrequencies(os.Stdout, freqs, "  -  ")
}

func parseArguments() {
	flag.StringVar(&stopwordsFile, "stopwords", "../stop_words.txt", "Stop words file name")
	flag.StringVar(&resultSep, "sep", "  -  ", "Seperation string of word and frequencies in the printing result")
	flag.IntVar(&commonN, "n", 25, "Most common words count")
	flag.Parse()

	if commonN <= 0 {
		fmt.Println("n must be a positive number")
		os.Exit(1)
	}

	if stopwordsFile == "" {
		fmt.Println("stop words file cannot be empty")
		os.Exit(1)
	}

	inputFile = flag.Arg(0)
	if inputFile == "" {
		fmt.Println("input file cannot be empty")
		os.Exit(1)
	}
}
