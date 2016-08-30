package markdown

import (
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/backlager/malt/ingredient"
)

var (
	regexInlineTitle    = regexp.MustCompile(`(?m:^(\#{1}\s*)(\w.+)$)`)
	regexInlineProperty = regexp.MustCompile(`(?m:^(\#{2}\s*)(\w.+)$)`)
	regexLineTitle      = regexp.MustCompile(`(?m:^(\w.+\r?\n)={2,}$)`)
	regexLineProperty   = regexp.MustCompile(`(?m:^(\w.+\r?\n)-{2,}$)`)
)

var (
	regexTitles = map[*regexp.Regexp]string{
		regexInlineTitle: "$2",
		regexLineTitle:   "$1",
	}
	regexProperties = map[*regexp.Regexp]string{
		regexInlineProperty: "$2",
		regexLineProperty:   "$1",
	}
)

// Markdown parser
type Markdown struct {
	breakpoints [][]int
}

// Parse the stream for markdown content
func (m Markdown) Parse(reader io.Reader) (ingredient.Ingredient, error) {
	recipe := &ingredient.Ingredient{}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return *recipe, err
	}

	for regex, replace := range regexTitles {
		if val := parseProperties(regex, b, replace); val != "" {
			recipe.Name = val
		}
	}

	for regex, replace := range regexProperties {
		if val := parseProperties(regex, b, replace); val != "" {
			if grain, ok := recipe.Grains[val]; !ok {
				grain = ingredient.Grain{
					Content: []string{},
				}
				if recipe.Grains == nil {
					recipe.Grains = make(map[string]ingredient.Grain)
				}
				recipe.Grains[val] = grain
			}
		}
	}

	return *recipe, nil
}

// Supports determines if a file is supported by the parser
func (m *Markdown) Supports(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".md")
}

func parseProperties(regex *regexp.Regexp, source []byte, replace string) string {
	if regex.Match(source) {
		return strings.TrimSpace(string(regex.ReplaceAll(regex.Find(source), []byte(replace))))
	}

	return ""
}
