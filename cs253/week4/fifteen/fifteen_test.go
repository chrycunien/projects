package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testInputFile    = "../../pride-and-prejudice.txt"
	testStopwordFile = "../../stop_words.txt"
	testResultFile   = "../../pride-and-prejudice-result.txt"
)

func TestFifteen(t *testing.T) {
	var expected []byte
	var buf bytes.Buffer
	var buf2 bytes.Buffer
	var err error

	expected, err = os.ReadFile(testResultFile)
	require.Nil(t, err)

	stdout = &buf
	stderr = &buf2
	inputFile = testInputFile
	stopwordFile = testStopwordFile
	run()
	require.Equal(t, expected, buf.Bytes(), string(buf.Bytes()))
	require.Equal(t, "File contains: 837 words with z.\n", buf2.String(), buf2.String())
}
