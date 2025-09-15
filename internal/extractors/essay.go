package extractors

import (
	"firefly-home-assigment/configs"
	"firefly-home-assigment/internal/counters"
	"firefly-home-assigment/internal/readers"
	"firefly-home-assigment/internal/transporters"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Essay struct {
	Extractor
	WordsBank map[string]interface{}
}

// NewEssay creates a new Essay extractor with the provided words bank.
func NewEssay(wordsBank map[string]interface{}) Essay {
	return Essay{
		Extractor: Extractor{
			Result: make(map[string]interface{}),
		},
		WordsBank: wordsBank,
	}
}

// Extract reads essay URLs from a file, fetches each essay, extracts text from the <article> tag,
// tokenizes the text into words (≥3 letters), and counts occurrences of words present in
func (e Essay) Extract() {
	wd, _ := os.Getwd()
	r := readers.FileReader{
		Reader: readers.Reader{
			InputChannel: make(chan []string),
			QuitChannel:  make(chan bool),
			Path:         path.Join(wd, configs.Env("ESSAYS_FILE", "")),
		},
	}

	// Start reading the file in a separate goroutine
	go func() {
		err := r.Read()
		if err != nil {
			panic(err)
		}
	}()
	counter := &counters.WordsCounter{Wc: e.Result}
	// These 2 wait groups manage concurrency between processing chunks and fetching URLs.
	//While the chunkWorkGroup ensures all URLs in a chunk are processed before moving to the next chunk, the processWaitGroup ensures all chunks are fully processed before exiting.
	chunkWorkGroup := sync.WaitGroup{}
	processWaitGroup := sync.WaitGroup{}
	for {
		select {
		case urls := <-r.InputChannel:
			processWaitGroup.Add(1)
			go func() {
				defer processWaitGroup.Done()
				for _, url := range urls {
					chunkWorkGroup.Add(1)
					go func() {
						defer chunkWorkGroup.Done()

						// Fetch URL with rate limiting
						t := transporters.NewHttp(
							http.MethodGet,
							url,
							nil,
							nil,
						)
						resp, err := t.Transport()
						// Here we won't halt the entire process if one URL fails.
						// Instead, we log the error and continue with the next URL.
						// This is crucial for robustness when dealing with many URLs.
						if err != nil {
							log.Printf("Failed to fetch URL %s: %s", url, err.Error())
							return
						}
						// Extract text from <article>
						text, err := extractArticleText(resp.(*http.Response).Body)
						resp.(*http.Response).Body.Close()
						// Log and skip if extraction fails
						if err != nil {
							log.Printf("Couldn't extract words for url: %s, due to following error: %s\n", url, err.Error())
							return
						}

						// tokenize (only words ≥3 letters)
						wordRe := regexp.MustCompile(`[A-Za-z]{3,}`)
						words := wordRe.FindAllString(text, -1)
						// Count words present in the words bank
						counter.Count(words, e.WordsBank)
					}()
				}
			}()
			// Wait for all URLs in the current chunk to be processed before reading the next chunk
			chunkWorkGroup.Wait()

		case <-r.QuitChannel:
			// Wait for all processing to complete before exiting
			processWaitGroup.Wait()
			return
		}
	}
}

// extractArticleText finds only <article> (essay body) and gets text inside it
func extractArticleText(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	var inArticle bool

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "article" {
			inArticle = true
		}

		//	 Capture text nodes only when inside <article>
		if inArticle && n.Type == html.TextNode {
			trimmed := strings.TrimSpace(n.Data)
			if trimmed != "" {
				sb.WriteString(trimmed)
				sb.WriteRune(' ')
			}
		}

		// skip script/style
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}

		// Traverse child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		// If we reach the end of the article tag, stop capturing text
		if n.Type == html.ElementNode && n.Data == "article" {
			inArticle = false
		}
	}
	f(doc)

	return sb.String(), nil
}
