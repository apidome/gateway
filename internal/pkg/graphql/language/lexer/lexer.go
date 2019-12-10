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

var punctuatorKindsMap = map[rune]TokenKind{
	'!': BANG,
	'$': DOLLAR,
	'(': PAREN_L,
	')': PAREN_R,
	':': COLON,
	'=': EQUALS,
	'@': AT,
	'[': BRACKET_L,
	']': BRACKET_R,
	'{': BRACE_L,
	'}': BRACE_R,
	'|': PIPE,
	'&': AMP,
}

func getPunctuatorKind(c rune) TokenKind {
	return punctuatorKindsMap[c]
}

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

func (l *Lexer) Get() *Token {
	token := &l.tokens[l.currentTokenIndex]

	if l.currentTokenIndex < len(l.tokens) {
		l.currentTokenIndex++
	}

	return token
}

func (l *Lexer) Current() *Token {
	return &l.tokens[l.currentTokenIndex]
}

func (l *Lexer) Source() *string {
	return &l.source
}

func lex(doc string) ([]Token, error) {
	var (
		whiteSpaceOn  bool
		doubleQuoteOn bool
		blockStringOn bool
		tok           string
		kind          TokenKind
		tokens        []Token
	)

	for index, char := range doc {
		switch char {
		// This case handles all the punctuator characters that stand as an independent token.
		case '!', '$', '(', ')', ':', '=', '@', '[', ']', '{', '}', '|', '&':
			{
				// Turn off the white space flag.
				whiteSpaceOn = false

				// If the token is not empty, it means that the current character
				// is a punctuator that ended a token.
				if tok != "" {
					// Append the token token to the tokens slice.
					tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})

					// Empty the token.
					tok = ""
				}

				// Append a punctuator token to the tokens slice.
				tokens = append(tokens, Token{getPunctuatorKind(char), index, index, string(char)})
			}
		// This case handles the "spread" operator and float's decimal point.
		case '.':
			{
				// If the kind of the token is INT, set it to FLOAT.
				if kind == INT {
					kind = FLOAT
				} else if kind == FLOAT {
					// If we met '.' when the token kind is FLOAT, return an error.
					return nil, errors.New("unexpected character '.' at position " + strconv.Itoa(index))
				}

				// Append the character to the token.
				tok = tok + string(char)

				// Check if the token is a spread operator.
				if tok == "..." {
					// Append the token to the tokens slice.
					tokens = append(tokens, Token{SPREAD, index, index, "..."})

					// Empty the token.
					tok = ""
				}
			}
		case '"':
			{
				// Turn off the white space flag.
				whiteSpaceOn = false

				// Append the character to the token variable
				tok = tok + string(char)

				// If double quote flag is off, it means that it is the start of a new STRING.
				// Else, it is the end of a STRING.
				if !doubleQuoteOn {
					// Turn on the double quote flag.
					doubleQuoteOn = true

					// Set the token kind to string.
					kind = STRING
				} else {
					// Turn off the double quote flag.
					doubleQuoteOn = false

					// Append the token to the tokens slice.
					tokens = append(tokens, Token{STRING, index - len(tok), index - 1, tok})

					// Empty the token.
					tok = ""
				}
			}
		// This case handles white spaces (single white space character only)
		case ' ':
			{
				// If the double quote flag is on, we append any white space
				// we meet to the token variable.
				// Else, we will treat the white space as a delimiter between tokens
				if doubleQuoteOn {
					tok = tok + string(char)
				} else {
					// If the white space flag is off, turn it on.
					if !whiteSpaceOn {
						whiteSpaceOn = true

						// If the token is not empty
						if tok != "" {
							// Append the token to tokens slice.
							tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})

							// Empty the token
							tok = ""
						}
					}
				}
			}
		// This case handles commas.
		case ',':
			{
				// Turn off the white space flag.
				whiteSpaceOn = false

				// If the token is not empty
				if tok != "" {
					// Append the token to the tokens slice.
					tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})

					// Empty the token
					tok = ""
				}
			}
		// This case handles line feed, carriage return and tab.
		case '\n', '\r', '\t':
			{
				// If the token is not empty, it means that the current character ends a token.
				if tok != "" {
					// Append the token to the tokens slice.
					tokens = append(tokens, Token{kind, index - len(tok), index - 1, tok})

					// Empty the token.
					tok = ""
				}
			}
		case '_', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
			'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
			'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			{
				// If the current character is the first character of the token,
				// set the token kind to NAME.
				if tok == "" {
					kind = NAME
				}

				// Append the character to the token variable.
				tok = tok + string(char)
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			{
				// If the current character (which is a digit) is the first character
				// of the token, set the kind as an INT.
				if tok == "" {
					kind = INT
				}

				// Append the character to the token variable.
				tok = tok + string(char)
			}
		// This case handles the rest of the unicode characters (hebrew, chinese, etc)
		default:
			{
				tok = tok + string(char)
			}
		}
	}

	// Append an EOF token to mark the end of the document.
	tokens = append(tokens, Token{EOF, len(doc), len(doc), EOF.String()})

	// Return the tokens slice.
	return tokens, nil
}
