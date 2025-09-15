package extractors

import (
	"firefly-home-assigment/configs"
	"firefly-home-assigment/internal/readers"
	"os"
	"path"
	"regexp"
)

type WordsBank struct {
	Extractor
	re *regexp.Regexp
}

func NewWordsBank() WordsBank {
	return WordsBank{
		Extractor: Extractor{
			Result: make(map[string]interface{}),
		},
		re: regexp.MustCompile(`^[a-zA-Z]{3,}$`), // words with 3 or more letters
	}
}

// Extract reads a file containing words and stores valid words in the Result map.
func (w WordsBank) Extract() {
	wd, _ := os.Getwd()
	r := readers.FileReader{
		Reader: readers.Reader{
			InputChannel: make(chan []string),
			QuitChannel:  make(chan bool),
			Path:         path.Join(wd, configs.Env("BANK_OF_WORDS_FILE", "")),
		},
	}

	go func() {
		err := r.Read()
		if err != nil {

			panic(err)
		}
	}()

	for {
		select {
		// Read words from the input channel and store valid words in the Result map
		case input := <-r.InputChannel:
			for _, word := range input {
				if w.re.Match([]byte(word)) {
					w.Result[word] = struct{}{}
				}
			}
		case <-r.QuitChannel:
			return
		}
	}
}
