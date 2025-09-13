package readers

import (
	"firefly-home-assigment/configs"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
)

type ReaderInterface interface {
	Read() error
}
type Reader struct {
	Channel chan string
	Result  map[string]string
	Wg      *sync.WaitGroup
	Path    string
}

// ChunkProcessor reads the resource in chunks and processes each chunk

func (r Reader) ChunkProcessor(resource io.ReadCloser) error {
	var err error
	var n int

	buf := make([]byte, configs.EnvInt("CHUNK_SIZE", "4096"))
	leftover := ""
	counter := 0
	for {
		n, err = resource.Read(buf)
		counter += 1
		fmt.Println(counter)
		if n > 0 {
			chunk := leftover + string(buf[:n])

			// Find last newline for natural boundary
			lastNewline := strings.LastIndex(chunk, "\n")

			var processChunk string
			if lastNewline == len(chunk)-1 {
				processChunk = chunk[:lastNewline]
			} else {
				processChunk = chunk[:lastNewline]
				leftover = chunk[lastNewline+1:]
			}
			fmt.Println(processChunk)
			r.processChunk(strings.Split(processChunk, "\n"))
		}

		if err == io.EOF {
			r.processChunk(strings.Split(leftover, "\n"))
			break

		} else if err != nil {
			return err
		}
	}

	return nil

}

func (r Reader) processChunk(chunk []string) {
	re, _ := regexp.Compile(`^[a-zA-Z]+$`)
	for _, word := range chunk {
		if len(word) >= 3 && re.Match([]byte(word)) {
			r.Result[word] = word
		}
	}
}
