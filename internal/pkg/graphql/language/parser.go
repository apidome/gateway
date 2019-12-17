package language

import (
	"reflect"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

var (
	errDoesntExist error = errors.New("Element does not exist")
)

func isNilInterface(i interface{}) bool {
	return reflect.ValueOf(i).IsNil()
}

func Parse(doc string) (*Document, error) {
	l, err := NewLexer(doc)
	if err != nil {
		return nil, err
	}

	astDoc, _ := parseDocument(l)

	return astDoc, nil
}

// https://graphql.github.io/graphql-spec/draft/#Document
func parseDocument(l *Lexer) (*Document, error) {
	def, err := parseDefinitions(l)

	if err != nil {
		return nil, err
	}

	doc := &Document{}

	doc.Definitions = *def

	return doc, nil
}

// https://graphql.github.io/graphql-spec/draft/#Definition
func parseDefinitions(l *Lexer) (*Definitions, error) {
	defs := &Definitions{}

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	for tok.Value != EOF.String() {
		var def Definition

		def, err := parseDefinition(l)

		if err != nil {
			if err == errDoesntExist {
				break
			}

			return nil, err
		}

		if def != nil {
			*defs = append(*defs, def)
		}

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}
	}

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value == EOF.String() {
		l.Get()
	}

	if len(*defs) == 0 {
		return nil, errors.New("No definitions found in document")
	}

	return defs, nil
}

// https://graphql.github.io/graphql-spec/draft/#Definition
func parseDefinition(l *Lexer) (Definition, error) {
	def, err := parseExecutableDefinition(l)

	if err == nil {
		return def, nil
	}

	def, err = parseTypeSystemDefinition(l)

	if err == nil {
		return def, nil
	}

	def, err = parseTypeSystemExtension(l)

	if err != nil {
		return nil, err
	}

	return def, nil
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDefinition
func parseExecutableDefinition(l *Lexer) (ExecutableDefinition, error) {
	var execDef ExecutableDefinition

	execDef, err := parseOperationDefinition(l)

	if err == nil {
		return execDef, nil
	}

	execDef, err = parseFragmentDefinition(l)

	if err != nil {
		return nil, err
	}

	return execDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDefinition
func parseTypeSystemDefinition(l *Lexer) (TypeSystemDefinition, error) {
	var def TypeSystemDefinition

	def, err := parseSchemaDefinition(l)

	if err == nil {
		return def, nil
	}

	def, err = parseTypeDefinition(l)

	if err == nil {
		return def, nil
	}

	def, err = parseDirectiveDefinition(l)

	if err != nil {
		return nil, err
	}

	return def, nil
}

// https://graphql.github.io/graphql-spec/draft/#SchemaDefinition
func parseSchemaDefinition(l *Lexer) (*SchemaDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != KW_SCHEMA {
		return nil, errors.New("Missing 'schema' keyword for a schema definition")
	}

	l.Get()

	dirs, _ := parseDirectives(l)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACE_L.String() {
		return nil, errors.New("Missing '{' for a schema definition")
	}

	l.Get()

	rOtd, err := parseRootOperationTypeDefinitions(l)

	if err != nil {
		return nil, err
	}

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACE_R.String() {
		return nil, errors.New("Missing '}' for schema definition")
	}

	l.Get()

	schDef := &SchemaDefinition{}

	schDef.Directives = dirs
	schDef.RootOperationTypeDefinitions = *rOtd
	schDef.Loc = Location{locStart, tok.End, l.Source()}

	return schDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#RootOperationTypeDefinition
func parseRootOperationTypeDefinitions(l *Lexer) (*RootOperationTypeDefinitions, error) {
	rotds := &RootOperationTypeDefinitions{}

	for rotd, err := parseRootOperationTypeDefinition(l); err != nil; rotd, err = parseRootOperationTypeDefinition(l) {
		*rotds = append(*rotds, *rotd)
	}

	return rotds, nil
}

// https://graphql.github.io/graphql-spec/draft/#RootOperationTypeDefinition
func parseRootOperationTypeDefinition(l *Lexer) (*RootOperationTypeDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	opType, err := parseOperationType(l)

	if err != nil {
		return nil, err
	}

	tok, err = l.Current()

	if tok.Value != COLON.String() {
		return nil, errors.New("Expecting ':' after operation type")
	}

	l.Get()

	namedType, err := parseNamedType(l)

	if err != nil {
		return nil, err
	}

	rotd := &RootOperationTypeDefinition{}

	rotd.OperationType = *opType
	rotd.NamedType = *namedType
	rotd.Loc = Location{locStart, namedType.Location().End, l.Source()}

	return rotd, nil
}

// https://graphql.github.io/graphql-spec/draft/#OperationType
func parseOperationType(l *Lexer) (*OperationType, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != string(OPERATION_MUTATION) &&
		tok.Value != string(OPERATION_QUERY) &&
		tok.Value != string(OPERATION_SUBSCRIPTION) {
		return nil,
			errors.New("Expecting 'query', 'mutation' or 'subscription' as operation type")
	}

	opType := new(OperationType)

	*opType = (OperationType)(tok.Value)

	return opType, nil
}

// ! Type system !! continue from here
// https://graphql.github.io/graphql-spec/draft/#TypeDefinition
func parseTypeDefinition(l *Lexer) (TypeDefinition, error) {
	scalarTd, err := parseScalarTypeDefinition(l)

	if err != nil {
		return nil, err
	} else {
		return scalarTd, nil
	}

	objectTd, err := parseObjectTypeDefinition(l)

	if err != nil {
		return nil, err
	} else {
		return objectTd, err
	}

	interfaceTd, err := parseInterfaceTypeDefinition(l)

	if err != nil {
		return nil, err
	} else {
		return interfaceTd, err
	}

	unionTd, err := parseUnionTypeDefinition(l)

	if err != nil {
		return nil, err
	} else {
		return unionTd, nil
	}

	enumTd, err := parseEnumTypeDefinition(l)

	if err != nil {
		return nil, err
	} else {
		return enumTd, nil
	}

	inputTd, err := parseInputObjectTypeDefinition(l)

	if err != nil {
		return nil, err
	} else {
		return inputTd, nil
	}

	return nil, errors.New("No type definition found")
}

// https://graphql.github.io/graphql-spec/draft/#ScalarTypeDefinition
func parseScalarTypeDefinition(l *Lexer) (*ScalarTypeDefinition, error) {
	desc, _ := parseDescription(l)

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != SCALAR {
		return nil, errors.New("Missing 'scalar' keyword for scalar type definition")
	}

	l.Get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	locEnd := name.Location().End

	dirs, _ := parseDirectives(l)

	if dirs != nil {
		if len(*dirs) > 0 {
			locEnd = (*dirs)[len(*dirs)-1].Location().End
		}
	}
	scalarTd := &ScalarTypeDefinition{}

	scalarTd.Description = desc
	scalarTd.Name = *name
	scalarTd.Directives = dirs
	scalarTd.Loc = Location{tok.Start, locEnd, l.Source()}

	return scalarTd, nil
}

// https://graphql.github.io/graphql-spec/draft/#Description
func parseDescription(l *Lexer) (*Description, error) {
	strVal, err := parseStringValue(l)

	if err != nil {
		return nil, err
	}

	return (*Description)(strVal), nil
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeDefinition
func parseObjectTypeDefinition(l *Lexer) (*ObjectTypeDefinition, error) {
	desc, _ := parseDescription(l)

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != KW_TYPE {
		return nil, errors.New("Expecting 'type' keyword for object type definition")
	}

	l.Get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	ii, _ := parseImplementsInterfaces(l)

	locEnd := 0

	if ii != nil {
		if len(*ii) > 0 {
			locEnd = (*ii)[len(*ii)-1].Location().End
		}
	}

	dirs, _ := parseDirectives(l)

	if dirs != nil {
		if len(*dirs) > 0 {
			locEnd = (*dirs)[len(*dirs)-1].Location().End
		}
	}

	fd, _ := parseFieldsDefinition(l)

	if fd != nil {
		if len(*fd) > 0 {
			locEnd = (*fd)[len(*fd)-1].Location().End
		}
	}

	locStart := 0

	if desc != nil {
		locStart = desc.Location().Start
	} else {
		locStart = tok.Start
	}

	objTd := &ObjectTypeDefinition{}

	objTd.Description = desc
	objTd.Directives = dirs
	objTd.FieldsDefinition = fd
	objTd.ImplementsInterfaces = ii
	objTd.Name = *name
	objTd.Loc = Location{locStart, locEnd, l.Source()}

	return objTd, nil
}

// https://graphql.github.io/graphql-spec/draft/#ImplementsInterfaces
func parseImplementsInterfaces(l *Lexer) (*ImplementsInterfaces, error) {

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != KW_IMPLEMENTS {
		return nil, errors.New("Expecting 'implements' keyword")
	}

	l.Get()

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value == AMP.String() {
		l.Get()
	}

	nt, err := parseNamedType(l)

	if err != nil {
		return nil, err
	}

	ii := &ImplementsInterfaces{}

	(*ii) = append(*ii, *nt)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	for tok.Value == AMP.String() {
		l.Get()

		nt, err := parseNamedType(l)

		if err != nil {
			return nil, err
		}

		(*ii) = append(*ii, *nt)

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}
	}

	return ii, nil
}

// https://graphql.github.io/graphql-spec/draft/#FieldsDefinition
func parseFieldsDefinition(l *Lexer) (*FieldsDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACE_L.String() {
		return nil, errors.New("Expecting '{' for fields definition")
	}

	l.Get()

	fds := &FieldsDefinition{}

	for fd, err := parseFieldDefinition(l); err != nil; fd, err = parseFieldDefinition(l) {
		(*fds) = append(*fds, *fd)
	}

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACE_R.String() {
		return nil, errors.New("Expecting '}' for fields definition")
	}

	l.Get()

	return fds, nil
}

// https://graphql.github.io/graphql-spec/draft/#FieldsDefinition
func parseFieldDefinition(l *Lexer) (*FieldDefinition, error) {
	desc, _ := parseDescription(l)

	locStart := 0

	if desc != nil {
		locStart = desc.Location().Start
	}

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	argsDef, _ := parseArgumentsDefinition(l)

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != COLON.String() {
		return nil, errors.New("Expecting ':' for field definition")
	}

	l.Get()

	_type, err := parseType(l)

	if err != nil {
		return nil, err
	}

	locEnd := _type.Location().End

	dirs, _ := parseDirectives(l)

	if dirs != nil {
		if len(*dirs) > 0 {
			locEnd = (*dirs)[len(*dirs)-1].Location().End
		}
	}

	fd := &FieldDefinition{}

	fd.Description = desc
	fd.Name = *name
	fd.ArgumentsDefinition = argsDef
	fd.Type = _type
	fd.Directives = dirs
	fd.Loc = Location{locStart, locEnd, l.Source()}

	return fd, nil
}

// https://graphql.github.io/graphql-spec/draft/#ArgumentsDefinition
func parseArgumentsDefinition(l *Lexer) (*ArgumentsDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != PAREN_L.String() {
		return nil, errors.New("Expecting '(' for arguments definition")
	}

	l.Get()

	argsDef := &ArgumentsDefinition{}

	for ivDef, err := parseInputValueDefinition(l); err != nil; ivDef, err = parseInputValueDefinition(l) {
		*argsDef = append(*argsDef, *ivDef)
	}

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != PAREN_R.String() {
		return nil, errors.New("Expecting ')' for arguments definition")
	}

	l.Get()

	return argsDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputValueDefinition
func parseInputValueDefinition(l *Lexer) (*InputValueDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	desc, _ := parseDescription(l)

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != COLON.String() {
		return nil, errors.New("Expecting ':' for input value definition")
	}

	l.Get()

	_type, err := parseType(l)

	if err != nil {
		return nil, err
	}

	defVal, _ := parseDefaultValue(l)

	dirs, _ := parseDirectives(l)

	ivDef := &InputValueDefinition{}

	ivDef.Description = desc
	ivDef.Name = *name
	ivDef.Type = _type
	ivDef.DefaultValue = defVal
	ivDef.Directives = dirs
	ivDef.Loc = Location{locStart, l.PrevLocation().End, l.Source()}

	return ivDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeDefinition
func parseInterfaceTypeDefinition(l *Lexer) (*InterfaceTypeDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	desc, _ := parseDescription(l)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != INTERFACE {
		return nil, errors.New("Expecting 'interface' keyword for interface type definition")
	}

	l.Get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	fds, _ := parseFieldsDefinition(l)

	itd := &InterfaceTypeDefinition{}

	itd.Description = desc
	itd.Directives = dirs
	itd.FieldsDefinition = fds
	itd.Name = *name
	itd.Loc = Location{locStart, l.PrevLocation().End, l.Source()}

	return itd, nil
}

func parseUnionTypeDefinition(l *Lexer) (*UnionTypeDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	desc, _ := parseDescription(l)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != UNION {
		return nil, errors.New("Expecting 'union' keyowrd for union type definition")
	}

	l.Get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	umt, _ := parseUnionMemberTypes(l)

	utd := &UnionTypeDefinition{}

	utd.Description = desc
	utd.Name = *name
	utd.Directives = dirs
	utd.UnionMemberTypes = umt
	utd.Loc = Location{locStart, l.PrevLocation().End, l.Source()}

	return utd, nil
}

// https://graphql.github.io/graphql-spec/draft/#UnionMemberTypes
func parseUnionMemberTypes(l *Lexer) (*UnionMemberTypes, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != EQUALS.String() {
		return nil, errors.New("Expecting '=' for union member types")
	}

	l.Get()

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value == PIPE.String() {
		l.Get()
	}

	nt, err := parseNamedType(l)

	if err != nil {
		return nil, err
	}

	umt := &UnionMemberTypes{}

	*umt = append(*umt, *nt)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	for tok.Value == PIPE.String() {
		l.Get()

		nt, err = parseNamedType(l)

		if err != nil {
			return nil, err
		}

		*umt = append(*umt, *nt)

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}
	}

	return umt, nil
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeDefinition
func parseEnumTypeDefinition(l *Lexer) (*EnumTypeDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	desc, _ := parseDescription(l)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != ENUM {
		return nil, errors.New("Expecting 'enum' keyword for enum type definition")
	}

	l.Get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	evd, _ := parseEnumValuesDefinition(l)

	etd := &EnumTypeDefinition{}

	etd.Description = desc
	etd.Name = *name
	etd.Directives = dirs
	etd.EnumValuesDefinition = evd
	etd.Loc = Location{locStart, l.PrevLocation().End, l.Source()}

	return etd, nil
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinition
func parseEnumValuesDefinition(l *Lexer) (*EnumValuesDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACE_L.String() {
		return nil, errors.New("Expecting '{' for enum values definition")
	}

	l.Get()

	evds := &EnumValuesDefinition{}

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	for tok.Value != BRACE_R.String() {
		evd, err := parseEnumValueDefinition(l)

		if err != nil {
			return nil, err
		}

		*evds = append(*evds, *evd)

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}
	}

	l.Get()

	return evds, nil
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinition
func parseEnumValueDefinition(l *Lexer) (*EnumValueDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	desc, _ := parseDescription(l)

	ev, err := parseEnumValue(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	evd := &EnumValueDefinition{}

	evd.Description = desc
	evd.EnumValue = *ev
	evd.Directives = dirs
	ev.Loc = Location{locStart, l.PrevLocation().End, l.Source()}

	return evd, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeDefinition
func parseInputObjectTypeDefinition(l *Lexer) (*InputObjectTypeDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	desc, _ := parseDescription(l)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != KW_INPUT {
		return nil, errors.New("Expecting 'input' keyword for input object type definition")
	}

	l.Get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	ifds, _ := parseInputFieldsDefinition(l)

	iotd := &InputObjectTypeDefinition{}

	iotd.Description = desc
	iotd.Directives = dirs
	iotd.Name = *name
	iotd.InputFieldsDefinition = ifds
	iotd.Loc = Location{locStart, l.PrevLocation().End, l.Source()}

	return iotd, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputFieldsDefinition
func parseInputFieldsDefinition(l *Lexer) (*InputFieldsDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACE_L.String() {
		return nil, errors.New("Expecting '{' for input fields definition")
	}

	l.Get()

	ifds := &InputFieldsDefinition{}

	for tok.Value != BRACE_R.String() {
		ivd, err := parseInputValueDefinition(l)

		if err != nil {
			return nil, err
		}

		*ifds = append(*ifds, *ivd)
	}

	l.Get()

	return ifds, nil
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveDefinition
func parseDirectiveDefinition(l *Lexer) (*DirectiveDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	desc, _ := parseDescription(l)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != KW_DIRECTIVE {
		return nil, errors.New("Expecting 'directive' keyword for directive definition")
	}

	l.Get()

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != AT.String() {
		return nil, errors.New("Expecting '@' for directive definition")
	}

	l.Get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	argsDef, _ := parseArgumentsDefinition(l)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != "on" {
		return nil, errors.New("Expecting 'on' keyworkd for directive definition")
	}

	l.Get()

	dls, err := parseDirectiveLocations(l)

	if err != nil {
		return nil, err
	}

	df := &DirectiveDefinition{}

	df.Description = desc
	df.Name = *name
	df.ArgumentsDefinition = argsDef
	df.DirectiveLocations = *dls
	df.Loc = Location{locStart, l.PrevLocation().End, l.Source()}

	return df, nil
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocations
func parseDirectiveLocations(l *Lexer) (*DirectiveLocations, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value == PIPE.String() {
		l.Get()
	}

	dls := &DirectiveLocations{}

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	for tok.Value != PIPE.String() {
		l.Get()

		tok, err := l.Current()

		if err != nil {
			return nil, err
		}

		if tok.Value != PIPE.String() {
			return nil, errors.New("Expecting '|' between directive locations")
		}

		l.Get()

		dl, err := parseDirectiveLocation(l)

		if err != nil {
			return nil, err
		}

		*dls = append(*dls, *dl)

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}
	}

	return dls, nil
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocation
func parseDirectiveLocation(l *Lexer) (*DirectiveLocation, error) {
	edl, err := parseExecutableDirectiveLocation(l)

	if err == nil {
		return (*DirectiveLocation)(edl), nil
	}

	tsdl, err := parseTypeSystemDirectiveLocation(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting a directive location")
	}

	return (*DirectiveLocation)(tsdl), nil
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDirectiveLocation
func parseExecutableDirectiveLocation(l *Lexer) (*ExecutableDirectiveLocation, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	for i := range executableDirectiveLocations {
		if string(executableDirectiveLocations[i]) == tok.Value {
			edl := executableDirectiveLocations[i]

			return &edl, nil
		}
	}

	return nil, errors.New("Expecting executable directive location")
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDirectiveLocation
func parseTypeSystemDirectiveLocation(l *Lexer) (*TypeSystemDirectiveLocation, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	for i := range typeSystemDirectiveLocations {
		if string(typeSystemDirectiveLocations[i]) == tok.Value {
			tsdl := typeSystemDirectiveLocations[i]

			return &tsdl, nil
		}
	}

	return nil, errors.New("Expecting type systen directive location")
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemExtension
func parseTypeSystemExtension(l *Lexer) (TypeSystemExtension, error) {
	se, err := parseSchemaExtension(l)

	if err == nil {
		return se, nil
	}

	te, err := parseTypeExtension(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting type system extension")
	}

	return te, nil
}

// ! need to create a commiting mechanism for lexer
// https://graphql.github.io/graphql-spec/draft/#SchemaExtension
func parseSchemaExtension(l *Lexer) (*SchemaExtension, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != KW_EXTEND {
		return nil, errors.New("Expecting 'extend' keyword for schema extension")
	}

	l.Get()

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != KW_SCHEMA {
		return nil, errors.New("Expecting 'schema' keyword for schema extension")
	}

	l.Get()

	dirs, _ := parseDirectives(l)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACE_L.String() {
		return nil, errors.New("Expecting '{' for schema extension")
	}

	l.Get()

	rotds, _ := parseRootOperationTypeDefinitions(l)

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACE_R.String() {
		return nil, errors.New("Expecting '}' for schema extension")
	}

	l.Get()

	if dirs == nil && rotds == nil {
		return nil, errors.New("Expecting directives or root operation type definitions for schema extension")
	}

	se := &SchemaExtension{}

	se.Directives = dirs
	se.RootOperationTypeDefinitions = rotds
	se.Loc = Location{locStart, l.PrevLocation().End, l.Source()}

	return se, nil
}

func parseTypeExtension(l *Lexer) (TypeExtension, error) {
	return nil, nil
}

// ! come back to this
// https://graphql.github.io/graphql-spec/draft/#OperationDefinition
func parseOperationDefinition(l *Lexer) (*OperationDefinition, error) {
	tok, err := l.Current()

	locStart := tok.Start

	if err != nil {
		return nil, err
	}

	// Shorthand query
	// https://graphql.github.io/graphql-spec/draft/#sec-Language.Operations.Query-shorthand
	if tok.Value == BRACE_L.String() {

		shorthandQuery, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		opDef := &OperationDefinition{}

		opDef.OperationType = KW_QUERY
		opDef.SelectionSet = *shorthandQuery
		opDef.Loc = Location{locStart, tok.End, l.Source()}

		return opDef, nil
	} else if tok.Value != KW_QUERY &&
		tok.Value != KW_MUTATION &&
		tok.Value != KW_SUBSCRIPTION {
		return nil, errDoesntExist
	} else {
		tok, err := l.Get()

		if err != nil {
			return nil, err
		}

		opType := tok.Value

		name, _ := parseName(l)

		varDef, _ := parseVariableDefinitions(l)

		directives, _ := parseDirectives(l)

		selSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		opDefinition := &OperationDefinition{}

		opDefinition.OperationType = OperationType(opType)
		opDefinition.Name = name
		opDefinition.VariableDefinitions = varDef
		opDefinition.Directives = directives
		opDefinition.SelectionSet = *selSet
		opDefinition.Loc = Location{locStart, tok.End, l.Source()}

		return opDefinition, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentDefinition
func parseFragmentDefinition(l *Lexer) (*FragmentDefinition, error) {
	tok, err := l.Current()

	locStart := tok.Start

	if err != nil {
		return nil, err
	}

	if tok.Value != KW_FRAGMENT {
		return nil, errors.New("Expecting fragment keyword")
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

		directives, _ := parseDirectives(l)

		selectionSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		fragDef := &FragmentDefinition{}

		fragDef.FragmentName = *name
		fragDef.TypeCondition = *typeCond
		fragDef.Directives = directives
		fragDef.SelectionSet = *selectionSet
		fragDef.Loc = Location{locStart, tok.End, l.Source()}

		return fragDef, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#Name
func parseName(l *Lexer) (*Name, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	pattern := "^[_A-Za-z][_0-9A-Za-z]*$"

	// If the current token is not a Name, return nil
	if tok.Kind != NAME {
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

	name := &Name{}

	// Populate the Name struct.
	name.Value = tok.Value
	name.Loc.Start = tok.Start
	name.Loc.End = tok.End
	name.Loc.Source = l.Source()

	// Return the AST Name object.
	return name, nil
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinition
func parseVariableDefinitions(l *Lexer) (*VariableDefinitions, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != PAREN_L.String() {
		return nil, errors.New("Expecting '(' opener for variable definitions")
	} else {
		l.Get()

		varDefs := &VariableDefinitions{}

		tok, err = l.Current()

		if err != nil {
			return nil, err
		}

		for tok.Value != PAREN_R.String() {
			varDef, err := parseVariableDefinition(l)

			if err != nil {
				break
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

		if tok.Value != PAREN_R.String() {
			return nil, errors.New("Expecting closing parentheses for variable definitions")
		}

		l.Get()

		return varDefs, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinition
func parseVariableDefinition(l *Lexer) (*VariableDefinition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	_var, err := parseVariable(l)

	if err != nil {
		return nil, err
	}

	locStart := _var.Location().Start

	if tok.Value != COLON.String() {
		return nil, errors.New("Expecting a colon after variable name")
	}

	_type, err := parseType(l)

	if err != nil {
		return nil, err
	}

	locEnd := _type.Location().End

	defVal, _ := parseDefaultValue(l)

	directives, _ := parseDirectives(l)

	varDef := &VariableDefinition{}

	varDef.Variable = *_var
	varDef.Type = _type
	varDef.DefaultValue = defVal
	varDef.Directives = directives
	varDef.Loc = Location{locStart, locEnd, l.Source()}

	return varDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#Type
func parseType(l *Lexer) (Type, error) {
	var _type Type

	_type, err := parseNamedType(l)

	if err == nil {
		return _type, nil
	}

	_type, err = parseListType(l)

	if err == nil {
		return _type, nil
	}

	_type, err = parseNonNullType(l)

	if err == nil {
		return _type, nil
	} else {
		return nil, err
	}
}

// https://graphql.github.io/graphql-spec/draft/#ListType
func parseListType(l *Lexer) (*ListType, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != BRACKET_L.String() {
		return nil, errors.New("Expecting '[' for list type")
	}

	l.Get()

	_type, err := parseType(l)

	if err != nil {
		return nil, err
	}

	tok, err = l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACKET_R.String() {
		return nil, errors.New("Expecting ']' for list type")
	}

	locEnd := tok.End

	l.Get()

	listType := &ListType{}

	listType.OfType = _type
	listType.Loc = Location{locStart, locEnd, l.Source()}

	return listType, nil
}

// https://graphql.github.io/graphql-spec/draft/#NonNullType
func parseNonNullType(l *Lexer) (*NonNullType, error) {
	var _type Type

	_type, err := parseNamedType(l)

	if err != nil {
		_type, err = parseListType(l)

		if err != nil {
			return nil, err
		}
	}

	locStart := _type.Location().Start

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BANG.String() {
		return nil, errors.New("Expecting '!' at the end of a non null type")
	}

	locEnd := tok.End

	l.Get()

	nonNull := &NonNullType{}

	nonNull.OfType = _type
	nonNull.Loc = Location{locStart, locEnd, l.Source()}

	return nonNull, nil
}

// https://graphql.github.io/graphql-spec/draft/#Directives
func parseDirectives(l *Lexer) (*Directives, error) {
	dirs := &Directives{}

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

// https://graphql.github.io/graphql-spec/draft/#Directive
func parseDirective(l *Lexer) (*Directive, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != AT.String() {
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

		locEnd := 0

		if err == nil {
			locEnd = (*args)[len(*args)-1].Location().End
		} else {
			locEnd = tok.End
		}

		dir := &Directive{}

		dir.Name = *name
		dir.Arguments = args
		dir.Loc = Location{locStart, locEnd, l.Source()}

		return dir, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#SelectionSet
func parseSelectionSet(l *Lexer) (*SelectionSet, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != BRACE_L.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		selSet := &SelectionSet{}

		tok, err := l.Current()

		if err != nil {
			return nil, err
		}

		for tok.Value != BRACE_R.String() {
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

		if tok.Value != BRACE_R.String() {
			return nil, errors.New("Expecting closing bracket for selection set")
		}

		l.Get()

		return selSet, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#Selection
func parseSelection(l *Lexer) (Selection, error) {
	var sel Selection

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

// https://graphql.github.io/graphql-spec/draft/#Variable
func parseVariable(l *Lexer) (*Variable, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != DOLLAR.String() {
		return nil, errDoesntExist
	} else {
		name, err := parseName(l)

		if err != nil {
			return nil, err
		}

		_var := &Variable{}

		_var.Name = *name
		_var.Loc = Location{locStart, name.Location().End, l.Source()}

		return _var, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#DefaultValue
func parseDefaultValue(l *Lexer) (*DefaultValue, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != EQUALS.String() {
		return nil, errDoesntExist
	} else {
		val, err := parseValue(l)

		if err != nil {
			return nil, err
		}

		dVal := &DefaultValue{}

		dVal.Value = val
		dVal.Loc = Location{locStart, val.Location().End, l.Source()}

		return dVal, nil
	}
}

// ! need to check variable type in order to parse its value
// https://graphql.github.io/graphql-spec/draft/#Value
func parseValue(l *Lexer) (Value, error) {
	_var, err := parseVariable(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	} else {
		// ! need to read variable value if it exists
		_var = _var
	}

	var val Value

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

// https://graphql.github.io/graphql-spec/draft/#Arguments
func parseArguments(l *Lexer) (*Arguments, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != PAREN_L.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		args := &Arguments{}

		tok, err := l.Current()

		if err != nil {
			return nil, err
		}

		for tok.Value != PAREN_R.String() {
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

		if tok.Value != PAREN_R.String() {
			return nil, errors.New("Expecting closing parentheses for arguments")
		}

		l.Get()

		return args, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#Argument
func parseArgument(l *Lexer) (*Argument, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	locStart := name.Location().Start

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != COLON.String() {
		return nil, errors.New("Expecting colon after argument name")
	}

	l.Get()

	val, err := parseValue(l)

	if err != nil {
		return nil, err
	}

	arg := &Argument{}

	arg.Name = *name
	arg.Value = val
	arg.Loc = Location{locStart, val.Location().End, l.Source()}

	return arg, nil
}

// https://graphql.github.io/graphql-spec/draft/#Field
func parseField(l *Lexer) (*Field, error) {
	alias, err := parseName(l)

	if err != nil {
		return nil, err
	}

	locStart := alias.Location().Start

	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	name := &Name{}

	if tok.Value == COLON.String() {
		l.Get()

		name, err = parseName(l)

		if err != nil {
			return nil, err
		}
	} else {
		*name = *alias

		alias = nil
	}

	locEnd := name.Location().End

	args, err := parseArguments(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	} else {
		locEnd = (*args)[len(*args)-1].Location().End
	}

	dirs, err := parseDirectives(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	} else {
		locEnd = (*dirs)[len(*dirs)-1].Location().End
	}

	selSet, err := parseSelectionSet(l)

	if err != nil {
		if err != errDoesntExist {
			return nil, err
		}
	} else {
		locEnd = (*selSet)[len(*selSet)-1].Location().End
	}

	field := &Field{}

	field.Alias = (*Alias)(alias)
	field.Name = *name
	field.Arguments = args
	field.Directives = dirs
	field.SelectionSet = selSet
	field.Loc = Location{locStart, locEnd, l.Source()}

	return field, nil
}

// https://graphql.github.io/graphql-spec/draft/#FragmentSpread
func parseFragmentSpread(l *Lexer) (*FragmentSpread, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != SPREAD.String() {
		return nil, errDoesntExist
	} else {
		l.Get()

		fname, err := parseFragmentName(l)

		if err != nil {
			return nil, err
		}

		locEnd := fname.Location().End

		directives, err := parseDirectives(l)

		if err != nil {
			if err != errDoesntExist {
				return nil, err
			}
		} else {
			locEnd = (*directives)[len(*directives)-1].Location().End
		}

		spread := &FragmentSpread{}

		spread.FragmentName = *fname
		spread.Directives = directives
		spread.Loc = Location{locStart, locEnd, l.Source()}

		return spread, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#InlineFragment
func parseInlineFragment(l *Lexer) (*InlineFragment, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != SPREAD.String() {
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

		locEnd := (*selSet)[len(*selSet)-1].Location().End

		inlineFrag := &InlineFragment{}

		inlineFrag.TypeCondition = typeCon
		inlineFrag.Directives = directives
		inlineFrag.SelectionSet = *selSet
		inlineFrag.Loc = Location{locStart, locEnd, l.Source()}

		return inlineFrag, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentName
func parseFragmentName(l *Lexer) (*FragmentName, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if name.Value == "on" {
		return nil, errors.New("Fragment name cannot be 'on'")
	}

	var fragName *FragmentName

	*fragName = FragmentName(*name)
	fragName.Loc = *name.Location()

	return fragName, nil
}

// https://graphql.github.io/graphql-spec/draft/#TypeCondition
func parseTypeCondition(l *Lexer) (*TypeCondition, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != "on" {
		return nil, errDoesntExist
	} else {
		namedType, err := parseNamedType(l)

		if err != nil {
			return nil, err
		}

		typeCond := &TypeCondition{}

		typeCond.NamedType = *namedType
		typeCond.Loc = Location{locStart, namedType.Location().End, l.Source()}

		return typeCond, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#NamedType
func parseNamedType(l *Lexer) (*NamedType, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	var namedType *NamedType

	*namedType = NamedType(*name)

	return namedType, nil
}

// https://graphql.github.io/graphql-spec/draft/#IntValue
func parseIntValue(l *Lexer) (*IntValue, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	intVal, err := strconv.ParseInt(tok.Value, 10, 64)

	if err != nil {
		return nil, err
	}

	l.Get()

	intValP := &IntValue{}

	intValP.Value = intVal
	intValP.Loc = Location{tok.Start, tok.End, l.Source()}

	return intValP, nil
}

// https://graphql.github.io/graphql-spec/draft/#FloatValue
func parseFloatValue(l *Lexer) (*FloatValue, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	floatVal, err := strconv.ParseFloat(tok.Value, 64)

	if err != nil {
		return nil, err
	}

	l.Get()

	floatValP := &FloatValue{}

	floatValP.Value = floatVal
	floatValP.Loc = Location{tok.Start, tok.End, l.Source()}

	return floatValP, nil
}

// ! Have a discussion about this function
// https://graphql.github.io/graphql-spec/draft/#StringValue
func parseStringValue(l *Lexer) (*StringValue, error) {
	tok, _ := l.Get()

	sv := &StringValue{}

	sv.Value = tok.Value
	sv.Loc = Location{tok.Start, tok.End, l.Source()}

	return sv, nil
}

// https://graphql.github.io/graphql-spec/draft/#BooleanValue
func parseBooleanValue(l *Lexer) (*BooleanValue, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	boolVal, err := strconv.ParseBool(tok.Value)

	if err != nil {
		return nil, err
	}

	l.Get()

	boolValP := &BooleanValue{}

	boolValP.Value = boolVal
	boolValP.Loc = Location{tok.Start, tok.End, l.Source()}

	return boolValP, nil
}

// ! Figure out what to do with a null value
// https://graphql.github.io/graphql-spec/draft/#NullValue
func parseNullValue(l *Lexer) (*NullValue, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	if tok.Value != "null" {
		return nil, errDoesntExist
	} else {
		l.Get()

		null := &NullValue{}
		null.Loc = Location{tok.Start, tok.End, l.Source()}

		return null, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#EnumValue
func parseEnumValue(l *Lexer) (*EnumValue, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	switch name.Value {
	case "true", "false", "null":
		return nil, errors.New("Enum value cannot be 'true', 'false' or 'null'")
	default:
		enumVal := &EnumValue{}

		enumVal.Name = *name
		enumVal.Loc = Location{name.Location().Start, name.Location().End, l.Source()}

		return enumVal, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#ListValue
func parseListValue(l *Lexer) (*ListValue, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != "[" {
		return nil, errDoesntExist
	} else {
		l.Get()

		lstVal := &ListValue{}

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

		locEnd := tok.End

		if tok.Value != "]" {
			return nil, errors.New("Missing closing bracket for list value")
		}

		lstVal.Loc = Location{locStart, locEnd, l.Source()}

		return lstVal, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectValue
func parseObjectValue(l *Lexer) (*ObjectValue, error) {
	tok, err := l.Current()

	if err != nil {
		return nil, err
	}

	locStart := tok.Start

	if tok.Value != "{" {
		return nil, errDoesntExist
	} else {
		l.Get()

		objVal := &ObjectValue{}

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

		objVal.Loc = Location{locStart, tok.End, l.Source()}

		return objVal, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectField
func parseObjectField(l *Lexer) (*ObjectField, error) {
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

	objField := &ObjectField{}

	objField.Name = *name
	objField.Value = val
	objField.Loc = Location{name.Location().Start, val.Location().End, l.Source()}

	return objField, nil
}
