package extractors

type ExtractorInterface interface {
	// Extract processes the input data and populates the Result field.
	Extract()
}

// Extractor serves as a base struct for different types of data extractors.
type Extractor struct {
	Result map[string]interface{}
}
