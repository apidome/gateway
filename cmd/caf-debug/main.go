package main

import (
	"fmt"

	"github.com/omeryahud/caf/internal/pkg/graphql/language"
)

func main() {
	query := `
	query inlineFragmentTyping {
		profiles(handles: ["zuck", "cocacola"]) {
		  handle
		  ... on User {
			friends {
			  count
			}
		  }
		  ... on Page {
			likers {
			  count
			}
		  }
		}
	  }
	`

	_, err := language.Parse(query)

	if err != nil {
		fmt.Println(err)
	}
}
