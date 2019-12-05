package parser

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/ast"
	"github.com/omeryahud/caf/internal/pkg/graphql/language/lexer"
	"github.com/pkg/errors"
)

func Parse(doc string) (*ast.Document, error) {
	l, err := lexer.NewLexer(doc)
	if err != nil {
		return nil, err
	}

	return parseDocument(l)
}

func parseDocument(l *lexer.Lexer) (*ast.Document, error) {
	var op ast.ExecutableDefinition
	var err error
	document := &ast.Document{}

	for token := l.Current(); token.Kind != lexer.EOF; {
		if token.Kind == lexer.BRACE_L {
			op, err = parseOperationDefinition(l, ast.OPERATION_QUERY)
			if err != nil {
				return nil, err
			}
		} else {
			switch token.Value {
			case lexer.QUERY:
				{
					op, err = parseOperationDefinition(l, ast.OPERATION_QUERY)
					if err != nil {
						return nil, err
					}
				}
			case lexer.MUTATION:
				{
					op, err = parseOperationDefinition(l, ast.OPERATION_MUTATION)
					if err != nil {
						return nil, err
					}
				}
			case lexer.SUBSCRIPTION:
				{
					op, err = parseOperationDefinition(l, ast.OPERATION_SUBSCRIPTION)
					if err != nil {
						return nil, err
					}
				}

			case lexer.FRAGMENT:
				{
					op, err = parseFragment(l)
					if err != nil {
						return nil, err
					}
				}
			default:
				{
					return nil, errors.New("invalid operation - " + token.Value)
				}
			}
		}

		document.Definitions = append(document.Definitions, op)
	}

	return document, nil
}

func parseFragment(l *lexer.Lexer) (*ast.FragmentDefinition, error) {
	return nil, nil
}

func parseOperationDefinition(l *lexer.Lexer, operationType ast.OperationType) (*ast.OperationDefinition, error) {
	return nil, nil
}

func parseFragmentDefinition(l *lexer.Lexer) (*ast.FragmentDefinition, error) {
	return nil, nil
}

func parseName(n string) (string, error) {
	return "", nil
}

func parseVariableDefinition(l *lexer.Lexer) (*ast.VariableDefinitions, error) {
	return nil, nil
}

func parseDirectives(l *lexer.Lexer) (*ast.Directives, error) {
	return nil, nil
}

func parseSelectionSet(l *lexer.Lexer) (*ast.SelectionSet, error) {
	return nil, nil
}
