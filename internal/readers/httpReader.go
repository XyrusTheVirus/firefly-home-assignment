package readers

import (
	"net/http"
)

type HTTPReader struct {
	Reader
}

func (r HTTPReader) Read() error {
	var err error
	var resp *http.Response
	//var n int64

	resp, err = http.Get(r.Path)
	if err != nil {
		return err
	}

	//numOfChunks := int(resp.ContentLength) / configs.EnvInt("CHUNK_SIZE", "4096")

	r.Wg.Add(1)
	go func() {
		defer r.Wg.Done()
		err = r.ChunkProcessor(resp.Body)
	}()

	r.Wg.Wait()
	return err
}
