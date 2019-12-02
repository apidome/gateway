package parser

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/ast"
	"github.com/omeryahud/caf/internal/pkg/graphql/language/lexer"
	"github.com/pkg/errors"
)

func Parse(doc string) (*ast.Document, error) {
	lexer, err := lexer.NewLexer(doc)
	if err != nil {
		return nil, err
	}

	return parseDocument(lexer)
}

func parseDocument(lexer *lexer.Lexer) (*ast.Document, error) {
	var op ast.ExecutableDefinition
	var err error
	document := &ast.Document{}

	for _, token := range lexer.Tokens {
		if token.Kind == lexer.BRACE_L {
			op, err = parseQuery(lexer)
			if err != nil {
				return nil, err
			}
		}

		switch token.Value {
		case lexer.QUERY:
			{
				op, err = parseQuery(lexer)
				if err != nil {
					return nil, err
				}
			}
		case lexer.MUTATION:
			{
				op, err = parseMutation(lexer)
				if err != nil {
					return nil, err
				}
			}
		case lexer.SUBSCRIPTION:
			{
				op, err = parseSubscription(lexer)
				if err != nil {
					return nil, err
				}
			}

		case lexer.FRAGMENT:
			{
				op, err = parseFragment(lexer)
				if err != nil {
					return nil, err
				}
			}
		default:
			{
				return nil, errors.New("invalid operation " + token.Value)
			}
		}

		document.Definitions = append(document.Definitions, op)
	}

	return document, nil
}

func parseQuery(lexer *lexer.Lexer) (*ast.OperationDefinition, error) {
	return nil, nil
}

func parseMutation(lexer *lexer.Lexer) (*ast.OperationDefinition, error) {
	return nil, nil
}

func parseSubscription(lexer *lexer.Lexer) (*ast.OperationDefinition, error) {
	return nil, nil
}

func parseFragment(lexer *lexer.Lexer) (*ast.FragmentDefinition, error) {
	return nil, nil
}

func parseOperationDefinition(lexer *lexer.Lexer) (*ast.OperationDefinition, error) {
	return nil, nil
}

func parseFragmentDefinition(lexer *lexer.Lexer) (*ast.FragmentDefinition, error) {
	return nil, nil
}

func parseName(n string) (string, error) {
	return "", nil
}

func parseVariableDefinition(lexer *lexer.Lexer) (*ast.VariableDefinitions, error) {
	return nil, nil
}

func parseDirectives(lexer *lexer.Lexer) (*ast.Directives, error) {
	return nil, nil
}

func parseSelectionSet(lexer *lexer.Lexer) (*ast.SelectionSet, error) {
	return nil, nil
}
