package json

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseJSON(t *testing.T) {
	bueno := `
	{
		"title" : "This is un bueno title",
		"description" : "#Some description\n of stuff or whatever",
		"data" : {
			"hello" : "1"
		}
	}
	`

	tests := []struct {
		JSONContent string
		Property    string
		Expected    int
	}{
		{
			bueno,
			"title",
			1,
		},
		{
			bueno,
			"description",
			1,
		},
		{
			bueno,
			"dog",
			0,
		},
		{
			bueno,
			"data",
			0,
		},
	}

	json := &JSON{}
	for _, test := range tests {
		reader := bytes.NewReader([]byte(test.JSONContent))
		recipe, err := json.Parse(reader)
		if err != nil {
			t.Fatal(err)
		}

		grain := recipe.Grains[test.Property]
		if len(grain.Content) != test.Expected {
			t.Errorf("For `%s` expected %d but was %d with %v",
				test.Property,
				test.Expected,
				len(grain.Content),
				strings.Join(grain.Content, "\n"),
			)
		}
	}
}
