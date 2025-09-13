package main

import (
	"firefly-home-assigment/configs"
	"firefly-home-assigment/internal/readers"
	"log"
	"sync"
)

func main() {
	var err error
	r := readers.HTTPReader{
		readers.Reader{
			Channel: make(chan string, 3),
			Result:  make(map[string]string),
			Wg:      &sync.WaitGroup{},
			Path:    configs.Env("BANK_OF_WORDS_URL", ""),
		},
	}

	err = r.Read()
	//x := readers.FileReader{
	//	readers.Reader{
	//		Channel: make(chan string, 3),
	//		Result:  make(map[string]string),
	//		Wg:      &sync.WaitGroup{},
	//		Path:    "/home/galb/GolandProjects/firefly-home-assignment/endg-urls",
	//	},
	//}
	//err = x.Read()
	if err != nil {
		log.Fatal(err)
	}

}
