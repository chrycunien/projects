package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	testInputFile  = "../../pride-and-prejudice.txt"
	testResultFile = "../../pride-and-prejudice-result.txt"
)

func TestThirty(t *testing.T) {
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

func BenchmarkThirty(b *testing.B) {
	var buf bytes.Buffer
	stdout = &buf
	inputFile = testInputFile
	for i := 0; i < b.N; i++ {
		run()
	}
}
