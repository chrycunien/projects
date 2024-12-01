package wordcount

import (
	"errors"
	"sort"
)

var _ WordCounter = (*ImperativeWordCounter)(nil)

type ImperativeWordCounter struct {
	inputStr     string
	stopwordsStr string
	words        []string
	stopwords    []string
	wordsMap     map[string]int
	stopwordsMap map[string]struct{}
}

func NewImperativeCounter() *ImperativeWordCounter {
  return &ImperativeWordCounter{
    wordsMap: map[string]int{},
    stopwordsMap: map[string]struct{}{},
  }
}

func (iwc *ImperativeWordCounter) Init(inputFile, stopwordsFile string) error {
	inputBytes, stopwordsBytes, err := readInput(inputFile, stopwordsFile)
	if err != nil {
		return err
	}

	iwc.inputStr = string(inputBytes)
	iwc.stopwordsStr = string(stopwordsBytes)
	return nil
}

func (iwc *ImperativeWordCounter) Process() error {
	iwc.stopwords = parseStopwords(iwc.stopwordsStr)
  iwc.words = parseInput(iwc.inputStr)

	for _, stopword := range iwc.stopwords {
		if _, ok := iwc.stopwordsMap[stopword]; !ok {
			iwc.stopwordsMap[stopword] = struct{}{}
		}
	}

	for _, word := range iwc.words {
		if _, ok := iwc.stopwordsMap[word]; !ok {
			if _, ok := iwc.wordsMap[word]; !ok {
				iwc.wordsMap[word] = 1
			} else {
				iwc.wordsMap[word]++
			}
		}
	}

  return nil
}

func (iwc *ImperativeWordCounter) CountCommon(n int) ([]Frequency, error) {
  if len(iwc.wordsMap) < n {
    return nil, errors.New("n too big: not enough entries")
  }
  
  freqs := make([]Frequency, 0, len(iwc.wordsMap))
  for word, count := range iwc.wordsMap {
    freqs = append(freqs, Frequency{
      Word:   word,
      Count: count,
    })
  }

  sort.Slice(freqs, func(i, j int) bool {
    return freqs[i].Count > freqs[j].Count
  })

  commons := make([]Frequency, 0, n)
  for i := 0; i < n; i++ {
    commons = append(commons, freqs[i])
  }

  return commons, nil
}
