package graphqlvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/omeryahud/caf/internal/pkg/graphql/language/parser"
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
		Query     json.RawMessage `json:"query"`
		Variables json.RawMessage `json:"variables"`
	}

	// Unmarshal the body into the struct.
	err := json.Unmarshal(bytes, &body)
	if err != nil {
		return err
	}

	// Remove the quote characters from the query string.
	query := string(body.Query)[1 : len(string(body.Query))-1]

	// Parse the query in order to get an AST object.
	ast, err := parser.Parse(query)
	if err != nil {
		return errors.Wrap(err, "graphQL parser failed: ")
	}

	fmt.Println(query)
	fmt.Println(ast)

	return nil
}
