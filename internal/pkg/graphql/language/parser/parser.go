package parser

import (
	"regexp"

	"github.com/omeryahud/caf/internal/pkg/graphql/language/ast"
	"github.com/omeryahud/caf/internal/pkg/graphql/language/lexer"
	"github.com/pkg/errors"
)

var (
	errDoesntExist error = errors.New("Element does not exist")
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

	if opType != lexer.QUERY &&
		opType != lexer.MUTATION &&
		opType != lexer.SUBSCRIPTION {
		return nil, errDoesntExist
	} else {
		l.Get()

		opDefinition := &ast.OperationDefinition{}

		// optional
		name, err := parseName(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		opDefinition.Name = name

		// optional
		varDef, err := parseVariableDefinitions(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		opDefinition.VariableDefinitions = varDef

		// optional
		directives, err := parseDirectives(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
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
		return nil, errDoesntExist
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
			if err != errDoesntExist {
				return nil, err
			}
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

//
func parseVariableDefinitions(l *lexer.Lexer) (*ast.VariableDefinitions, error) {
	if l.Current().Value != lexer.PAREN_L.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		varDefs := &ast.VariableDefinitions{}

		for l.Current().Value != lexer.PAREN_R.String() {
			varDef, err := parseVariableDefinition(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			*varDefs = append(*varDefs, *varDef)
		}

		// Get closing parentheses
		l.Get()

		return varDefs, nil
	}
}

//
func parseVariableDefinition(l *lexer.Lexer) (*ast.VariableDefinition, error) {
	varDef := &ast.VariableDefinition{}

	_var, err := parseVariable(l)

	if err != nil {
		return nil, err
	}

	varDef.Variable = *_var

	defVal, err := parseDefaultValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	varDef.DefaultValue = defVal

	directives, err := parseDirectives(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	varDef.Directives = directives

	return varDef, nil
}

//
func parseDirectives(l *lexer.Lexer) (*ast.Directives, error) {
	dirs := &ast.Directives{}

	for {
		dir, err := parseDirective(l)

		if err != nil {
			if err == errDoesntExist {
				break
			}

			return nil, err
		}

		*dirs = append(*dirs, *dir)
	}
	return nil, nil
}

//
func parseDirective(l *lexer.Lexer) (*ast.Directive, error) {
	if l.Current().Value != lexer.AT.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		dir := &ast.Directive{}

		name, err := parseName(l)

		if err != nil {
			return nil, err
		}

		dir.Name = *name

		args, err := parseArguments(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		dir.Arguments = args

		return dir, nil
	}
}

//
func parseSelectionSet(l *lexer.Lexer) (*ast.SelectionSet, error) {
	if l.Current().Value != lexer.BRACE_L.String() {
		return nil, errDoesntExist
	} else {
		selSet := &ast.SelectionSet{}

		for l.Current().Value != lexer.BRACE_R.String() {
			sel, err := parseSelection(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			*selSet = append(*selSet, *sel)
		}

		if l.Current().Value != lexer.BRACE_R.String() {
			return nil, errors.New("Expecting closing bracket for selection set")
		}

		return selSet, nil
	}
}

//
func parseSelection(l *lexer.Lexer) (*ast.Selection, error) {
	var sel ast.Selection

	sel, err := parseField(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	sel, err = parseFragmentSpread(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	sel, err = parseInlineFragment(l)

	if err != nil {
		return nil, err
	}

	return &sel, nil
}

//
func parseVariable(l *lexer.Lexer) (*ast.Variable, error) {
	if l.Current().Value != lexer.DOLLAR.String() {
		return nil, errDoesntExist
	} else {
		_var := &ast.Variable{}

		name, err := parseName(l)

		if err != nil {
			return nil, err
		}

		_var.Name = *name

		return _var, nil
	}
}

//
func parseDefaultValue(l *lexer.Lexer) (*ast.DefaultValue, error) {
	if l.Current().Value != lexer.EQUALS.String() {
		return nil, errDoesntExist
	} else {
		dVal := &ast.DefaultValue{}

		val, err := parseValue(l)

		if err != nil {
			return nil, err
		}

		dVal.Value = *val

		return dVal, nil
	}
}

//
func parseValue(l *lexer.Lexer) (*ast.Value, error) {
	_var, err := parseVariable(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	} else {
		// ! need to read variable value if it exists
		_var = _var
	}

	var val ast.Value

	val, err = parseIntValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	val, err = parseFloatValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	val, err = parseStringValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	val, err = parseBooleanValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	val, err = parseNullValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	val, err = parseEnumValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	val, err = parseListValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	val, err = parseObjectValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	if val == nil {
		return nil, errDoesntExist
	} else {
		return &val, nil
	}
}

func parseArguments(l *lexer.Lexer) (*ast.Arguments, error) {
	if l.Current().Value != lexer.PAREN_L.String() {

	}
}

//
func parseField(l *lexer.Lexer) (*ast.Field, error) {
	field := &ast.Field{}

	alias, err := parseAlias(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	field.Alias = alias

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	field.Name = *name

	args, err := parseArguments(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	field.Arguments = args

	dirs, err := parseDirectives(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	field.Directives = dirs

	selSet, err := parseSelectionSet(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	field.SelectionSet = selSet

	return field, nil
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

func parseIntValue(l *lexer.Lexer) (*ast.IntValue, error) {
	return nil, nil
}

func parseFloatValue(l *lexer.Lexer) (*ast.FloatValue, error) {
	return nil, nil
}

func parseStringValue(l *lexer.Lexer) (*ast.StringValue, error) {
	return nil, nil
}

func parseBooleanValue(l *lexer.Lexer) (*ast.BooleanValue, error) {
	return nil, nil
}

func parseNullValue(l *lexer.Lexer) (*ast.NullValue, error) {
	return nil, nil
}

func parseEnumValue(l *lexer.Lexer) (*ast.EnumValue, error) {
	return nil, nil
}

func parseListValue(l *lexer.Lexer) (*ast.ListValue, error) {
	return nil, nil
}

func parseObjectValue(l *lexer.Lexer) (*ast.ObjectValue, error) {
	return nil, nil
}
