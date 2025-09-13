package readers

import (
	"firefly-home-assigment/configs"
	"os"
)

type FileReader struct {
	Reader
}

func (r FileReader) Read() error {
	var err error
	var file *os.File

	file, err = os.Open(r.Path)
	if err != nil {
		return err
	}

	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	numOfChunks := int(info.Size()) / configs.EnvInt("CHUNK_SIZE", "4096")
	r.Wg.Add(numOfChunks)
	go func() {
		defer r.Wg.Done()
		err = r.ChunkProcessor(file)
	}()
	r.Wg.Wait()
	return err
}
