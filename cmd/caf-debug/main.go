package main

import (
	"fmt"

	"github.com/omeryahud/caf/internal/pkg/graphql/language"
)

func main() {
	query := `
	{
<<<<<<< HEAD
		"""Hello"""
		user(id: 4) {
=======
		user(id 4) {
>>>>>>> b05ca5301ce48e96fd6360f15543291fa55e8ce7
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
