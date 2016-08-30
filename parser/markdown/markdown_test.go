package markdown

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseLineTitle(t *testing.T) {
	good := `
this is a title
===============

don't use this title
====================

and some context`

	md := &Markdown{}
	reader := bytes.NewReader([]byte(good))
	recipe, err := md.Parse(reader)

	if err != nil || recipe.Name != "this is a title" {
		t.Fail()
	}
}

func TestParseInlineTitle(t *testing.T) {
	good := `
# this is a title
and some context in a really
long format that should be grouped together
- item 1
- item 2

* item a
* item b

# don't use this title
`

	md := &Markdown{}
	reader := bytes.NewReader([]byte(good))
	recipe, err := md.Parse(reader)

	if err != nil || recipe.Name != "this is a title" {
		t.Fail()
	}
}

func TestParseProperties(t *testing.T) {
	good := `
# this is a title
and some context

## Property 1
Property 1 value

## Property 1b

Prop 1b value

Property 2
----------

Property 2 value

Invalid property
-
astastast

### Invalid Property
asteastast`

	md := &Markdown{}
	reader := bytes.NewReader([]byte(good))
	recipe, err := md.Parse(reader)

	if err != nil || recipe.Grains == nil || len(recipe.Grains) != 3 {
		t.Fail()
	}
}

func TestBlocks(t *testing.T) {
	good := `
# this is a title
and some context

## Property 1
Property 1 value

## Property 1b

Prop 1b value
and one more line

Property 2
----------

Property 2 value

- one
- two
- three`

	md := &Markdown{}
	reader := bytes.NewReader([]byte(good))
	recipe, err := md.Parse(reader)

	if err != nil {
		t.Fail()
	}

	tests := map[string]int{
		"Property 1":  1,
		"Property 1b": 1,
		"Property 2":  4,
	}

	for prop, expected := range tests {
		grain := recipe.Grains[prop]
		if len(grain.Content) != expected {
			t.Errorf("For `%s` expected %d but was %d with %v", prop, expected, len(grain.Content), strings.Join(grain.Content, "\n"))
		}
	}
}
