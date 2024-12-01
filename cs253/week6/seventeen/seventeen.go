package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

var inputFile string
var stopwordsFile string
var stdout io.Writer
var stderr io.Writer

type Freq struct {
	word  string
	count int
}

type DataStorageManager struct {
	data      string
	words     []string
	processed bool
}

func NewDataStorageManager(inputFile string) *DataStorageManager {
	wordsBytes, err := os.ReadFile(inputFile)
	if err != nil {
		os.Exit(1)
	}
	return &DataStorageManager{
		data: strings.ToLower(string(wordsBytes)),
	}
}

func (dsm *DataStorageManager) Words() []string {
	if !dsm.processed {
		pattern := regexp.MustCompile(`[a-z]{2,}`)
		dsm.words = pattern.FindAllString(dsm.data, -1)
		dsm.processed = true
	}
	return dsm.words
}

type StopwordManager struct {
	stopwords map[string]struct{}
}

func NewStopwordManager(stopwordFile string) *StopwordManager {
	stopwordsBytes, err := os.ReadFile(stopwordFile)
	if err != nil {
		os.Exit(1)
	}

	stopwords := strings.Split(string(stopwordsBytes), ",")
	stopwordMap := map[string]struct{}{}
	for _, stopword := range stopwords {
		stopwordMap[stopword] = struct{}{}
	}

	return &StopwordManager{
		stopwords: stopwordMap,
	}
}

func (sm *StopwordManager) IsStopword(word string) bool {
	_, ok := sm.stopwords[word]
	return ok
}

type WordFrequencyManager struct {
	freqs map[string]int
}

func NewWordFrequencyManager() *WordFrequencyManager {
	return &WordFrequencyManager{
		freqs: map[string]int{},
	}
}

func (wfm *WordFrequencyManager) IncrementCount(word string) {
	wfm.freqs[word]++
}

func (wfm *WordFrequencyManager) Sorted() []Freq {
	wordFreqs := []Freq{}
	for k, v := range wfm.freqs {
		wordFreqs = append(wordFreqs, Freq{k, v})
	}

	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].count >= wordFreqs[j].count
	})

	return wordFreqs
}

type WordFrequencyController struct {
	dsm reflect.Value
	sm  reflect.Value
	wfm reflect.Value
}

func NewWordFrequencyController(inputFile, stopwordFile string) *WordFrequencyController {
	return &WordFrequencyController{
		dsm: reflect.ValueOf(NewDataStorageManager(inputFile)),
		sm:  reflect.ValueOf(NewStopwordManager(stopwordsFile)),
		wfm: reflect.ValueOf(NewWordFrequencyManager()),
	}
}

func (wfc *WordFrequencyController) Run() {
	wordsRet := wfc.dsm.MethodByName("Words").Call([]reflect.Value{})
	for _, w := range wordsRet[0].Interface().([]string) {
		wVal := reflect.ValueOf(w)
		values := wfc.sm.MethodByName("IsStopword").Call([]reflect.Value{wVal})
		if !values[0].Bool() {
			wfc.wfm.MethodByName("IncrementCount").Call([]reflect.Value{wVal})
		}
	}

	sortedRet := wfc.wfm.MethodByName("Sorted").Call([]reflect.Value{})
	for _, freq := range sortedRet[0].Interface().([]Freq)[0:25] {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.word, freq.count)
	}
}

type TypeRegistry struct {
	p map[string]reflect.Type
	v map[string]reflect.Type
	w io.Writer
}

func NewTypeRegistry(w io.Writer) *TypeRegistry {
	return &TypeRegistry{
		p: make(map[string]reflect.Type),
		v: make(map[string]reflect.Type),
		w: w,
	}
}

func (reg *TypeRegistry) register(typedNil any) {
	t := reflect.TypeOf(typedNil)
	e := t.Elem()
	reg.p[e.PkgPath()+"."+e.Name()] = t
	reg.v[e.PkgPath()+"."+e.Name()] = e
}

func (reg *TypeRegistry) make(name string) any {
	return reflect.New(reg.p[name]).Elem().Interface()
}

func (reg *TypeRegistry) print(name string) {
	// print fields
	vt, ok := reg.v[name]
	if !ok {
		panic("struct not exists")
	}
	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		fmt.Fprintf(reg.w, "field: %s - %v\n", f.Name, f.Type)
	}
	// print methods
	pt, ok := reg.p[name]
	if !ok {
		panic("struct not exists")
	}
	for i := 0; i < pt.NumMethod(); i++ {
		m := pt.Method(i)
		fmt.Fprintf(reg.w, "method: %s - %v\n", m.Name, m.Type)
	}

	// superclasses: go don't have inheritance
	// interfaces: not available in reflect package
	//             go implements duck-typing so it only implicitly
	//             knows if has an interface it you call it
}

func main() {
	stdout = os.Stdout
	stderr = os.Stderr
	inputFile = os.Args[1]
	stopwordsFile = "../../stop_words.txt"
	run()
	printStruct()
}

func run() {
	wfc := NewWordFrequencyController(inputFile, stopwordsFile)
	wfc.Run()
}

func printStruct() {
	typeRegistry := NewTypeRegistry(stderr)
	typeRegistry.register((*DataStorageManager)(nil))
	typeRegistry.register((*StopwordManager)(nil))
	typeRegistry.register((*WordFrequencyManager)(nil))

	var name string
	fmt.Fprintln(stderr, "Enter the struct name: ")
	fmt.Scanln(&name)
	typeRegistry.print(name)
}
