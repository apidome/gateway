package language

import (
	"testing"
)

func TestExtractUsedFragmentsNames(t *testing.T) {
	rawSchema := `
type Query {
  dog: Dog
  findDog(complex: ComplexInput): Dog
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
  booleanList(booleanListArg: [Boolean!]): Boolean
}

type Dog implements Pet {
  name: String!
  nickname: String
  barkVolume: Int
  doesKnowCommand(dogCommand: DogCommand!): Boolean!
  isHousetrained(atOtherHomes: Boolean): Boolean!
  owner: Human
}
`

	query := `
fragment goodBooleanArg on Arguments {
  booleanArgField(booleanArg: true)
}

fragment coercedIntIntoFloatArg on Arguments {
  # Note: The input coercion rules for Float allow Int literals.
  floatArgField(floatArg: 123)
}

query goodComplexDefaultValue($search: ComplexInput = { name: "Fido" }) {
  findDog(complex: $search)
}
`

	schemaAST, err := ParseSchema(rawSchema)
	if err != nil {
		t.Error("schema parse failed: ", err)
	}

	queryAST, err := Parse(schemaAST, query)
	if err != nil {
		t.Fatal("query parse failed: ", err)
	}

	validateValuesOfCorrectType(schemaAST, queryAST)

	t.Log("Validation Succeeded")
}
