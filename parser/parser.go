package parser

import (
	"errors"
	"io"
	"os"

	"github.com/backlager/malt/ingredient"
	"github.com/backlager/malt/parser/markdown"
)

// Parser is a stream for parsing work unit files
type Parser interface {
	Parse(reader io.Reader) (ingredient.Ingredient, error)
	Supports(filePath string) bool
}

// AvailableParsers is an array of default parsers within the package
var AvailableParsers = []Parser{
	new(markdown.Markdown),
}

// GetParser helps find the appropiate parser for a file
func GetParser(fileName string, parsers ...Parser) (Parser, error) {
	if !fileExists(fileName) {
		return nil, errors.New("File does not exist")
	}

	// Loop through the parsers doing a quick sanity check
	// to see if they can handle the file type. Take the first
	// supported parser found.
	for _, p := range parsers {
		if p.Supports(fileName) {
			return p, nil
		}
	}

	return nil, errors.New("File not supported by available parsers")
}

// Read the stream with the appropriate parser
func Read(filePath string) (ingredient.Ingredient, error) {
	// Try and find a suitable parser for the file
	p, err := GetParser(filePath, AvailableParsers...)
	if err != nil {
		return ingredient.Ingredient{}, err
	}

	// Open the file
	f, err := os.Open(filePath)
	if err != nil {
		return ingredient.Ingredient{}, err
	}
	defer f.Close()

	// Parse the contents
	return p.Parse(f)
}

// Check to see if the file exists and is accessible
func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return false
	}

	return true
}
