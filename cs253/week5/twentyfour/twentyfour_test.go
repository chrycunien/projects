package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testInputFile  = "../../pride-and-prejudice.txt"
	testResultFile = "../../pride-and-prejudice-result.txt"
)

func TestTwentyFour(t *testing.T) {
	var expected []byte
	var buf bytes.Buffer
	var err error

	expected, err = os.ReadFile(testResultFile)
	require.Nil(t, err)

	stdout = &buf
	inputFile = testInputFile
	run()
	require.Equal(t, expected, buf.Bytes(), string(buf.Bytes()))
}

func BenchmarkTwentyFour(b *testing.B) {
	var buf bytes.Buffer
	stdout = &buf
	inputFile = testInputFile
	for i := 0; i < b.N; i++ {
		run()
	}
}
