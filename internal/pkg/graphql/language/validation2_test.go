package language

import (
	"testing"
)

func TestExtractUsedFragmentsNames(t *testing.T) {
	rawSchema := `
type Query {
  dog: Dog
}

enum DogCommand { SIT, DOWN, HEEL }

type Dog implements Pet {
  name: String!
  nickname: String
  barkVolume: Int
  doesKnowCommand(dogCommand: DogCommand!): Boolean!
  isHousetrained(atOtherHomes: Boolean): Boolean!
  owner: Human
}

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
`

	query := `
query A($atOtherHomes: Boolean) {
  ...HouseTrainedFragment
}

query B($atOtherHomes: Boolean) {
  ...HouseTrainedFragment
}

fragment HouseTrainedFragment on Dog {
  dog {
    isHousetrained(atOtherHomes: $atOtherHomes)
  }
}
`

	schemaAST, err := ParseSchema(rawSchema)
	if err != nil {
		t.Error(err)
	}

	queryAST, err := Parse(schemaAST, query)
	if err != nil {
		t.Fatal(err)
	}

	validateVariableUniqueness(*queryAST)

	t.Log("Validation Succeeded")
}
