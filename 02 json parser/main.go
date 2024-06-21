package main

import (
	"fmt"
	"os"
)

// 1- Defining the JSON Lexer tokens
type TokenType int

const (
	NULL TokenType = iota
	TRUE
	FALSE
	NUMBER
	STRING
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	COMMA
	COLON
	EOF
	ILLEGAL
)

type Token struct {
	value     string
	tokenType TokenType
}

var mapToken = map[TokenType]string{
	NULL:     "NULL",
	TRUE:     "TRUE",
	FALSE:    "FALSE",
	NUMBER:   "NUMBER",
	STRING:   "STRING",
	LBRACE:   "LBRACE",
	RBRACE:   "RBRACE",
	LBRACKET: "LBRACKET",
	RBRACKET: "RBRACKET",
	COMMA:    "COMMA",
	COLON:    "COLON",
	EOF:      "EOF",
	ILLEGAL:  "ILLEGAL",
}

// 2- Defining the JSON Lexer structure
type JSONLexer struct {
	input  string // JSON input
	start  int    // Start position of the current token
	pos    int    // Current position in the input
	tokens []Token
}

// 3- Defining the JSON Lexer constructor
func NewJSONLexer(input string) *JSONLexer {
	return &JSONLexer{input: input, tokens: make([]Token, 0)}
}

// 4- Defining the JSON Lexer run method
func (l *JSONLexer) run() {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\n', '\r':
			l.ignore()
		case '{':
			l.emit(LBRACE, "{")
		case '}':
			l.emit(RBRACE, "}")
		case '[':
			l.emit(LBRACKET, "[")
		case ']':
			l.emit(RBRACKET, "]")
		case ',':
			l.emit(COMMA, ",")
		case ':':
			l.emit(COLON, ":")
		case 't', 'f', 'n':
			l.parseKeyword()
		case '"':
			l.parseString()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			l.parseNumber()
		default:
			l.emit(ILLEGAL, string(l.input[l.pos]))
		}
	}
	l.emit(EOF, "EOF")
}

// 5- Defining the JSON Lexer ignore method
func (l *JSONLexer) ignore() {
	l.start = l.pos
	l.pos++
}

// 6- Defining the JSON Lexer emit method
func (l *JSONLexer) emit(token TokenType, value string) {
	l.tokens = append(l.tokens, Token{value: value, tokenType: token})
	l.start = l.pos
	l.pos++
}

// 7- Defining the JSON Lexer parseKeyword method
func (l *JSONLexer) parseKeyword() {
	l.pos++
	for l.pos < len(l.input) && l.input[l.pos] >= 'a' && l.input[l.pos] <= 'z' {
		l.pos++
	}
	keyword := l.input[l.start+1 : l.pos]
	fmt.Println("Keyword found", keyword, len(keyword), string(l.input[l.pos]))
	switch keyword {
	case "true":
		l.emit(TRUE, "true")
	case "false":
		l.emit(FALSE, "false")
	case "null":
		l.emit(NULL, "null")
	// unknown keyword found
	default:
		l.emit(ILLEGAL, keyword)
	}
	l.pos--
}

func (l *JSONLexer) parseString() {
	l.pos++
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		l.pos++
	}
	str := l.input[l.start+2 : l.pos]
	l.emit(STRING, str)
}

func (l *JSONLexer) parseNumber() {
	for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
		l.pos++
	}
	num := l.input[l.start+1 : l.pos]
	l.emit(NUMBER, num)
	l.pos--
}

func main() {
	// input := `{"key": "value", "name": null, "bool": false, "key2": 123, "key3": true, "key4": [1, 2, 3, 4, 5]}`
	file, err := os.ReadFile("./tests/step3/invalid.json")
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	lexer := NewJSONLexer(string(file))
	lexer.run()
	for _, token := range lexer.tokens {
		fmt.Printf("Token: %s \t Type: %s\n", token.value, mapToken[token.tokenType])
	}

}

// create a generic stack
type Stack struct {
	stack []string
}

func NewStack() *Stack {
	return &Stack{stack: make([]string, 0)}
}

func (s *Stack) Push(value string) {
	s.stack = append(s.stack, value)
}

func (s *Stack) Pop() string {
	if len(s.stack) == 0 {
		return ""
	}
	value := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return value
}
