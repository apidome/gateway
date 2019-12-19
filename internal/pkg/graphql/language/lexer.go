package language

import (
	"strconv"

	"github.com/pkg/errors"
)

type token struct {
	Kind  tokenKind
	Start int
	End   int
	Value string
}

type tokenKind int

// NAME -> keyword relationship
const (
	kwFragment     = "fragment"
	kwQuery        = "query"
	kwMutation     = "mutation"
	kwSubscription = "subscription"
	kwSchema       = "schema"
	kwScalar       = "scalar"
	kwType         = "type"
	kwInterface    = "interface"
	kwUnion        = "union"
	kwEnum         = "enum"
	kwInput        = "input"
	kwExtend       = "extend"
	kwDirective    = "directive"
	kwImplements   = "implements"
	kwOn           = "on"
	kwTrue         = "true"
	kwFalse        = "false"
	kwNull         = "null"
)

const (
	tokEOF tokenKind = iota + 1
	tokBang
	tokDollar
	tokParenL
	tokParenR
	tokSpread
	tokColon
	tokEquals
	tokAt
	tokBracketL
	tokBracketR
	tokBraceL
	tokPipe
	tokBraceR
	tokName
	tokInt
	tokFloat
	tokString
	tokBlockString
	tokAmp
)

var punctuatorKindsMap = map[rune]tokenKind{
	'!': tokBang,
	'$': tokDollar,
	'(': tokParenL,
	')': tokParenR,
	':': tokColon,
	'=': tokEquals,
	'@': tokAt,
	'[': tokBracketL,
	']': tokBracketR,
	'{': tokBraceL,
	'}': tokBraceR,
	'|': tokPipe,
	'&': tokAmp,
}

func getPunctuatorKind(c rune) tokenKind {
	return punctuatorKindsMap[c]
}

var tokenDescription = map[tokenKind]string{
	tokEOF:         "EOF",
	tokBang:        "!",
	tokDollar:      "$",
	tokParenL:      "(",
	tokParenR:      ")",
	tokSpread:      "...",
	tokColon:       ":",
	tokEquals:      "=",
	tokAt:          "@",
	tokBracketL:    "[",
	tokBracketR:    "]",
	tokBraceL:      "{",
	tokPipe:        "|",
	tokBraceR:      "}",
	tokName:        "Name",
	tokInt:         "Int",
	tokFloat:       "Float",
	tokString:      "String",
	tokBlockString: "BlockString",
	tokAmp:         "&",
}

func (kind tokenKind) string() string {
	return tokenDescription[kind]
}

type lexer struct {
	source            string
	tokens            []token
	currentTokenIndex int
}

func newlexer(src string) (*lexer, error) {
	lexer := &lexer{
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

func (l *lexer) get() *token {
	tok := l.current()

	if tok.Value != "EOF" {
		l.currentTokenIndex++
	}

	return tok
}

func (l *lexer) put() {
	if l.currentTokenIndex > 0 {
		l.currentTokenIndex--
	}
}

func (l *lexer) putMany(c int) {
	if l.currentTokenIndex-c >= 0 {
		l.currentTokenIndex -= c
	}
}

func (l *lexer) prevLocation() *location {
	if l.currentTokenIndex == 0 {
		loc := &location{0, 0, l.source}

		return loc
	} else {
		return &location{l.tokens[l.currentTokenIndex-1].Start,
			l.tokens[l.currentTokenIndex-1].End,
			l.source}
	}
}

func (l *lexer) current() *token {
	if l.currentTokenIndex >= len(l.tokens) {
		return &token{tokEOF, len(l.source), len(l.source), "EOF"}
	}

	return &l.tokens[l.currentTokenIndex]
}

func (l *lexer) location() *location {
	tok := l.current()

	return &location{tok.Start, tok.End, l.source}
}

func (l *lexer) tokenEquals(tokVals ...string) bool {
	if l.currentTokenIndex+len(tokVals)-1 >= len(l.tokens) {
		return false
	}

	for i, val := range tokVals {
		if l.tokens[l.currentTokenIndex+i].Value != val {
			return false
		}
	}

	return true
}

func lex(doc string) ([]token, error) {
	var (
		whiteSpaceOn  bool
		stringOn      bool
		blockStringOn bool
		commentOn     bool
		tok           string
		kind          tokenKind
		tokens        []token
	)

	runes := []rune(doc)

	//for index, char := range doc {
	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		// This case handles all the punctuator characters that stand as an independent token.
		case '!', '$', '(', ')', ':', '=', '@', '[', ']', '{', '}', '|', '&':
			{
				if commentOn {
					continue
				}

				// Turn off the white space flag.
				whiteSpaceOn = false

				// If the token is not empty, it means that the current character
				// is a punctuator that ended a token.
				if tok != "" {
					// Append the token token to the tokens slice.
					tokens = append(tokens, token{kind, i - len(tok), i - 1, tok})

					// Empty the token.
					tok = ""
				}

				// Append a punctuator token to the tokens slice.
				tokens = append(tokens, token{getPunctuatorKind(runes[i]), i, i, string(runes[i])})
			}
		// This case handles the "spread" operator and float's decimal point.
		case '.':
			{
				if commentOn {
					continue
				}

				// If the kind of the token is INT, set it to FLOAT.
				if kind == tokInt {
					kind = tokFloat
				} else if kind == tokFloat {
					// If we met '.' when the token kind is FLOAT, return an error.
					return nil, errors.New("unexpected character '.' at position " + strconv.Itoa(i))
				}

				// Append the character to the token.
				tok = tok + string(runes[i])

				// Check if the token is a spread operator.
				if tok == "..." {
					// Append the token to the tokens slice.
					tokens = append(tokens, token{tokSpread, i, i, "..."})

					// Empty the token.
					tok = ""
				}
			}
		// This case is responsible for detection of STRING and BLOCK_STRING values.
		case '"':
			{
				if commentOn {
					continue
				}

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
						kind = tokBlockString

						// Append the '"' runes to the token.
						tok = tok + string(runes[i:i+3])

						// Increment the loop index twice.
						i += 2
					} else {
						// Turn the string flag on
						stringOn = true

						// Set the token king to string
						kind = tokString

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

					// If both of the flags are off (which means that we are inside any string value
					// anymore), append the token to the tokens slice.
					if !stringOn && !blockStringOn {
						// Append the token to the tokens slice.
						tokens = append(tokens, token{kind, i - len(tok), i - 1, tok})

						// Empty the token.
						tok = ""
					}
				}
			}
		// This case handles white spaces (single white space character only)
		case ' ', '\t':
			{
				if commentOn {
					continue
				}

				// If the double quote flag is on, we append any white space
				// we meet to the token variable.
				// Else, we will treat the white space as a delimiter between tokens
				if stringOn || blockStringOn {
					tok = tok + string(runes[i])
				} else {
					// If the white space flag is off, turn it on.
					if !whiteSpaceOn {
						whiteSpaceOn = true

						// If the token is not empty
						if tok != "" {
							// Append the token to tokens slice.
							tokens = append(tokens, token{kind, i - len(tok), i - 1, tok})

							// Empty the token
							tok = ""
						}
					}
				}
			}
		// This case handles commas.
		case ',':
			{
				if commentOn {
					continue
				}

				// Turn off the white space flag.
				whiteSpaceOn = false

				// If the token is not empty
				if tok != "" {
					// Append the token to the tokens slice.
					tokens = append(tokens, token{kind, i - len(tok), i - 1, tok})

					// Empty the token
					tok = ""
				}
			}
		// This case handles line feed and carriage return.
		case '\n', '\r':
			{
				if commentOn {
					commentOn = false
					continue
				}

				// If we are inside a string value, line feed and carriage return are disallowed.
				if stringOn {
					return nil, errors.New("line feed and carriage return characters are " +
						"disallowed in a string value.")
				}

				// If one if the string flags are on, append the character to the token.
				// Else, treat is as the end of the previous token.
				if blockStringOn {
					tok = tok + string(runes[i])
				} else {
					// If the token is not empty, it means that the current character ends a token.
					if tok != "" {
						// Append the token to the tokens slice.
						tokens = append(tokens, token{kind, i - len(tok), i - 1, tok})

						// Empty the token.
						tok = ""
					}
				}
			}
		case '_', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
			'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
			'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			{
				if commentOn {
					continue
				}

				// If the current character is the first character of the token,
				// set the token kind to NAME.
				if tok == "" {
					kind = tokName
				}

				// Append the character to the token variable.
				tok = tok + string(runes[i])
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			{
				if commentOn {
					continue
				}

				// If the current character (which is a digit) is the first character
				// of the token, set the kind as an INT.
				if tok == "" {
					kind = tokInt
				}

				// Append the character to the token variable.
				tok = tok + string(runes[i])
			}
		case '#':
			{
				commentOn = true
			}
		// This case handles the rest of the unicode characters (hebrew, chinese, etc)
		default:
			{
				tok = tok + string(runes[i])
			}
		}
	}

	// Append an EOF token to mark the end of the document.
	tokens = append(tokens, token{tokEOF, len(doc), len(doc), tokEOF.string()})

	// Return the tokens slice.
	return tokens, nil
}
