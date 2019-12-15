package main

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/parser"
)

func main() {
	query := `
	{
		user(id: 4) {
		  id
		  name
		  profilePic(width: 100, height: 50)
		}
	  }
	  
	`

	parser.Parse(query)
}
