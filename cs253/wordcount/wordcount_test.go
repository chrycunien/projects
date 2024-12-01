package wordcount

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
  testStopwordsFile = "../stop_words.txt"
  testInputFile = "../pride-and-prejudice.txt"
  testResultFile = "../pride-and-prejudice-result.txt"
  testResultSep = "  -  "
  testCommonN = 25
)

func readTestResult() ([]byte, error) {
  return os.ReadFile(testResultFile)
}

func testWordCounterInterface(t *testing.T, wc WordCounter) {
  var expected []byte
  var freqs []Frequency
  var buf bytes.Buffer
  var err error

  expected, err = readTestResult()
  require.Nil(t, err)

  err = wc.Init(testInputFile, testStopwordsFile)
  require.Nil(t, err)

  err = wc.Process()
  require.Nil(t, err)

  freqs, err = wc.CountCommon(testCommonN)
  require.Nil(t, err)

  PrintFrequencies(&buf, freqs, testResultSep)
  require.Equal(t, expected, buf.Bytes())
}

func TestWordCounters(t *testing.T) {
  t.Parallel()
  
  testcases := []struct {
    name string
    wc WordCounter
  }{
    {
      name: "ImperativeCounter",
      wc: NewImperativeCounter(),
    },
  }
  
  for _, tc := range testcases {
    tc := tc
    t.Run(tc.name, func(t *testing.T) {
      t.Parallel()
      testWordCounterInterface(t, tc.wc)
    })
  }
}