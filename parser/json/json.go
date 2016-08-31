package json

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/backlager/malt/ingredient"
)

// JSON parser
type JSON struct{}

// Parse the stream for json content
func (j *JSON) Parse(reader io.Reader) (ingredient.Ingredient, error) {
	// Initialize map
	ing := ingredient.Ingredient{
		Grains: make(map[string]ingredient.Grain),
	}

	// Initialize map
	props := make(map[string]interface{})

	// Decode the reader. If something goes wrong, it's most likely bad json
	dec := json.NewDecoder(reader)
	if err := dec.Decode(&props); err != nil {
		return ing, err
	}

	// Iteratate through obtained properties and add them to the Grains
	for k := range props {

		// Ensure that the top level key maps to a string value
		val, ok := props[k].(string)
		if !ok {
			fmt.Printf("Skipping property `%s` because it contains a value that is not of type string. Value is of type `%T`\n", k, props[k])
			continue
		}
		// TODO How are we setting the name of the ingredient?
		if k == "title" {
			ing.Name = val
		}
		ing.Grains[k] = ingredient.Grain{
			Content: []string{val},
		}
	}
	return ing, nil
}

// Supports determines if a file is supported by the parser
func (j *JSON) Supports(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".json")
}
