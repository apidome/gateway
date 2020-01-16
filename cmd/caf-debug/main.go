package main

import (
	"fmt"

	"github.com/omeryahud/caf/internal/pkg/graphql/language"
)

func main() {
	query := `
	query getDogName {
		dog {
		  name
		  color
		}
	  }
	  
	  extend type Dog {
		color: String
	  }
	  
	`

	_, err := language.Parse(query)

	if err != nil {
		fmt.Println(err)
	}
}
