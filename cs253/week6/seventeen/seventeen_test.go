package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testInputFile     = "../../pride-and-prejudice.txt"
	testStopwordsFile = "../../stop_words.txt"
	testResultFile    = "../../pride-and-prejudice-result.txt"
)

func TestSeventeen(t *testing.T) {
	var expected []byte
	var buf bytes.Buffer
	var buf2 bytes.Buffer
	var err error

	expected, err = os.ReadFile(testResultFile)
	require.Nil(t, err)

	stdout = &buf
	stderr = &buf2
	inputFile = testInputFile
	stopwordsFile = testStopwordsFile
	run()
	// printStruct()
	require.Equal(t, expected, buf.Bytes(), string(buf.Bytes()))
}

func BenchmarkSeventeen(b *testing.B) {
	var buf bytes.Buffer
	stdout = &buf
	inputFile = testInputFile
	for i := 0; i < b.N; i++ {
		wfc := NewWordFrequencyController(inputFile, stopwordsFile)
		wfc.Run()
	}
}
