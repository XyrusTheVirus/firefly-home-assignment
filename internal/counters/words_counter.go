package counters

import (
	"sync"
)

// Mutex for synchronizing access to the map
var mu = &sync.Mutex{}

type WordsCounter struct {
	Wc map[string]interface{}
}

// Count Responsible for counting the occurrences of words in the provided slice.
func (wc *WordsCounter) Count(words []string, wordsBank map[string]interface{}) {
	for _, word := range words {
		if _, ok := wordsBank[word]; ok {
			// Lock the map for safe concurrent access
			mu.Lock()
			// If the word is not in the map, initialize it to 1
			if _, exists := wc.Wc[word]; !exists {
				wc.Wc[word] = 1
			} else {
				wc.Wc[word] = wc.Wc[word].(int) + 1
			}
			mu.Unlock()
		}
	}
}
