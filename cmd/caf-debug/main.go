package main

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/parser"
)

func main() {
	query := `
	mutation {
		likeStory(storyID: 12345) {
		  story {
			likeCount
		  }
		}
	  }	   
	`

	parser.Parse(query)
}
