package language

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

func Lex(doc string) ([]Token, error) {
	for index := range doc {
		switch doc[index] {

		}
	}

	return nil, nil
}
