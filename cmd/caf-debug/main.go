package main

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/parser"
)

func main() {
	query := `
	{
		#commet
		empireHero: hero(episode: EMPIRE) {
		  name
		}
		jediHero: hero(episode: JEDI) {
		  name
		}
	}
	`

	parser.Parse(query)
}
