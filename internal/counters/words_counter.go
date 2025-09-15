package counters

import (
	"sort"
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

// TopN returns the top N words sorted by frequency.
func (wc *WordsCounter) TopN(n int) []struct {
	Word  string
	Count int
} {
	var list []struct {
		Word  string
		Count int
	}
	// Convert the map to a slice of structs
	for w, c := range wc.Wc {
		list = append(list, struct {
			Word  string
			Count int
		}{
			Word:  w,
			Count: c.(int),
		},
		)
	}
	// Sort the list by count in descending order, and by word in ascending order for ties
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count == list[j].Count {
			return list[i].Word < list[j].Word
		}
		return list[i].Count > list[j].Count
	})
	// Return only the top N words
	if len(list) > n {
		return list[:n]
	}
	return list
}
