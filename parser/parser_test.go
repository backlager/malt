package parser

import "testing"

func TestGetParser2(t *testing.T) {
	filePaths := []string{
		"../.brewery/grain/parser_for_markdown.md",
		"../.brewery/grain/parser_for_json.json",
	}

	for _, filePath := range filePaths {
		if _, err := GetParser(filePath, AvailableParsers...); err != nil {
			t.Fail()
		}
	}
}

func TestRead(t *testing.T) {
	filePaths := []string{
		"../.brewery/grain/parser_for_markdown.md",
		"../.brewery/grain/parser_for_json.json",
	}

	for _, filePath := range filePaths {
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
}
