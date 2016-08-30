package markdown

import (
	"bytes"
	"testing"
)

func TestParseLineTitle(t *testing.T) {
	good := `
this is a title
===============

and some context`

	reader := bytes.NewReader([]byte(good))
	ingredient, err := Markdown{}.Parse(reader)

	if err != nil || ingredient.Name != "this is a title" {
		t.Fail()
	}
}

func TestParseInlineTitle(t *testing.T) {
	good := `
# this is a title
and some context`

	reader := bytes.NewReader([]byte(good))
	ingredient, err := Markdown{}.Parse(reader)

	if err != nil || ingredient.Name != "this is a title" {
		t.Fail()
	}
}

func TestParseProperties(t *testing.T) {
	good := `
# this is a title
and some context

## Property 1
Property 1 value

Property 2
----------

Property 2 value

Invalid property
-
astastast

### Invalid Property
asteastast`

	reader := bytes.NewReader([]byte(good))
	ingredient, err := Markdown{}.Parse(reader)

	if err != nil || ingredient.Grains == nil || len(ingredient.Grains) != 2 {
		t.Logf("%v", ingredient)
		t.Fail()
	}
}
