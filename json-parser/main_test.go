package main

import (
	"os"
	"testing"
)

func TestJSONLexer(t *testing.T) {
	lexer := NewJSONLexer(`{"key": "value"}`)

	lexer.run()

	expected := []Token{
		{value: "{", tokenType: LBRACE},
		{value: "key", tokenType: STRING},
		{value: ":", tokenType: COLON},
		{value: "value", tokenType: STRING},
		{value: "}", tokenType: RBRACE},
		{value: "", tokenType: EOF},
	}

	for i, token := range lexer.tokens {
		if token != expected[i] {
			t.Errorf("Expected %v, got %v", expected[i], token)
		}
	}
}

func TestJSONParser(t *testing.T) {

	tc := []struct {
		filename string
		valid    bool
	}{
		{"step1/valid", true},
		{"step1/invalid", false},
		{"step2/valid", true},
		{"step2/valid2", true},
		{"step2/invalid", false},
		{"step2/invalid2", false},
		{"step3/valid", true},
		{"step3/invalid", false},
		{"step4/valid", true},
		{"step4/invalid", false},
		{"step4/valid2", true},
	}

	for _, c := range tc {
		path := "./tests/" + c.filename + ".json"
		file, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		lex := NewJSONLexer(string(file))
		lex.run()
		parser := NewJSONParser(lex.tokens)
		_, err = parser.Parse()
		if c.valid && err != nil {
			t.Errorf("Expected valid json, got error: %v, path: %s", err, path)
		}
		if !c.valid && err == nil {
			t.Errorf("Expected invalid json, got no error, path: %s", path)
		}
	}
}
