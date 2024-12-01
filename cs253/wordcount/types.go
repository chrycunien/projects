package wordcount

type Frequency struct {
  Word string
  Count int
}

type WordCounter interface {
  Init(inputFile, stopwordsFile string) error
  Process() error
  CountCommon(n int) ([]Frequency, error)
}