package parser

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/ast"
	"github.com/omeryahud/caf/internal/pkg/graphql/language/lexer"
	"github.com/pkg/errors"
	"regexp"
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

	for token := l.Current(); token.Kind != lexer.EOF; token = l.Get() {
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

func parseName(l *lexer.Lexer) (*ast.Name, error) {
	name := new(ast.Name)
	token := l.Current()
	pattern := "^[_A-Za-z][_0-9A-Za-z]*$"

	// If the current token is not a Name, return nil
	if token.Kind != lexer.NAME {
		return nil, nil
	}

	// Check if the given name matches the regex provided by graphql spec at
	// https://graphql.github.io/graphql-spec/draft/#Name
	match, err := regexp.MatchString(pattern, token.Value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse name: ")
	}

	// If the name does not match the requirements, return an error.
	if !match {
		return nil, errors.New("invalid name - " + token.Value)
	}

	// Populate the Name struct.
	name.Value = token.Value
	name.Loc.Start = token.Start
	name.Loc.End = token.End
	name.Loc.Source = l.Source()

	// Return the AST Name object.
	return name, nil
}

func parseVariableDefinitions(l *lexer.Lexer) (*ast.VariableDefinitions, error) {
	return nil, nil
}

func parseDirectives(l *lexer.Lexer) (*ast.Directives, error) {
	return nil, nil
}

func parseSelectionSet(l *lexer.Lexer) (*ast.SelectionSet, error) {
	return nil, nil
}

func parseVariable(l *lexer.Lexer) (*ast.Variable, error) {
	return nil, nil
}

func parseDefaultValue(l *lexer.Lexer) (*ast.DefaultValue, error) {
	return nil, nil
}

func parseValue(lexer2 *lexer.Lexer) (*ast.Value, error) {
	return nil, nil
}

func parseArguments(l *lexer.Lexer) (*ast.Arguments, error) {
	return nil, nil
}

func parseField(l *lexer.Lexer) (*ast.Field, error) {
	return nil, nil
}

func parseFragmentSpread(l *lexer.Lexer) (*ast.FragmentSpread, error) {
	return nil, nil
}

func parseInlineFragment(l *lexer.Lexer) (*ast.InlineFragment, error) {
	return nil, nil
}

func parseAlias(l *lexer.Lexer) (*ast.Alias, error) {
	return nil, nil
}

func parseFragmentName(l *lexer.Lexer) (*ast.FragmentName, error) {
	return nil, nil
}

func parseTypeCondition(l *lexer.Lexer) (*ast.TypeCondition, error) {
	return nil, nil
}

func parseNamedType(l *lexer.Lexer) (*ast.NamedType, error) {
	return nil, nil
}
