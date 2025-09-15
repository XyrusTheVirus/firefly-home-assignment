package sorts

import "sort"

// SortByFrequency takes a map of words and their counts and returns a slice of structs sorted by frequency.
func SortByFrequency(data map[string]interface{}) []struct {
	Word  string
	Count int
} {
	var list []struct {
		Word  string
		Count int
	}

	// Convert the map to a slice of structs
	for w, c := range data {
		list = append(list, struct {
			Word  string
			Count int
		}{
			Word:  w,
			Count: c.(int),
		},
		)
	}
	sort.Slice(list, func(i, j int) bool {
		// Sort by count in descending order, and by word in ascending order for ties
		if list[i].Count == list[j].Count {
			return list[i].Word < list[j].Word
		}
		return list[i].Count > list[j].Count
	})
	return list
}
