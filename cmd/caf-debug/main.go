package main

import (
	"fmt"

	"github.com/omeryahud/caf/internal/pkg/graphql/language"
)

func main() {
	query := `
	query inlineFragmentNoType($expandedInfo: Boolean) {
		user(handle: "zuck") {
		  id
		  name
		  ... @include(if: $expandedInfo) {
			firstName
			lastName
			birthday
		  }
		}
	  }
	`

	_, err := language.Parse(nil, query)

	if err != nil {
		fmt.Println(err)
	}
}
