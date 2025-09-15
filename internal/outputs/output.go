package outputs

type OutputInterface interface {
	Print()
}

type Output struct {
	Data []struct {
		Word  string
		Count int
	}
}

// TopN returns the top N words by frequency. If there are fewer than N words, it returns all of them.
func (o Output) TopN(n int) []struct {
	Word  string
	Count int
} {
	// Return n entries if there are more than n, otherwise return all
	if len(o.Data) > n {
		return o.Data[:n]
	}
	return o.Data
}
