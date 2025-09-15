package main

import (
	"firefly-home-assigment/internal/extractors"
	"firefly-home-assigment/internal/outputs"
	"firefly-home-assigment/internal/sorts"
	"log"
	"runtime/debug"
)

func init() {
	// Global panic recovery to log stack trace
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("Panic occurred: ", err, ", Trace: ", string(debug.Stack()))
		}
	}()
}

func main() {
	// Extract words bank
	wordsBank := extractors.NewWordsBank()
	wordsBank.Extract()
	// Extract essays and count words
	essays := extractors.NewEssay(wordsBank.Result)
	essays.Extract()
	// Sort and output results
	output := &outputs.JSONOutput{Output: outputs.Output{Data: sorts.SortByFrequency(essays.Result)}}
	output.Print()
}
