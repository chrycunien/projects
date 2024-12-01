package commons

type Freq struct {
	Word  string
	Count int
}

type DataStorageManager interface {
	Words() []string
}

type StopwordManager interface {
	IsStopword(word string) bool
}

type WordFrequencyManager interface {
	IncrementCount(word string)
	Sorted() []Freq
}

type DataStorageManagerInitFunc = func(inputFile string) DataStorageManager
type StopwordManagerInitFunc = func(stopwordFile string) StopwordManager
type WordFrequencyManagerInitFunc = func() WordFrequencyManager
