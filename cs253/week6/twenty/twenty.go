package main

import (
	"fmt"
	"io"
	"os"
	"plugin"

	"gopkg.in/ini.v1"
	"main/week6/twenty/commons"
)

var inputFile string
var stopwordsFile string
var stdout io.Writer
var stderr io.Writer

type WordFrequencyController struct {
	dsm commons.DataStorageManager
	sm  commons.StopwordManager
	wfm commons.WordFrequencyManager
}

func NewWordFrequencyController(inputFile, stopwordFile string) *WordFrequencyController {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Fprintf(stderr, "Fail to read file: %v", err)
		os.Exit(1)
	}
	app := cfg.Section("options").Key("app").String()

	plug, err := plugin.Open(app)
	if err != nil {
		fmt.Fprintf(stderr, "Fail to open plugin: %v", err)
		os.Exit(1)
	}

	symDsmInit, err := plug.Lookup("NewDataStorageManager")
	if err != nil {
		fmt.Fprintf(stderr, "Fail to lookup symbol: %v", err)
		os.Exit(1)
	}
	dsmInit := symDsmInit.(commons.DataStorageManagerInitFunc)

	symSwmInit, err := plug.Lookup("NewStopwordManager")
	if err != nil {
		fmt.Fprintf(stderr, "Fail to lookup symbol: %v", err)
		os.Exit(1)
	}
	swmInit := symSwmInit.(commons.StopwordManagerInitFunc)

	symWfmInit, err := plug.Lookup("NewWordFrequencyManager")
	if err != nil {
		fmt.Fprintf(stderr, "Fail to lookup symbol: %v", err)
		os.Exit(1)
	}
	wfmInit := symWfmInit.(commons.WordFrequencyManagerInitFunc)

	return &WordFrequencyController{
		dsm: dsmInit(inputFile),
		sm:  swmInit(stopwordsFile),
		wfm: wfmInit(),
	}
}

func (wfc *WordFrequencyController) Run() {
	for _, w := range wfc.dsm.Words() {
		if !wfc.sm.IsStopword(w) {
			wfc.wfm.IncrementCount(w)
		}
	}

	for _, freq := range wfc.wfm.Sorted()[0:25] {
		fmt.Fprintf(stdout, "%s  -  %d\n", freq.Word, freq.Count)
	}
}

func main() {
	stdout = os.Stdout
	stderr = os.Stderr
	inputFile = os.Args[1]
	stopwordsFile = "../../stop_words.txt"
	run()
}

func run() {
	wfc := NewWordFrequencyController(inputFile, stopwordsFile)
	wfc.Run()
}
