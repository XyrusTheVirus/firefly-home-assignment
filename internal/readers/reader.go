package readers

import (
	"firefly-home-assigment/configs"
	"fmt"
	"io"
	"log"
	"strings"
)

type ReaderInterface interface {
	Read() error
}
type Reader struct {
	InputChannel chan []string
	QuitChannel  chan bool
	Path         string
}

// ChunkProcessor reads the resource in chunks and processes each chunk
func (r Reader) ChunkProcessor(resource io.ReadCloser) error {
	var err error
	var n int

	buf := make([]byte, configs.EnvInt("CHUNK_SIZE", "4096"))
	leftover := ""

	for {
		n, err = resource.Read(buf)
		if n > 0 {
			chunk := leftover + string(buf[:n])

			// Find last newline for natural boundary
			lastNewline := strings.LastIndex(chunk, "\n")

			var processChunk string
			//If no newline found, process whole chunk
			if lastNewline == len(chunk)-1 {
				processChunk = chunk[:lastNewline]
				// Else process up to last newline
			} else if lastNewline >= 0 {
				processChunk = chunk[:lastNewline]
				leftover = chunk[lastNewline+1:]
			}

			r.InputChannel <- strings.Split(processChunk, "\n")
		}

		if err == io.EOF {
			// Process any leftover data
			r.InputChannel <- strings.Split(leftover, "\n")
			// Signal completion
			r.QuitChannel <- true
			break

		} else if err != nil {
			log.Println(fmt.Sprintf("Encountering error while trying to read chunk: %s", err.Error()))
		}
	}

	return nil

}
