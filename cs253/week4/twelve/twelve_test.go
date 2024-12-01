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

func TestTwelve(t *testing.T) {
	var expected []byte
	var buf bytes.Buffer
	var err error

	expected, err = os.ReadFile(testResultFile)
	require.Nil(t, err)

	stdout = &buf
	inputFile = testInputFile
	stopwordsFile = testStopwordsFile
	wfc := &WordFrequencyController{}

	_, err = wfc.Dispatch(NewMessage("init", []string{inputFile, stopwordsFile}))
	require.Nil(t, err)

	_, err = wfc.Dispatch(NewMessage("run", nil))
	require.Nil(t, err)

	require.Equal(t, expected, buf.Bytes(), string(buf.Bytes()))
}
