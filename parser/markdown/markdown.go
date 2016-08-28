package markdown

import (
	"io"
	"strings"

	"github.com/backlager/malt/ingredient"
)

// Markdown parser
type Markdown struct{}

// Parse the stream for markdown content
func (m *Markdown) Parse(reader io.Reader) (ingredient.Ingredient, error) {
	// TODO: Parse markdown
	return nil, nil
}

// Supports determines if a file is supported by the parser
func (m *Markdown) Supports(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".md")
}
