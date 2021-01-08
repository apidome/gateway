package main

import (
	"fmt"

	"github.com/apidome/gateway/internal/pkg/graphql/language"
)

func main() {
	query := `
	type Query {
		dog: Dog
		cat: Cat
	  }
	  
	  enum DogCommand { SIT, DOWN, HEEL }
	  
	  interface Sentient {
		name: String!
	  }
	  
	  interface Pet {
		name: String!
	  }
	  
	  type Alien implements Sentient {
		name: String!
		homePlanet: String
	  }
	  
	  type Human implements Sentient {
		name: String!
		pets: [Pet!]
	  }
	  
	  enum CatCommand { JUMP }
	  
	  type Cat implements Pet {
		name: String!
		nickname: String
		doesKnowCommand(catCommand: CatCommand!): Boolean!
		meowVolume: Int
	  }
	  
	  union CatOrDog = Cat | Dog
	  union DogOrHuman = Dog | Human
	  union HumanOrAlien = Human | Alien
	  
	  input ComplexInput { name: String, owner: String }
	  
	  extend type Query {
		findDog(complex: ComplexInput): Dog
		booleanList(booleanListArg: [Boolean!]): Boolean
	  }
	  
	  type Dog implements Pet {
		name: String!
		nickname: String
		barkVolume: Int
		doesKnowCommand(dogCommand: DogCommand!): Boolean!
		isHousetrained(atOtherHomes: [Boolean]!): Boolean!
		owner: Human
	  }
`

	_, err := language.Parse(nil, query)

	if err != nil {
		fmt.Println(err)
	}
}
