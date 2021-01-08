package graphqlvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/apidome/gateway/internal/pkg/graphql/language"
	"github.com/pkg/errors"
)

type GraphQLValidator struct {
	schema []byte
}

func NewGraphQLValidator() (GraphQLValidator, error) {
	return GraphQLValidator{}, nil
}

func (gv GraphQLValidator) LoadSchema(path, method string, bytes []byte) error {
	return nil
}

func (gv GraphQLValidator) Validate(path, method string, bytes []byte) error {
	// Create a struct that can contain an http body that contains a graphql document.
	var body struct {
		Query     string `json:"query"`
		Variables string `json:"variables"`
	}

	// Unmarshal the body into the struct.
	err := json.Unmarshal(bytes, &body)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal request body: ")
	}

	// Parse the query in order to get an AST object.
	ast, err := language.Parse(nil, body.Query)
	if err != nil {
		return errors.Wrap(err, "GraphQL parser failed: ")
	}

	fmt.Println("Raw Document:\n", body.Query)
	fmt.Println("\nParsed Document:\n", ast)

	return nil
}
