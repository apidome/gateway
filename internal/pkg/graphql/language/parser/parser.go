package parser

import (
	"regexp"

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

// idk what this is
func parseFragment(l *lexer.Lexer) (*ast.FragmentDefinition, error) {
	return nil, nil
}

//
func parseOperationDefinition(l *lexer.Lexer, operationType ast.OperationType) (*ast.OperationDefinition, error) {
	opType := l.Current().Value

	if opType != lexer.QUERY ||
		opType != lexer.MUTATION ||
		opType != lexer.SUBSCRIPTION {
		return nil, errors.New("No operation type found")
	} else {
		opDefinition := &ast.OperationDefinition{}

		// optional
		name, err := parseName(l)

		if err != nil {
			return nil, err
		}

		opDefinition.Name = name

		// optional
		varDef, err := parseVariableDefinitions(l)

		if err != nil {
			return nil, err
		}

		opDefinition.VariableDefinitions = varDef

		// optional
		directives, err := parseDirectives(l)

		if err != nil {
			return nil, err
		}

		opDefinition.Directives = directives

		selSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		opDefinition.SelectionSet = *selSet

		return opDefinition, nil
	}
}

//
func parseFragmentDefinition(l *lexer.Lexer) (*ast.FragmentDefinition, error) {
	if l.Current().Value != lexer.FRAGMENT {
		return nil, errors.New("Token is not a fragment definition")
	} else {
		l.Get()

		fragDef := &ast.FragmentDefinition{}

		name, err := parseFragmentName(l)

		if err != nil {
			return nil, err
		}

		if name.Value == "on" {
			return nil, errors.New("Fragment name cannot be 'on'")
		}

		fragDef.FragmentName = *name

		typeCond, err := parseTypeCondition(l)

		if err != nil {
			return nil, err
		}

		fragDef.TypeCondition = *typeCond

		// optional
		directives, err := parseDirectives(l)

		if err != nil {
			return nil, err
		}

		fragDef.Directives = directives

		selectionSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		fragDef.SelectionSet = *selectionSet

		return fragDef, nil
	}
}

//
func parseName(l *lexer.Lexer) (*ast.Name, error) {
	name := &ast.Name{}
	token := l.Current()
	pattern := "^[_A-Za-z][_0-9A-Za-z]*$"

	// If the current token is not a Name, return nil
	if token.Kind != lexer.NAME {
		return nil, errors.New("Not a name")
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

	l.Get()

	// Populate the Name struct.
	name.Value = token.Value
	name.Loc.Start = token.Start
	name.Loc.End = token.End
	name.Loc.Source = l.Source()

	// Return the AST Name object.
	return name, nil
}

func parseVariableDefinitions(l *lexer.Lexer) (*ast.VariableDefinitions, error) {
	if l.Current().Value != lexer.PAREN_L.String() {
		return nil, errors.New("Expecting '(' to parse variable definitions")
	} else {
		l.Get()
	}

	return nil, nil
}

func parseVariableDefinition(l *lexer.Lexer) (*ast.VariableDefinition, error) {
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
