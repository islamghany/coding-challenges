package main

import (
	"fmt"
	"strconv"
)

// =====================================================================================
// Lexer

// Defining the JSON Lexer tokens
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

type JSONLexer struct {
	input  string // JSON input
	start  int    // Start position of the current token
	pos    int    // Current position in the input
	tokens []Token
}

func NewJSONLexer(input string) *JSONLexer {
	return &JSONLexer{input: input, tokens: make([]Token, 0)}
}

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
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			l.parseNumber()
		default:
			l.emit(ILLEGAL, string(l.input[l.pos]))
		}
	}
	l.emit(EOF, "")
}

func (l *JSONLexer) ignore() {
	l.start = l.pos
	l.pos++
}

func (l *JSONLexer) emit(token TokenType, value string) {
	l.tokens = append(l.tokens, Token{value: value, tokenType: token})
	l.start = l.pos
	l.pos++
}

func (l *JSONLexer) parseKeyword() {
	l.pos++
	for l.pos < len(l.input) && l.input[l.pos] >= 'a' && l.input[l.pos] <= 'z' {
		l.pos++
	}
	keyword := l.input[l.start+1 : l.pos]
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

// =====================================================================================
// Parser

type JSONParser struct {
	tokens []Token
	pos    int
}

// JSONValue represents bacis JSON types like string, number, boolean, null
type JSONValue interface{}

// JSONObject represents a JSON object
type JSONObject map[string]JSONValue

// JSONArray represents a JSON array
type JSONArray []JSONValue

// Paser Implementation

func NewJSONParser(tokens []Token) *JSONParser {
	return &JSONParser{tokens: tokens}
}

func (p *JSONParser) Parse() (JSONValue, error) {
	return p.parseValue()
}

// parseValue parses the JSON value
func (p *JSONParser) parseValue() (JSONValue, error) {
	token := p.tokens[p.pos]
	p.pos++
	switch token.tokenType {
	case LBRACE:
		return p.parseObject()
	case LBRACKET:
		return p.parseArray()
	case STRING:
		return token.value, nil
	case NUMBER:
		num, err := strconv.ParseFloat(token.value, 64)
		if err != nil {
			return nil, err
		}
		return num, nil
	case TRUE:
		return true, nil
	case FALSE:
		return false, nil
	case NULL:
		return nil, nil
	default:
		return nil, fmt.Errorf("Unexpected token %s", token.value)
	}
}

// parseObject parses the JSON object
// An object is an unordered set of name/value pairs. An object begins with {left brace and ends with }right brace. Each name is followed by :colon and the name/value pairs are separated by ,comma.
func (p *JSONParser) parseObject() (JSONObject, error) {
	obj := JSONObject{}
	prevComma := false
	for {
		// check if the object is empty before we hit the right brace
		if p.end() {
			return nil, fmt.Errorf("unexpected end of tokens")
		}
		nextToken := p.tokens[p.pos]
		// check if we hit the right brace
		if nextToken.tokenType == RBRACE {
			// check if there is a comma before the right brace
			if prevComma {
				return nil, fmt.Errorf("unexpected comma")
			}
			p.pos++
			break
		}
		// check if the next token is a string key otherwise return an error
		if nextToken.tokenType != STRING {
			return nil, fmt.Errorf("expected string key")
		}
		p.pos++
		if p.end() || p.tokens[p.pos].tokenType != COLON {
			return nil, fmt.Errorf("expected colon")
		}
		p.pos++
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[nextToken.value] = value
		if p.end() {
			return nil, fmt.Errorf("unexpected end of tokens")
		}
		prevComma = false
		if p.tokens[p.pos].tokenType == COMMA {
			p.pos++
			prevComma = true
		}

	}
	return obj, nil

}
func (p *JSONParser) parseArray() (JSONArray, error) {
	arr := JSONArray{}
	prevComma := false
	for {
		if p.end() {
			return nil, fmt.Errorf("unexpected end of tokens")
		}
		nextToken := p.tokens[p.pos]
		if nextToken.tokenType == RBRACKET {
			if prevComma {
				return nil, fmt.Errorf("unexpected comma")
			}
			p.pos++
			break
		}
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
		if p.end() {
			return nil, fmt.Errorf("unexpected end of tokens")
		}
		prevComma = false
		if p.tokens[p.pos].tokenType == COMMA {
			p.pos++
			prevComma = true
		}
	}
	return arr, nil
}

func (p *JSONParser) end() bool {
	return p.pos >= len(p.tokens)
}

func main() {
	// input := `{"key": "value", "name": null, "bool": false, "key2": 123,"arr":[1,2, "islam",{"name":"islam"}], "key3": true, "data":{"name":"islam","age":21}}`

	// file, err := os.ReadFile("./tests/step3/invalid.json")
	// if err != nil {
	// 	fmt.Println("Error reading file", err)
	// 	return
	// }
	lexer := NewJSONLexer(`{"key": "value", "s":[1,2,3,]}`)
	lexer.run()
	// for _, token := range lexer.tokens {
	// 	fmt.Printf("Token: %s \t Type: %s\n", token.value, mapToken[token.tokenType])
	// }
	parser := NewJSONParser(lexer.tokens)
	fmt.Println(parser.Parse())

}
