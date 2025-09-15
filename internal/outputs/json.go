package outputs

import (
	"encoding/json"
	"firefly-home-assigment/configs"
	"log"
)

type JSONOutput struct {
	Output
}

// Print outputs the top N words in JSON format, falling back to regular print on error.
func (o *JSONOutput) Print() {
	topN := o.TopN(configs.EnvInt("TOP_NUMBER_OF_WORDS", "10"))
	out, err := json.MarshalIndent(topN, "", "  ")
	// Fallback to regular print if JSON marshalling fails
	if err != nil {
		log.Printf("Couldn't parse json, due to: %s, printing Regulary", err.Error())
		log.Println(topN)
		return
	}
	log.Println(string(out))
}
