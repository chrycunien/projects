package wordcount

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func PrintFrequencies(w io.Writer, freqs []Frequency, sep string) {
	for _, freq := range freqs {
		fmt.Fprintf(w, "%s%s%d\n", freq.Word, sep, freq.Count)
	}
}

func readInput(inputFile, stopwordsFile string) ([]byte, []byte, error) {
	var err error
	var inputBytes, stopwordsBytes []byte

	stopwordsBytes, err = os.ReadFile(stopwordsFile)
	if err != nil {
		return nil, nil, fmt.Errorf("Error reading stop words file: %s. Err: %w.\n", stopwordsFile, err)
	}

	inputBytes, err = os.ReadFile(inputFile)
	if err != nil {
    return nil, nil, fmt.Errorf("Error reading input file: %s. Err: %w.\n", inputFile, err)
	}

	return inputBytes, stopwordsBytes, nil
}

func parseStopwords(s string) []string {
	return strings.Split(s, ",")
}

func parseInput(s string) []string {
	lowerS := strings.ToLower(s)
	pattern := regexp.MustCompile(`[a-z]{2,}`)
	return pattern.FindAllString(lowerS, -1)
}
