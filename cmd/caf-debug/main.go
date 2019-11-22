package main

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language"
)

func main() {
	language.Lex(`"""Hello There!"""`)
}
