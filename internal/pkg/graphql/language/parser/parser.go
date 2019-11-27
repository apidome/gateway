package parser

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/ast"
	"github.com/omeryahud/caf/internal/pkg/graphql/language/lexer"
)

func Parse(doc string) (*ast.Document, error) {
	tokens, err := lexer.Lex(doc)
	if err != nil {
		return nil, err
	}

	return parseDocument(tokens)
}

func parseDocument(tokens []lexer.Token) (*ast.Document, error) {
	return nil, nil
}

func parseOperationDefinition(tokens []lexer.Token) (*ast.OperationDefinition, error) {
	return nil, nil
}

func parseFragmentDefinition(tokens []lexer.Token) (*ast.FragmentDefinition, error) {
	return nil, nil
}
