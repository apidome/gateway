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
		stringOn      bool
		blockStringOn bool
		tok           string
		kind          TokenKind
		tokens        []Token
	)

	runes := []rune(doc)

	//for index, char := range doc {
	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		// This case handles all the punctuator characters that stand as an independent token.
		case '!', '$', '(', ')', ':', '=', '@', '[', ']', '{', '}', '|', '&':
			{
				// Turn off the white space flag.
				whiteSpaceOn = false

				// If the token is not empty, it means that the current character
				// is a punctuator that ended a token.
				if tok != "" {
					// Append the token token to the tokens slice.
					tokens = append(tokens, Token{kind, i - len(tok), i - 1, tok})

					// Empty the token.
					tok = ""
				}

				// Append a punctuator token to the tokens slice.
				tokens = append(tokens, Token{getPunctuatorKind(runes[i]), i, i, string(runes[i])})
			}
		// This case handles the "spread" operator and float's decimal point.
		case '.':
			{
				// If the kind of the token is INT, set it to FLOAT.
				if kind == INT {
					kind = FLOAT
				} else if kind == FLOAT {
					// If we met '.' when the token kind is FLOAT, return an error.
					return nil, errors.New("unexpected character '.' at position " + strconv.Itoa(i))
				}

				// Append the character to the token.
				tok = tok + string(runes[i])

				// Check if the token is a spread operator.
				if tok == "..." {
					// Append the token to the tokens slice.
					tokens = append(tokens, Token{SPREAD, i, i, "..."})

					// Empty the token.
					tok = ""
				}
			}
		// This case is responsible for detection of STRING and BLOCK_STRING values.
		case '"':
			{
				// Turn off the white space flag.
				whiteSpaceOn = false

				// If the token is empty, this is a start of a STRING or a BLOCK_STRING
				// Else, this is the end of a value.
				if tok == "" {
					// If the next two characters are also double quotes, this is a BLOCK_STRING.
					// Else, this is a STRING.
					if runes[i+1] == '"' && runes[i+2] == '"' {
						// Turn the block string flag on.
						blockStringOn = true

						// Set the token kind to block string.
						kind = BLOCK_STRING

						// Append the '"' runes to the token.
						tok = tok + string(runes[i:i+3])

						// Increment the loop index twice.
						i += 2
					} else {
						// Turn the string flag on
						stringOn = true

						// Set the token king to string
						kind = STRING

						// Append the '"' rune to the token.
						tok = tok + string(runes[i])
					}
				} else {
					if blockStringOn {
						if runes[i+1] == '"' && runes[i+2] == '"' {
							// Turn the block string flag off.
							blockStringOn = false

							// Append the '"' runes to the token.
							tok = tok + string(runes[i:i+3])

							// Increment the loop index twice.
							i += 2
						} else {
							// Append the '"' as part of the value.
							tok = tok + string(runes[i])
						}
					} else if stringOn {
						// Turn the string flag off
						stringOn = false

						// Append the '"' as the end of the string value.
						tok = tok + string(runes[i])
					} else {
						return nil, errors.New("unexpected character '\"' at position " + strconv.Itoa(i))
					}

					// Append the token to the tokens slice.
					tokens = append(tokens, Token{kind, i - len(tok), i - 1, tok})

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
				if stringOn {
					tok = tok + string(runes[i])
				} else {
					// If the white space flag is off, turn it on.
					if !whiteSpaceOn {
						whiteSpaceOn = true

						// If the token is not empty
						if tok != "" {
							// Append the token to tokens slice.
							tokens = append(tokens, Token{kind, i - len(tok), i - 1, tok})

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
					tokens = append(tokens, Token{kind, i - len(tok), i - 1, tok})

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
					tokens = append(tokens, Token{kind, i - len(tok), i - 1, tok})

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
				tok = tok + string(runes[i])
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			{
				// If the current character (which is a digit) is the first character
				// of the token, set the kind as an INT.
				if tok == "" {
					kind = INT
				}

				// Append the character to the token variable.
				tok = tok + string(runes[i])
			}
		// This case handles the rest of the unicode characters (hebrew, chinese, etc)
		default:
			{
				tok = tok + string(runes[i])
			}
		}
	}

	// Append an EOF token to mark the end of the document.
	tokens = append(tokens, Token{EOF, len(doc), len(doc), EOF.String()})

	// Return the tokens slice.
	return tokens, nil
}
