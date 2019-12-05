package lexer

import (
	"github.com/pkg/errors"
	"strconv"
)

type Token struct {
	Kind  TokenKind
	Start int
	End   int
	Value string
}

type TokenKind int

// NAME -> keyword relationship
const (
	FRAGMENT     = "fragment"
	QUERY        = "query"
	MUTATION     = "mutation"
	SUBSCRIPTION = "subscription"
	SCHEMA       = "schema"
	SCALAR       = "scalar"
	TYPE         = "type"
	INTERFACE    = "interface"
	UNION        = "union"
	ENUM         = "enum"
	INPUT        = "input"
	EXTEND       = "extend"
	DIRECTIVE    = "directive"
)

const (
	EOF TokenKind = iota + 1
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

var tokenDescription = map[TokenKind]string{
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

func (kind TokenKind) String() string {
	return tokenDescription[kind]
}

type Lexer struct {
	source            string
	tokens            []Token
	currentTokenIndex int
}

func NewLexer(src string) (*Lexer, error) {
	lexer := &Lexer{
		source: src,
	}

	tokenizedDocument, err := lex(src)
	if err != nil {
		return nil, err
	}

	lexer.tokens = tokenizedDocument
	lexer.currentTokenIndex = 0

	return lexer, nil
}

func (l Lexer) Get() *Token {
	token := &l.tokens[l.currentTokenIndex]

	if l.currentTokenIndex < len(l.tokens) {
		l.currentTokenIndex++
	}

	return token
}

func (l Lexer) Current() *Token {
	return &l.tokens[l.currentTokenIndex]
}

func lex(doc string) ([]Token, error) {
	var (
		whiteSpaceOn  bool
		doubleQuoteOn bool
		tok           string
		kind          TokenKind
		tokens        []Token
	)

	for index, char := range doc {
		switch char {
		case '!':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{BANG, index, index, string(char)})
		case '$':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{DOLLAR, index, index, string(char)})
		case '(':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{PAREN_L, index, index, string(char)})
		case ')':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{PAREN_R, index, index, string(char)})
		case '.':
			tok = tok + string(char)
			if tok == "..." {
				tokens = append(tokens, Token{SPREAD, index, index, "..."})
				tok = ""
			} else if kind == INT {
				kind = FLOAT
			} else if kind == FLOAT {
				return nil, errors.New("unexpected character '.' at position " + strconv.Itoa(index))
			}
		case ':':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{COLON, index, index, string(char)})
		case '=':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{EQUALS, index, index, string(char)})
		case '@':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{AT, index, index, string(char)})
		case '[':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{BRACKET_L, index, index, string(char)})
		case ']':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{BRACKET_R, index, index, string(char)})
		case '{':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{BRACE_L, index, index, string(char)})
		case '}':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{BRACE_R, index, index, string(char)})
		case '|':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{PIPE, index, index, string(char)})
		case '&':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			tokens = append(tokens, Token{AMP, index, index, string(char)})
		case '"':
			whiteSpaceOn = false
			tok = tok + string(char)
			if !doubleQuoteOn {
				doubleQuoteOn = true
				kind = STRING
			} else {
				doubleQuoteOn = false
				tokens = append(tokens, Token{STRING, index - len(tok), index - 1, tok})
				tok = ""
			}
		case ' ':
			if doubleQuoteOn {
				tok = tok + string(char)
			} else {
				if !whiteSpaceOn {
					whiteSpaceOn = true
					if tok != "" && tok != "..." {
						tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
						tok = ""
					}
				}
			}
		case ',':
			whiteSpaceOn = false
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
			//tokens = append(tokens, Token{BANG, index, index, string(char)})
		case '\n', '\r':
			if tok != "" {
				tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})
				tok = ""
			}
		case '_', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
			'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
			'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			if tok == "" {
				kind = NAME
			}
			tok = tok + string(char)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if tok == "" {
				kind = INT
			}
			tok = tok + string(char)
		default:
			tok = tok + string(char)
		}
	}

	tokens = append(tokens, Token{EOF, len(doc), len(doc), EOF.String()})

	return tokens, nil
}
