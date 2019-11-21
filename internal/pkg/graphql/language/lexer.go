package language

import (
	"errors"
	"strings"
)

// tokenKind is a type to represent types of tokens
type tokenKind int

// Token Kinds
const (
	EOF tokenKind = iota + 1
	BANG
	DOLLAR
	PAREN_L
	PAREN_R
	SPREAD
	COLON
	EQUALS
	AT
	BRACKET_L
	BRACKET_R
	BRACE_L
	PIPE
	BRACE_R
	NAME
	INT
	FLOAT
	STRING
	BLOCK_STRING
	AMP
)

type keyword string

// Keywords
const (
	FRAGMENT     keyword = "fragment"
	QUERY        keyword = "query"
	MUTATION     keyword = "mutation"
	SUBSCRIPTION keyword = "subscription"
	SCHEMA       keyword = "schema"
	SCALAR       keyword = "scalar"
	TYPE         keyword = "type"
	INTERFACE    keyword = "interface"
	UNION        keyword = "union"
	ENUM         keyword = "enum"
	INPUT        keyword = "input"
	EXTEND       keyword = "extend"
	DIRECTIVE    keyword = "directive"
)

// Token is a struct that holds details about a token
type Token struct {
	Kind  tokenKind
	Start int
	End   int
	Value string
}

// Descriptions of all token kinds
var tokenDescription = map[tokenKind]string{
	EOF:          "EOF",
	BANG:         "!",
	DOLLAR:       "$",
	PAREN_L:      "(",
	PAREN_R:      ")",
	SPREAD:       "...",
	COLON:        ":",
	EQUALS:       "=",
	AT:           "@",
	BRACKET_L:    "[",
	BRACKET_R:    "]",
	BRACE_L:      "{",
	PIPE:         "|",
	BRACE_R:      "}",
	NAME:         "Name",
	INT:          "Int",
	FLOAT:        "Float",
	STRING:       "String",
	BLOCK_STRING: "BlockString",
	AMP:          "&",
}

func (tk tokenKind) String() string {
	return tokenDescription[tk]
}

func createToken(start, end int, kind tokenKind, value string) Token {
	return Token{
		Start: start,
		End:   end,
		Kind:  kind,
		Value: value,
	}
}

// readSpread reads a spread token from the document and
// returns the value of the token and the end index of it
func readSpread(doc string, index int) (string, int, error) {
	var tokenVal strings.Builder

	if index+2 > len(doc) {
		return "", index, errors.New("End of document reached")
	}

	for i := 0; i < index+3; i++ {
		tokenVal.WriteByte(doc[index+i])
	}

	tokenStr := tokenVal.String()

	if tokenStr != "..." {
		return "", index, errors.New("Not a spread token")
	}

	return tokenStr, index + 2, nil
}

func Lex(doc string) ([]Token, error) {
	var (
		tokens []Token
	)
	for index := 0; index < len(doc); index++ {
		switch doc[index] {
		case '!', '$', '(', ')', ':', '=',
			'@', '[', ']', '{', '}', '|', '&':
			token := createToken(index, index+1,
				BANG, string(doc[index]))

			tokens = append(tokens, token)
			break
		case '.':
			tokenVal, newIndex, err := readSpread(doc, index)

			if err != nil {
				return nil, err
			}

			token := createToken(index, newIndex,
				SPREAD, tokenVal)

			index = newIndex

			tokens = append(tokens, token)
			break
		}
	}

	return tokens, nil
}
