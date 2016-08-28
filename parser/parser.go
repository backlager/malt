package parser

import (
	"errors"
	"io"
	"os"
	"time"

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
		return nil, err
	}

	// Open the file
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Parse the contents
	return p.Parse(f)
}

// Check to see if the file exists and is accessible
// If it does not exist we create it
func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		_, err = os.Create(fileName)
		if err != nil {
			return false
		}
	}

	return true
}

// Parse the given string to extract a proper date
func parseDate(in string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"2006-1-2",
		"2006/1/2",
		"01-02-2006",
		"01/02/2006",
		"1-2-2006",
		"1/2/2006",
		"Jan 2, 2006",
		"Jan 02, 2006",
		"2 Jan 2006",
		"02 Jan 2006",
	}

	for _, f := range formats {
		d, err := time.Parse(f, in)
		if err == nil {
			return d, nil
		}
	}

	return time.Now().UTC(), errors.New("No valid date provided")
}
