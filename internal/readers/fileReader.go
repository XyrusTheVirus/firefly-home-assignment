package readers

import (
	"fmt"
	"os"
)

type FileReader struct {
	Reader
}

// Read opens the file at the specified path and processes it in chunks
func (r FileReader) Read() error {
	var err error
	var file *os.File

	file, err = os.Open(r.Path)
	if err != nil {
		panic(fmt.Sprintf("Couldn't open file: %s", err.Error())) // Panic here as this is a critical error
	}

	defer file.Close()
	err = r.ChunkProcessor(file)
	return err
}
