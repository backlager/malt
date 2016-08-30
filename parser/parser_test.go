package parser

import "testing"

func TestGetParser(t *testing.T) {
	filePath := "../.brewery/grain/parser_for_markdown.md"
	if _, err := GetParser(filePath, AvailableParsers...); err != nil {
		t.Fail()
	}
}

func TestRead(t *testing.T) {
	filePath := "../.brewery/grain/parser_for_markdown.md"
	if recipe, err := Read(filePath); err != nil {
		t.Logf("%s", err)
		t.Fail()
	} else {
		if recipe.Name == "" {
			t.Logf("%v", recipe)
			t.Fail()
		}
	}
}
