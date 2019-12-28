package main

import (
	"fmt"

	"github.com/omeryahud/caf/internal/pkg/graphql/language"
)

func main() {
	query := `
	{
		"""Hello"""
		user(id: 4) {
		  id
		  name
		  profilePic(width: 100, height: 50)
		}
	  }
	`

	doc, err := language.Parse(query)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(doc)
}
