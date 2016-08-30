package markdown

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/backlager/malt/ingredient"
)

var (
	regexInlineTitle    = regexp.MustCompile(`(?m:^(\#{1}\s*)(\w.+)$)`)
	regexInlineProperty = regexp.MustCompile(`(?m:^(\#{2}\s*)(\w.+)$)`)
	regexLineTitle      = regexp.MustCompile(`(?m:^(\w.+\r?\n)={2,}$)`)
	regexLineProperty   = regexp.MustCompile(`(?m:^(\w.+\r?\n)-{2,}$)`)
	regexBlock          = regexp.MustCompile(`(?m:(.*?[^\:\-\,])(?:$|\n{2,}))`)
	regexBlockHeading   = regexp.MustCompile(`^([#]{1,}\s+)(.*)`)
	regexBlockUnordered = regexp.MustCompile(`^([\-\+\*]{1}\s+)(.*)`)
	regexBlockQuote     = regexp.MustCompile(`^([\>]{1}\s+)(.*)`)
	regexBlockOrdered   = regexp.MustCompile(`^(\d{1,}\.\s+)(.*)`)
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
func (m *Markdown) Parse(reader io.Reader) (ingredient.Ingredient, error) {
	recipe := &ingredient.Ingredient{
		Name:   "",
		Grains: map[string]ingredient.Grain{},
	}

	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return *recipe, err
	}

	for regex, replace := range regexTitles {
		titles, _ := parseProperties(regex, buffer, replace)
		for _, val := range titles {
			recipe.Name = val
			break
		}
	}

	postNames := []string{}
	postBounds := [][]int{}
	for regex, replace := range regexProperties {
		properties, bounds := parseProperties(regex, buffer, replace)
		for i, val := range properties {
			if grain, ok := recipe.Grains[val]; !ok {
				grain = ingredient.Grain{
					Content: []string{},
				}

				recipe.Grains[val] = grain
				postNames = append(postNames, val)
				postBounds = append(postBounds, bounds[i])
			}
		}
	}

	length := len(buffer)
	total := len(postBounds)
	for i := 0; i < total; i++ {
		// Get the outer of the current property
		start := postBounds[i][1]

		// Start with the length of the entire document
		stop := length
		if i+1 < total {
			// If there is another property then use the inner of it
			// as the stopping point
			stop = postBounds[i+1][0]
		}

		// TODO: This should not be necessary, but sometimes the bounds come in screwed up.
		if stop < start {
			temp := start
			start = stop
			stop = temp
		}

		// Retrieve the grain and create a scope buffer to be used
		// for parsing the associated content
		prop := postNames[i]
		scope := buffer[start:stop]
		grain := recipe.Grains[prop]
		grain.Content = parseBlock(regexBlock, scope)
		recipe.Grains[prop] = grain
	}

	return *recipe, nil
}

// Supports determines if a file is supported by the parser
func (m *Markdown) Supports(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".md")
}

func parseProperties(regex *regexp.Regexp, source []byte, replace string) ([]string, [][]int) {
	properties := []string{}
	breakpoints := [][]int{}

	if regex.Match(source) {
		for i, match := range regex.FindAll(source, -1) {
			bp := regex.FindAllIndex(source, -1)[i]
			breakpoints = append(breakpoints, bp)

			prop := string(bytes.TrimSpace(regex.ReplaceAll(match, []byte(replace))))
			properties = append(properties, prop)
		}
	}

	return properties, breakpoints
}

func parseBlock(regex *regexp.Regexp, source []byte) []string {
	blocks := []string{}

	if regex.Match(source) {
		for _, match := range regex.FindAll(source, -1) {
			block := bytes.TrimSpace(match)

			if len(block) > 0 {
				if len(blocks) == 0 {
					blocks = append(blocks, string(block))
					continue
				}

				// Special handling for
				//  - unordered list items (-, +, *)
				//  - ordered list items (-, +, *)
				//  - headers (#)
				//  - blockquotes (>)
				if regexBlockUnordered.Match(block) {
					blocks = append(blocks, string(regexBlockUnordered.ReplaceAll(block, []byte("$2"))))
				} else if regexBlockOrdered.Match(block) {
					blocks = append(blocks, string(regexBlockOrdered.ReplaceAll(block, []byte("$2"))))
				} else if regexBlockHeading.Match(block) {
					blocks = append(blocks, string(regexBlockHeading.ReplaceAll(block, []byte("$2"))))
				} else if regexBlockQuote.Match(block) {
					blocks = append(blocks, string(regexBlockQuote.ReplaceAll(block, []byte("$2"))))
				} else {
					blocks[len(blocks)-1] = fmt.Sprintf("%s\n%s", blocks[len(blocks)-1], string(block))
				}

				log.Printf("%v", blocks)
			}
		}
	}

	return blocks
}
