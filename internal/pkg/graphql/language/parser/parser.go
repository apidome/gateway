package parser

import (
	"regexp"
	"strconv"

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

		// optional
		name, err := parseName(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		// optional
		varDef, err := parseVariableDefinitions(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		// optional
		directives, err := parseDirectives(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		selSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		opDefinition := &ast.OperationDefinition{}

		opDefinition.Name = name
		opDefinition.VariableDefinitions = varDef
		opDefinition.Directives = directives
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

		name, err := parseFragmentName(l)

		if err != nil {
			return nil, err
		}

		if name.Value == "on" {
			return nil, errors.New("Fragment name cannot be 'on'")
		}

		typeCond, err := parseTypeCondition(l)

		if err != nil {
			return nil, err
		}

		// optional
		directives, err := parseDirectives(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		selectionSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		fragDef := &ast.FragmentDefinition{}

		fragDef.FragmentName = *name
		fragDef.TypeCondition = *typeCond
		fragDef.Directives = directives
		fragDef.SelectionSet = *selectionSet

		return fragDef, nil
	}
}

//
func parseName(l *lexer.Lexer) (*ast.Name, error) {
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

	name := &ast.Name{}

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
	_var, err := parseVariable(l)

	if err != nil {
		return nil, err
	}

	defVal, err := parseDefaultValue(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	directives, err := parseDirectives(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	varDef := &ast.VariableDefinition{}

	varDef.Variable = *_var
	varDef.DefaultValue = defVal
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

		name, err := parseName(l)

		if err != nil {
			return nil, err
		}

		args, err := parseArguments(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		dir := &ast.Directive{}

		dir.Name = *name
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
		name, err := parseName(l)

		if err != nil {
			return nil, err
		}

		_var := &ast.Variable{}

		_var.Name = *name

		return _var, nil
	}
}

//
func parseDefaultValue(l *lexer.Lexer) (*ast.DefaultValue, error) {
	if l.Current().Value != lexer.EQUALS.String() {
		return nil, errDoesntExist
	} else {
		val, err := parseValue(l)

		if err != nil {
			return nil, err
		}

		dVal := &ast.DefaultValue{}

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

//
func parseArguments(l *lexer.Lexer) (*ast.Arguments, error) {
	if l.Current().Value != lexer.PAREN_L.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		args := &ast.Arguments{}

		for l.Current().Value != lexer.PAREN_R.String() {
			arg, err := parseArgument(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			*args = append(*args, *arg)
		}

		if l.Current().Value != lexer.PAREN_R.String() {
			return nil, errors.New("Expecting closing parentheses for arguments")
		}

		return args, nil
	}
}

//
func parseArgument(l *lexer.Lexer) (*ast.Argument, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if l.Current().Value != lexer.COLON.String() {
		return nil, errors.New("Expecting colon after argument name")
	}

	val, err := parseValue(l)

	if err != nil {
		return nil, err
	}

	arg := &ast.Argument{}

	arg.Name = *name
	arg.Value = *val

	return arg, nil
}

//
func parseField(l *lexer.Lexer) (*ast.Field, error) {
	alias, err := parseAlias(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	args, err := parseArguments(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	dirs, err := parseDirectives(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	selSet, err := parseSelectionSet(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	field := &ast.Field{}

	field.Alias = alias
	field.Name = *name
	field.Arguments = args
	field.Directives = dirs
	field.SelectionSet = selSet

	return field, nil
}

//
func parseFragmentSpread(l *lexer.Lexer) (*ast.FragmentSpread, error) {
	if l.Current().Value != lexer.SPREAD.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		fname, err := parseFragmentName(l)

		if err != nil {
			return nil, err
		}

		directives, err := parseDirectives(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		spread := &ast.FragmentSpread{}

		spread.FragmentName = *fname
		spread.Directives = directives

		return spread, nil
	}
}

func parseInlineFragment(l *lexer.Lexer) (*ast.InlineFragment, error) {
	if l.Current().Value != lexer.SPREAD.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		typeCon, err := parseTypeCondition(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		directives, err := parseDirectives(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		}

		selSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		inlineFrag := &ast.InlineFragment{}

		inlineFrag.TypeCondition = typeCon
		inlineFrag.Directives = directives
		inlineFrag.SelectionSet = *selSet

		return inlineFrag, nil
	}
}

//
func parseAlias(l *lexer.Lexer) (*ast.Alias, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if l.Current().Value != lexer.COLON.String() {
		return nil, errors.New("Expecting colon after alias name")
	} else {
		var alias *ast.Alias

		*alias = ast.Alias(*name)

		return alias, nil
	}
}

//
func parseFragmentName(l *lexer.Lexer) (*ast.FragmentName, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if name.Value == "on" {
		return nil, errors.New("Fragment name cannot be 'on'")
	}

	var fragName *ast.FragmentName

	*fragName = ast.FragmentName(*name)

	return fragName, nil
}

//
func parseTypeCondition(l *lexer.Lexer) (*ast.TypeCondition, error) {
	if l.Current().Value != "on" {
		return nil, errDoesntExist
	} else {
		namedType, err := parseNamedType(l)

		if err != nil {
			return nil, err
		}

		typeCond := &ast.TypeCondition{}

		typeCond.NamedType = *namedType

		return typeCond, nil
	}
}

//
func parseNamedType(l *lexer.Lexer) (*ast.NamedType, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	var namedType *ast.NamedType

	*namedType = ast.NamedType(*name)

	return namedType, nil
}

//
func parseIntValue(l *lexer.Lexer) (*ast.IntValue, error) {
	intVal, err := strconv.ParseInt(l.Current().Value, 10, 64)

	if err != nil {
		return nil, err
	}

	l.Get()

	intValP := &ast.IntValue{}

	intValP.Value = intVal

	return intValP, nil
}

//
func parseFloatValue(l *lexer.Lexer) (*ast.FloatValue, error) {
	floatVal, err := strconv.ParseFloat(l.Current().Value, 64)

	if err != nil {
		return nil, err
	}

	l.Get()

	floatValP := &ast.FloatValue{}

	floatValP.Value = floatVal

	return floatValP, nil
}

// ! Have a discussion about this function
func parseStringValue(l *lexer.Lexer) (*ast.StringValue, error) {
	return nil, nil
}

//
func parseBooleanValue(l *lexer.Lexer) (*ast.BooleanValue, error) {
	boolVal, err := strconv.ParseBool(l.Current().Value)

	if err != nil {
		return nil, err
	}

	l.Get()

	boolValP := &ast.BooleanValue{}

	boolValP.Value = boolVal

	return boolValP, nil
}

// ! Figure out what to do with a null value
func parseNullValue(l *lexer.Lexer) (*ast.NullValue, error) {
	if l.Current().Value != "null" {
		return nil, errDoesntExist
	} else {
		l.Get()

		null := &ast.NullValue{}

		return null, nil
	}
}

//
func parseEnumValue(l *lexer.Lexer) (*ast.EnumValue, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	switch name.Value {
	case "true", "false", "null":
		return nil, errors.New("Enum value cannot be 'true', 'false' or 'null'")
	default:
		enumVal := &ast.EnumValue{}

		enumVal.Name = *name

		return enumVal, nil
	}
}

//
func parseListValue(l *lexer.Lexer) (*ast.ListValue, error) {
	if l.Current().Value != "[" {
		return nil, errDoesntExist
	} else {
		l.Get()

		lstVal := &ast.ListValue{}

		for l.Current().Value != "]" {
			val, err := parseValue(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			lstVal.Values = append(lstVal.Values, *val)
		}

		if l.Current().Value != "]" {
			return nil, errors.New("Missing closing bracket for list value")
		}

		return lstVal, nil
	}
}

//
func parseObjectValue(l *lexer.Lexer) (*ast.ObjectValue, error) {
	if l.Current().Value != "{" {
		return nil, errDoesntExist
	} else {
		l.Get()

		objVal := &ast.ObjectValue{}

		for l.Current().Value != "}" {
			objField, err := parseObjectField(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			objVal.Values = append(objVal.Values, *objField)
		}

		if l.Current().Value != "}" {
			return nil, errors.New("Expecting a closing curly brace for an object value")
		}

		return objVal, nil
	}
}

//
func parseObjectField(l *lexer.Lexer) (*ast.ObjectField, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if l.Current().Value != ":" {
		return nil, errors.New("Expecting color after object field name")
	}

	l.Get()

	val, err := parseValue(l)

	if err != nil {
		return nil, err
	}

	objField := &ast.ObjectField{}

	objField.Name = *name
	objField.Value = *val

	return objField, nil
}
