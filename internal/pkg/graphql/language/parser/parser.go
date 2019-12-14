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

	astDoc, _ := parseDocument(l)

	return astDoc, nil
}

//
func parseDocument(l *lexer.Lexer) (*ast.Document, error) {
	def, err := parseDefinitions(l)

	if err != nil {
		return nil, err
	}

	doc := &ast.Document{}

	doc.Definitions = *def

	return doc, nil
}

//
func parseDefinitions(l *lexer.Lexer) (*ast.Definitions, error) {
	defs := &ast.Definitions{}

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	for tok.Value != lexer.EOF.String() {
		var def ast.Definition

		def, err := parseDefinition(l)

		if err != nil {
			if err == errDoesntExist {
				break
			}
		}

		if def != nil {
			*defs = append(*defs, def)
		}

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}
	}

	if len(*defs) == 0 {
		return nil, errors.New("No definitions found in document")
	}

	return defs, nil
}

//
func parseDefinition(l *lexer.Lexer) (ast.Definition, error) {
	def, err := parseExecutableDefinition(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	if def != nil {
		return def, nil
	}

	def, err = parseTypeSystemDefinition(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	if def != nil {
		return def, nil
	}

	def, err = parseTypeSystemExtension(l)

	if err != nil {
		return nil, err
	}

	return def, nil
}

//
func parseExecutableDefinition(l *lexer.Lexer) (ast.ExecutableDefinition, error) {
	var execDef ast.ExecutableDefinition

	execDef, err := parseOperationDefinition(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	if execDef != nil {
		return execDef, nil
	}

	execDef, err = parseFragmentDefinition(l)

	if err != nil {
		return nil, err
	}

	return execDef, nil
}

//
func parseTypeSystemDefinition(l *lexer.Lexer) (ast.TypeSystemDefinition, error) {
	var def ast.TypeSystemDefinition

	def, err := parseSchemaDefinition(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	if def != nil {
		return def, nil
	}

	def, err = parseTypeDefinition(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	if def != nil {
		return def, nil
	}

	def, err = parseDirectiveDefinition(l)

	if err != nil {
		return nil, err
	}

	return def, nil
}

// ! Type system
func parseSchemaDefinition(l *lexer.Lexer) (*ast.SchemaDefinition, error) {
	return nil, nil
}

// ! Type system
func parseTypeDefinition(l *lexer.Lexer) (ast.TypeDefinition, error) {
	return nil, nil
}

// ! Type system
func parseDirectiveDefinition(l *lexer.Lexer) (*ast.DirectiveDefinition, error) {
	return nil, nil
}

// ! Type system
func parseTypeSystemExtension(l *lexer.Lexer) (ast.TypeSystemExtension, error) {
	return nil, nil
}

// ! idk what this is
func parseFragment(l *lexer.Lexer) (*ast.FragmentDefinition, error) {
	return nil, nil
}

// ! come back to this
func parseOperationDefinition(l *lexer.Lexer) (*ast.OperationDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value == lexer.BRACE_L.String() {
		// this is a query
		shorthandQuery, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		opDef := &ast.OperationDefinition{}

		opDef.SelectionSet = *shorthandQuery

		return opDef, nil
	} else if tok.Value != lexer.QUERY &&
		tok.Value != lexer.MUTATION &&
		tok.Value != lexer.SUBSCRIPTION {
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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.FRAGMENT {
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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	pattern := "^[_A-Za-z][_0-9A-Za-z]*$"

	// If the current token is not a Name, return nil
	if tok.Kind != lexer.NAME {
		return nil, errors.New("Not a name")
	}

	// Check if the given name matches the regex provided by graphql spec at
	// https://graphql.github.io/graphql-spec/draft/#Name
	match, err := regexp.MatchString(pattern, tok.Value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse name: ")
	}

	// If the name does not match the requirements, return an error.
	if !match {
		return nil, errors.New("invalid name - " + tok.Value)
	}

	l.Get()

	name := &ast.Name{}

	// Populate the Name struct.
	name.Value = tok.Value
	name.Loc.Start = tok.Start
	name.Loc.End = tok.End
	name.Loc.Source = l.Source()

	// Return the AST Name object.
	return name, nil
}

//
func parseVariableDefinitions(l *lexer.Lexer) (*ast.VariableDefinitions, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.PAREN_L.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		varDefs := &ast.VariableDefinitions{}

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}

		for tok.Value != lexer.PAREN_R.String() {
			varDef, err := parseVariableDefinition(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			*varDefs = append(*varDefs, *varDef)

			tok, err = l.Current()

			if err != nil {
				return nil, err
			}
		}

		// Get closing parentheses
		tok, err := l.Current()

		if err != nil {
			return nil, err
		}

		if tok.Value != lexer.PAREN_R.String() {
			return nil, errors.New("Expecting closing parentheses for variabel definitions")
		}

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

	if len(*dirs) == 0 {
		return nil, errDoesntExist
	}

	return dirs, nil
}

//
func parseDirective(l *lexer.Lexer) (*ast.Directive, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.AT.String() {
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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.BRACE_L.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		selSet := &ast.SelectionSet{}

		tok, err := l.Current()

		if err != nil {
			return nil, err
		}

		for tok.Value != lexer.BRACE_R.String() {
			sel, err := parseSelection(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			*selSet = append(*selSet, sel)

			tok, err = l.Current()

			if err != nil {
				return nil, err
			}
		}

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}

		if tok.Value != lexer.BRACE_R.String() {
			return nil, errors.New("Expecting closing bracket for selection set")
		}

		l.Get()

		return selSet, nil
	}
}

//
func parseSelection(l *lexer.Lexer) (ast.Selection, error) {
	var sel ast.Selection

	sel, err := parseField(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	if sel != nil {
		return sel, nil
	}

	sel, err = parseFragmentSpread(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	}

	if sel != nil {
		return sel, nil
	}

	sel, err = parseInlineFragment(l)

	if err != nil {
		return nil, err
	}

	return sel, nil
}

//
func parseVariable(l *lexer.Lexer) (*ast.Variable, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.DOLLAR.String() {
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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.EQUALS.String() {
		return nil, errDoesntExist
	} else {
		val, err := parseValue(l)

		if err != nil {
			return nil, err
		}

		dVal := &ast.DefaultValue{}

		dVal.Value = val

		return dVal, nil
	}
}

//
func parseValue(l *lexer.Lexer) (ast.Value, error) {
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

	if err == nil {
		return val, nil
	}

	val, err = parseFloatValue(l)

	if err == nil {
		return val, nil
	}

	val, err = parseStringValue(l)

	if err == nil {
		return val, nil
	}

	val, err = parseBooleanValue(l)

	if err == nil {
		return val, nil
	}

	val, err = parseNullValue(l)

	if err == nil {
		return val, nil
	}

	val, err = parseEnumValue(l)

	if err == nil {
		return val, nil
	}

	val, err = parseListValue(l)

	if err == nil {
		return val, nil
	}

	val, err = parseObjectValue(l)

	if err == nil {
		return val, nil
	}

	if err != nil {
		return nil, errDoesntExist
	} else {
		return val, nil
	}
}

//
func parseArguments(l *lexer.Lexer) (*ast.Arguments, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.PAREN_L.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		args := &ast.Arguments{}

		tok, err := l.Current()

		if err != nil {
			return nil, err
		}

		for tok.Value != lexer.PAREN_R.String() {
			arg, err := parseArgument(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			*args = append(*args, *arg)

			tok, err = l.Current()

			if err != nil {
				return nil, err
			}

		}

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}

		if tok.Value != lexer.PAREN_R.String() {
			return nil, errors.New("Expecting closing parentheses for arguments")
		}

		l.Get()

		return args, nil
	}
}

//
func parseArgument(l *lexer.Lexer) (*ast.Argument, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.COLON.String() {
		return nil, errors.New("Expecting colon after argument name")
	}

	l.Get()

	val, err := parseValue(l)

	if err != nil {
		return nil, err
	}

	arg := &ast.Argument{}

	arg.Name = *name
	arg.Value = val

	return arg, nil
}

//
func parseField(l *lexer.Lexer) (*ast.Field, error) {
	alias, err := parseName(l)

	if err != nil {
		return nil, err
	}

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	name := &ast.Name{}

	if tok.Value == lexer.COLON.String() {
		l.Get()

		name, err = parseName(l)

		if err != nil {
			return nil, err
		}
	} else {
		*name = *alias

		alias = nil
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

	field.Alias = (*ast.Alias)(alias)
	field.Name = *name
	field.Arguments = args
	field.Directives = dirs
	field.SelectionSet = selSet

	return field, nil
}

//
func parseFragmentSpread(l *lexer.Lexer) (*ast.FragmentSpread, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.SPREAD.String() {
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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != lexer.SPREAD.String() {
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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != "on" {
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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	intVal, err := strconv.ParseInt(tok.Value, 10, 64)

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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	floatVal, err := strconv.ParseFloat(tok.Value, 64)

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
	tok, _ := l.Get()

	sv := &ast.StringValue{}

	sv.Value = tok.Value

	return sv, nil
}

//
func parseBooleanValue(l *lexer.Lexer) (*ast.BooleanValue, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	boolVal, err := strconv.ParseBool(tok.Value)

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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != "null" {
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
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != "[" {
		return nil, errDoesntExist
	} else {
		l.Get()

		lstVal := &ast.ListValue{}

		tok, err := l.Current()

		if err != nil {
			return nil, err
		}

		for tok.Value != "]" {
			val, err := parseValue(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			lstVal.Values = append(lstVal.Values, val)

			tok, err = l.Current()

			if err != nil {
				return nil, err
			}
		}

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}

		if tok.Value != "]" {
			return nil, errors.New("Missing closing bracket for list value")
		}

		return lstVal, nil
	}
}

//
func parseObjectValue(l *lexer.Lexer) (*ast.ObjectValue, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != "{" {
		return nil, errDoesntExist
	} else {
		l.Get()

		objVal := &ast.ObjectValue{}

		tok, err := l.Current()

		if err != nil {
			return nil, err
		}

		for tok.Value != "}" {
			objField, err := parseObjectField(l)

			if err != nil {
				if err == errDoesntExist {
					break
				}

				return nil, err
			}

			objVal.Values = append(objVal.Values, *objField)

			tok, err = l.Current()

			if err != nil {
				return nil, err
			}
		}

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}

		if tok.Value != "}" {
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

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != ":" {
		return nil, errors.New("Expecting color after object field name")
	}

	l.Get()

	val, err := parseValue(l)

	if err != nil {
		return nil, err
	}

	objField := &ast.ObjectField{}

	objField.Name = *name
	objField.Value = val

	return objField, nil
}
