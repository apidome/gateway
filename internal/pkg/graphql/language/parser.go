package language

import (
	"reflect"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

func isNilInterface(i interface{}) bool {
	return reflect.ValueOf(i).IsNil()
}

func Parse(doc string) (*Document, error) {
	l, err := newlexer(doc)
	if err != nil {
		return nil, err
	}

	astDoc, err := parseDocument(l)

	if err != nil {
		return nil, err
	}

	return astDoc, nil
}

// https://graphql.github.io/graphql-spec/draft/#Document
func parseDocument(l *lexer) (*Document, error) {
	def, err := parseDefinitions(l)

	if err != nil {
		return nil, err
	}

	doc := &Document{}

	doc.Definitions = *def

	return doc, nil
}

// https://graphql.github.io/graphql-spec/draft/#Definition
func parseDefinitions(l *lexer) (*Definitions, error) {
	defs := &Definitions{}

	for !l.tokenEquals(EOF.string()) {
		var def Definition

		def, err := parseDefinition(l)

		if err != nil {
			return nil, err
		}

		if def != nil {
			*defs = append(*defs, def)
		}
	}

	if len(*defs) == 0 {
		return nil, errors.New("No definitions found in document")
	}

	return defs, nil
}

// https://graphql.github.io/graphql-spec/draft/#Definition
func parseDefinition(l *lexer) (Definition, error) {
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
		return nil,
			errors.Wrap(err, "Expecting one of 'executable definition', 'type system definition', 'type system extension'")
	}

	return def, nil
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDefinition
func parseExecutableDefinition(l *lexer) (ExecutableDefinition, error) {
	var execDef ExecutableDefinition

	execDef, err := parseOperationDefinition(l)

	if err == nil {
		return execDef, nil
	}

	execDef, err = parseFragmentDefinition(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting one of 'operation definition', 'fragment definition'")
	}

	return execDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDefinition
func parseTypeSystemDefinition(l *lexer) (TypeSystemDefinition, error) {
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
		return nil,
			errors.Wrap(err, "Expecting one of 'schema definition', 'type definition', 'directive definition'")
	}

	return def, nil
}

// https://graphql.github.io/graphql-spec/draft/#SchemaDefinition
func parseSchemaDefinition(l *lexer) (*SchemaDefinition, error) {
	locStart := l.current().Start

	if !l.tokenEquals(KW_SCHEMA) {
		return nil, errors.New("Missing 'schema' keyword for a schema definition")
	}

	l.get()

	dirs, _ := parseDirectives(l)

	if !l.tokenEquals(BRACE_L.string()) {
		return nil, errors.New("Missing '{' for a schema definition")
	}

	l.get()

	rOtd, err := parseRootOperationTypeDefinitions(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(BRACE_R.string()) {
		return nil, errors.New("Missing '}' for schema definition")
	}

	locEnd := l.current().End

	l.get()

	schDef := &SchemaDefinition{}

	schDef.Directives = dirs
	schDef.RootOperationTypeDefinitions = *rOtd
	schDef.Loc = Location{locStart, locEnd, l.source}

	return schDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#RootOperationTypeDefinition
func parseRootOperationTypeDefinitions(l *lexer) (*RootOperationTypeDefinitions, error) {
	rotds := &RootOperationTypeDefinitions{}

	for !l.tokenEquals(BRACE_R.string()) {
		rotd, err := parseRootOperationTypeDefinition(l)

		if err != nil {
			return nil, err
		}

		*rotds = append(*rotds, *rotd)
	}

	if len(*rotds) == 0 {
		return nil, errors.New("Expecting at least one root operation type definition")
	}

	return rotds, nil
}

// https://graphql.github.io/graphql-spec/draft/#RootOperationTypeDefinition
func parseRootOperationTypeDefinition(l *lexer) (*RootOperationTypeDefinition, error) {
	locStart := l.current().Start

	opType, err := parseOperationType(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(COLON.string()) {
		return nil, errors.New("Expecting ':' after operation type")
	}

	l.get()

	namedType, err := parseNamedType(l)

	if err != nil {
		return nil, err
	}

	rotd := &RootOperationTypeDefinition{}

	rotd.OperationType = *opType
	rotd.NamedType = *namedType
	rotd.Loc = Location{locStart, namedType.Location().End, l.source}

	return rotd, nil
}

// https://graphql.github.io/graphql-spec/draft/#OperationType
func parseOperationType(l *lexer) (*OperationType, error) {
	tok := l.current()

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

// https://graphql.github.io/graphql-spec/draft/#TypeDefinition
func parseTypeDefinition(l *lexer) (TypeDefinition, error) {
	scalarTd, err := parseScalarTypeDefinition(l)

	if scalarTd != nil {
		return scalarTd, nil
	}

	objectTd, err := parseObjectTypeDefinition(l)

	if objectTd != nil {
		return objectTd, nil
	}

	interfaceTd, err := parseInterfaceTypeDefinition(l)

	if interfaceTd != nil {
		return interfaceTd, nil
	}

	unionTd, err := parseUnionTypeDefinition(l)

	if unionTd != nil {
		return unionTd, nil
	}

	enumTd, err := parseEnumTypeDefinition(l)

	if enumTd != nil {
		return enumTd, nil
	}

	inputTd, err := parseInputObjectTypeDefinition(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting a type definition")
	} else {
		return inputTd, nil
	}

	return nil, errors.New("No type definition found")
}

// https://graphql.github.io/graphql-spec/draft/#ScalarTypeDefinition
func parseScalarTypeDefinition(l *lexer) (*ScalarTypeDefinition, error) {
	desc, _ := parseDescription(l)

	if !l.tokenEquals(SCALAR) {
		return nil, errors.New("Missing 'scalar' keyword for scalar type definition")
	}

	tok := l.get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	scalarTd := &ScalarTypeDefinition{}

	scalarTd.Description = desc
	scalarTd.Name = *name
	scalarTd.Directives = dirs
	scalarTd.Loc = Location{tok.Start, l.prevLocation().End, l.source}

	return scalarTd, nil
}

// https://graphql.github.io/graphql-spec/draft/#Description
func parseDescription(l *lexer) (*Description, error) {
	strVal, err := parseStringValue(l)

	if err != nil {
		return nil, err
	}

	return (*Description)(strVal), nil
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeDefinition
func parseObjectTypeDefinition(l *lexer) (*ObjectTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(KW_TYPE) {
		return nil, errors.New("Expecting 'type' keyword for object type definition")
	}

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	ii, _ := parseImplementsInterfaces(l)

	dirs, _ := parseDirectives(l)

	fd, _ := parseFieldsDefinition(l)

	objTd := &ObjectTypeDefinition{}

	objTd.Description = desc
	objTd.Directives = dirs
	objTd.FieldsDefinition = fd
	objTd.ImplementsInterfaces = ii
	objTd.Name = *name
	objTd.Loc = Location{locStart, l.prevLocation().End, l.source}

	return objTd, nil
}

// https://graphql.github.io/graphql-spec/draft/#ImplementsInterfaces
func parseImplementsInterfaces(l *lexer) (*ImplementsInterfaces, error) {
	if !l.tokenEquals(KW_IMPLEMENTS) {
		return nil, errors.New("Expecting 'implements' keyword")
	}

	if l.tokenEquals(AMP.string()) {
		l.get()
	}

	nt, err := parseNamedType(l)

	if err != nil {
		return nil, err
	}

	ii := &ImplementsInterfaces{}

	(*ii) = append(*ii, *nt)

	for l.tokenEquals(AMP.string()) {
		l.get()

		nt, err := parseNamedType(l)

		if err != nil {
			return nil, err
		}

		(*ii) = append(*ii, *nt)
	}

	return ii, nil
}

// https://graphql.github.io/graphql-spec/draft/#FieldsDefinition
func parseFieldsDefinition(l *lexer) (*FieldsDefinition, error) {
	if l.tokenEquals(BRACE_L.string()) {
		return nil, errors.New("Expecting '{' for fields definition")
	}

	l.get()

	fds := &FieldsDefinition{}

	for !l.tokenEquals(BRACE_R.string()) {
		fd, err := parseFieldDefinition(l)

		if err != nil {
			return nil, err
		}

		(*fds) = append(*fds, *fd)
	}

	l.get()

	if len(*fds) == 0 {
		return nil, errors.New("Expecting at lease one field definition")
	}

	return fds, nil
}

// https://graphql.github.io/graphql-spec/draft/#FieldsDefinition
func parseFieldDefinition(l *lexer) (*FieldDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	argsDef, _ := parseArgumentsDefinition(l)

	if !l.tokenEquals(COLON.string()) {
		return nil, errors.New("Expecting ':' for field definition")
	}

	l.get()

	_type, err := parseType(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	fd := &FieldDefinition{}

	fd.Description = desc
	fd.Name = *name
	fd.ArgumentsDefinition = argsDef
	fd.Type = _type
	fd.Directives = dirs
	fd.Loc = Location{locStart, l.prevLocation().End, l.source}

	return fd, nil
}

// https://graphql.github.io/graphql-spec/draft/#ArgumentsDefinition
func parseArgumentsDefinition(l *lexer) (*ArgumentsDefinition, error) {
	if !l.tokenEquals(PAREN_L.string()) {
		return nil, errors.New("Expecting '(' for arguments definition")
	}

	l.get()

	argsDef := &ArgumentsDefinition{}

	for !l.tokenEquals(PAREN_R.string()) {
		ivDef, err := parseInputValueDefinition(l)

		if err != nil {
			return nil, err
		}

		*argsDef = append(*argsDef, *ivDef)
	}

	l.get()

	if len(*argsDef) == 0 {
		return nil, errors.New("Expecting at least one input value definitions")
	}

	return argsDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputValueDefinition
func parseInputValueDefinition(l *lexer) (*InputValueDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(COLON.string()) {
		return nil, errors.New("Expecting ':' for input value definition")
	}

	l.get()

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
	ivDef.Loc = Location{locStart, l.prevLocation().End, l.source}

	return ivDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeDefinition
func parseInterfaceTypeDefinition(l *lexer) (*InterfaceTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(INTERFACE) {
		return nil, errors.New("Expecting 'interface' keyword for interface type definition")
	}

	l.get()

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
	itd.Loc = Location{locStart, l.prevLocation().End, l.source}

	return itd, nil
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeDefinition
func parseUnionTypeDefinition(l *lexer) (*UnionTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(UNION) {
		return nil, errors.New("Expecting 'union' keyowrd for union type definition")
	}

	l.get()

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
	utd.Loc = Location{locStart, l.prevLocation().End, l.source}

	return utd, nil
}

// https://graphql.github.io/graphql-spec/draft/#UnionMemberTypes
func parseUnionMemberTypes(l *lexer) (*UnionMemberTypes, error) {
	if !l.tokenEquals(EQUALS.string()) {
		return nil, errors.New("Expecting '=' for union member types")
	}

	l.get()

	if l.tokenEquals(PIPE.string()) {
		l.get()
	}

	nt, err := parseNamedType(l)

	if err != nil {
		return nil, err
	}

	umt := &UnionMemberTypes{}

	*umt = append(*umt, *nt)

	for l.tokenEquals(PIPE.string()) {
		l.get()

		nt, err = parseNamedType(l)

		if err != nil {
			return nil, err
		}

		*umt = append(*umt, *nt)
	}

	if len(*umt) == 0 {
		return nil, errors.New("Expecting at least one union member type")
	}

	return umt, nil
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeDefinition
func parseEnumTypeDefinition(l *lexer) (*EnumTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(ENUM) {
		return nil, errors.New("Expecting 'enum' keyword for enum type definition")
	}

	l.get()

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
	etd.Loc = Location{locStart, l.prevLocation().End, l.source}

	return etd, nil
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinition
func parseEnumValuesDefinition(l *lexer) (*EnumValuesDefinition, error) {
	if !l.tokenEquals(BRACE_L.string()) {
		return nil, errors.New("Expecting '{' for enum values definition")
	}

	l.get()

	evds := &EnumValuesDefinition{}

	for !l.tokenEquals(BRACE_R.string()) {
		evd, err := parseEnumValueDefinition(l)

		if err != nil {
			return nil, err
		}

		*evds = append(*evds, *evd)
	}

	l.get()

	if len(*evds) == 0 {
		return nil, errors.New("Expecting at least one enum value definition")
	}

	return evds, nil
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinition
func parseEnumValueDefinition(l *lexer) (*EnumValueDefinition, error) {
	locStart := l.location().Start

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
	ev.Loc = Location{locStart, l.prevLocation().End, l.source}

	return evd, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeDefinition
func parseInputObjectTypeDefinition(l *lexer) (*InputObjectTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(KW_INPUT) {
		return nil, errors.New("Expecting 'input' keyword for input object type definition")
	}

	l.get()

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
	iotd.Loc = Location{locStart, l.prevLocation().End, l.source}

	return iotd, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputFieldsDefinition
func parseInputFieldsDefinition(l *lexer) (*InputFieldsDefinition, error) {
	if !l.tokenEquals(BRACE_L.string()) {
		return nil, errors.New("Expecting '{' for input fields definition")
	}

	l.get()

	ifds := &InputFieldsDefinition{}

	for !l.tokenEquals(BRACE_R.string()) {
		ivd, err := parseInputValueDefinition(l)

		if err != nil {
			return nil, err
		}

		*ifds = append(*ifds, *ivd)
	}

	l.get()

	if len(*ifds) == 0 {
		return nil, errors.New("Expecting at least one input field definition")
	}

	return ifds, nil
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveDefinition
func parseDirectiveDefinition(l *lexer) (*DirectiveDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(KW_DIRECTIVE) {
		return nil, errors.New("Expecting 'directive' keyword for directive definition")
	}

	l.get()

	if !l.tokenEquals(AT.string()) {
		return nil, errors.New("Expecting '@' for directive definition")
	}

	l.get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	argsDef, _ := parseArgumentsDefinition(l)

	if !l.tokenEquals(KW_ON) {
		return nil, errors.New("Expecting 'on' keyworkd for directive definition")
	}

	l.get()

	dls, err := parseDirectiveLocations(l)

	if err != nil {
		return nil, err
	}

	df := &DirectiveDefinition{}

	df.Description = desc
	df.Name = *name
	df.ArgumentsDefinition = argsDef
	df.DirectiveLocations = *dls
	df.Loc = Location{locStart, l.prevLocation().End, l.source}

	return df, nil
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocations
func parseDirectiveLocations(l *lexer) (*DirectiveLocations, error) {
	dls := &DirectiveLocations{}

	if l.tokenEquals(PIPE.string()) {
		l.get()
	}

	dl, err := parseDirectiveLocation(l)

	if err != nil {
		return nil, err
	}

	*dls = append(*dls, *dl)

	for l.tokenEquals(PIPE.string()) {
		l.get()

		dl, err := parseDirectiveLocation(l)

		if err != nil {
			return nil, err
		}

		*dls = append(*dls, *dl)
	}

	return dls, nil
}

// https://graphql.github.io/graphql-spec/draft/#TypeExtension
func parseTypeExtension(l *lexer) (TypeExtension, error) {
	ste, _ := parseScalarTypeExtension(l)

	if ste != nil {
		return ste, nil
	}

	ote, _ := parseObjectTypeExtension(l)

	if ote != nil {
		return ote, nil
	}

	ite, _ := parseInterfaceTypeExtension(l)

	if ite != nil {
		return ite, nil
	}

	ute, _ := parseUnionTypeExtension(l)

	if ute != nil {
		return ute, nil
	}

	ete, _ := parseEnumTypeExtension(l)

	if ete != nil {
		return ete, nil
	}

	iote, _ := parseInputObjectTypeExtension(l)

	if iote != nil {
		return iote, nil
	}

	return nil, errors.New("Expecting type extension")
}

// https://graphql.github.io/graphql-spec/draft/#ScalarTypeExtension
func parseScalarTypeExtension(l *lexer) (*ScalarTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(KW_EXTEND, KW_SCALAR) {
		return nil, errors.New("Expecting 'extend scalar' keywords for scalar type extension")
	}

	l.get()
	l.get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, err := parseDirectives(l)

	if err != nil {
		return nil, err
	}

	ste := &ScalarTypeExtension{}

	ste.Name = *name
	ste.Directives = *dirs
	ste.Loc = Location{locStart, l.prevLocation().End, l.source}

	return ste, nil
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeExtension
func parseObjectTypeExtension(l *lexer) (*ObjectTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(KW_EXTEND, KW_TYPE) {
		return nil, errors.New("Expecting 'extend type' keywords for object type extension")
	}

	l.get()
	l.get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	ii, _ := parseImplementsInterfaces(l)

	dirs, _ := parseDirectives(l)

	fds, _ := parseFieldsDefinition(l)

	if ii == nil && dirs == nil && fds == nil {
		return nil, errors.New("Expecting at least one of 'implements interface', 'directives', 'fields definition' for object type extension")
	}

	ote := &ObjectTypeExtension{}

	ote.Name = *name
	ote.ImplementsInterfaces = ii
	ote.Directives = dirs
	ote.FieldsDefinition = fds
	ote.Loc = Location{locStart, l.prevLocation().End, l.source}

	return ote, nil
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeExtension
func parseInterfaceTypeExtension(l *lexer) (*InterfaceTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(KW_EXTEND, KW_INTERFACE) {
		return nil, errors.New("Expecting 'extend interface' keywords for interface type extension")
	}

	l.get()
	l.get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	fds, _ := parseFieldsDefinition(l)

	if dirs == nil && fds == nil {
		return nil, errors.New("Expecting at least one of 'directives', 'fields definition' for interface type extension")
	}

	ite := &InterfaceTypeExtension{}

	ite.Name = *name
	ite.Directives = dirs
	ite.FieldsDefinition = fds
	ite.Loc = Location{locStart, l.prevLocation().End, l.source}

	return ite, nil
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeExtension
func parseUnionTypeExtension(l *lexer) (*UnionTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(KW_EXTEND, UNION) {
		return nil, errors.New("Expecting 'extend union' keywords for union type extension")
	}

	l.get()
	l.get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	umt, _ := parseUnionMemberTypes(l)

	if dirs == nil && umt == nil {
		return nil, errors.New("Expecting at  least one of 'directives', 'union member types' for union type extension")
	}

	ute := &UnionTypeExtension{}

	ute.Name = *name
	ute.Directives = dirs
	ute.UnionMemberTypes = umt
	ute.Loc = Location{locStart, l.prevLocation().End, l.source}

	return ute, nil
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeExtension
func parseEnumTypeExtension(l *lexer) (*EnumTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(KW_EXTEND, ENUM) {
		return nil, errors.New("Expecting 'extend enum' keywords for enum type extension")
	}

	l.get()
	l.get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	evd, _ := parseEnumValuesDefinition(l)

	if dirs == nil && evd == nil {
		return nil, errors.New("Expecting at least one of 'directives', 'enum values definition' for enum type extension")
	}

	ete := &EnumTypeExtension{}

	ete.Name = *name
	ete.Directives = dirs
	ete.EnumValuesDefinition = evd
	ete.Loc = Location{locStart, l.prevLocation().End, l.source}

	return ete, nil
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocation
func parseDirectiveLocation(l *lexer) (*DirectiveLocation, error) {
	edl, err := parseExecutableDirectiveLocation(l)

	if edl != nil {
		return (*DirectiveLocation)(edl), nil
	}

	tsdl, err := parseTypeSystemDirectiveLocation(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting a directive location")
	}

	return (*DirectiveLocation)(tsdl), nil
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDirectiveLocation
func parseExecutableDirectiveLocation(l *lexer) (*ExecutableDirectiveLocation, error) {
	tok := l.current()

	for i := range executableDirectiveLocations {
		if string(executableDirectiveLocations[i]) == tok.Value {
			l.get()

			edl := executableDirectiveLocations[i]

			return &edl, nil
		}
	}

	return nil, errors.New("Expecting executable directive location")
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDirectiveLocation
func parseTypeSystemDirectiveLocation(l *lexer) (*TypeSystemDirectiveLocation, error) {
	tok := l.current()

	for i := range typeSystemDirectiveLocations {
		if string(typeSystemDirectiveLocations[i]) == tok.Value {
			l.get()

			tsdl := typeSystemDirectiveLocations[i]

			return &tsdl, nil
		}
	}

	return nil, errors.New("Expecting type systen directive location")
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemExtension
func parseTypeSystemExtension(l *lexer) (TypeSystemExtension, error) {
	se, err := parseSchemaExtension(l)

	if se != nil {
		return se, nil
	}

	te, err := parseTypeExtension(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting type system extension")
	}

	return te, nil
}

// https://graphql.github.io/graphql-spec/draft/#SchemaExtension
func parseSchemaExtension(l *lexer) (*SchemaExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(KW_EXTEND, KW_SCHEMA) {
		return nil, errors.New("Expecting 'extend schema' keywords for schema extension")
	}

	l.get()
	l.get()

	dirs, _ := parseDirectives(l)

	if !l.tokenEquals(BRACE_L.string()) {
		return nil, errors.New("Expecting '{' for schema extension")
	}

	l.get()

	rotds, _ := parseRootOperationTypeDefinitions(l)

	if !l.tokenEquals(BRACE_R.string()) {
		return nil, errors.New("Expecting '}' for schema extension")
	}

	l.get()

	if dirs == nil && rotds == nil {
		return nil, errors.New("Expecting directives or root operation type definitions for schema extension")
	}

	se := &SchemaExtension{}

	se.Directives = dirs
	se.RootOperationTypeDefinitions = rotds
	se.Loc = Location{locStart, l.prevLocation().End, l.source}

	return se, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeExtension
func parseInputObjectTypeExtension(l *lexer) (*InputObjectTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(KW_EXTEND, KW_INPUT) {
		return nil, errors.New("Expecting 'extend' keyword for input object type extension")
	}

	l.get()
	l.get()

	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	idfs, _ := parseInputFieldsDefinition(l)

	if dirs == nil && idfs == nil {
		return nil, errors.New("Expecting at lease one of 'directives', 'input fields definition' fo input object type extension")
	}

	iote := &InputObjectTypeExtension{}

	iote.Name = *name
	iote.Directives = dirs
	iote.InputFieldsDefinition = idfs
	iote.Loc = Location{locStart, l.prevLocation().End, l.source}

	return iote, nil
}

// https://graphql.github.io/graphql-spec/draft/#OperationDefinition
func parseOperationDefinition(l *lexer) (*OperationDefinition, error) {
	locStart := l.location().Start

	// Shorthand query
	// https://graphql.github.io/graphql-spec/draft/#sec-Language.Operations.Query-shorthand
	if l.tokenEquals(BRACE_L.string()) {

		shorthandQuery, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		opDef := &OperationDefinition{}

		opDef.OperationType = KW_QUERY
		opDef.SelectionSet = *shorthandQuery
		opDef.Loc = Location{locStart, l.prevLocation().End, l.source}

		return opDef, nil
	} else if !l.tokenEquals(KW_QUERY) &&
		!l.tokenEquals(KW_MUTATION) &&
		!l.tokenEquals(KW_SUBSCRIPTION) {
		return nil, errors.New("Expecting one of 'query', 'mutation', 'subscription' for operation definition")
	} else {
		tok := l.get()

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
		opDefinition.Loc = Location{locStart, tok.End, l.source}

		return opDefinition, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentDefinition
func parseFragmentDefinition(l *lexer) (*FragmentDefinition, error) {
	locStart := l.location().Start

	if !l.tokenEquals(KW_FRAGMENT) {
		return nil, errors.New("Expecting fragment keyword")
	} else {
		l.get()

		name, err := parseFragmentName(l)

		if err != nil {
			return nil, err
		}

		if name.Value == KW_ON {
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
		fragDef.Loc = Location{locStart, l.prevLocation().End, l.source}

		return fragDef, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#Name
func parseName(l *lexer) (*Name, error) {
	tok := l.current()

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

	l.get()

	name := &Name{}

	// Populate the Name struct.
	name.Value = tok.Value
	name.Loc.Start = tok.Start
	name.Loc.End = tok.End
	name.Loc.Source = l.source

	// Return the AST Name object.
	return name, nil
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinition
func parseVariableDefinitions(l *lexer) (*VariableDefinitions, error) {
	if !l.tokenEquals(PAREN_L.string()) {
		return nil, errors.New("Expecting '(' opener for variable definitions")
	} else {
		l.get()

		varDefs := &VariableDefinitions{}

		for !l.tokenEquals(PAREN_R.string()) {
			varDef, err := parseVariableDefinition(l)

			if err != nil {
				break
			}

			*varDefs = append(*varDefs, *varDef)
		}

		l.get()

		return varDefs, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinition
func parseVariableDefinition(l *lexer) (*VariableDefinition, error) {
	locStart := l.location().Start

	_var, err := parseVariable(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(COLON.string()) {
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
	varDef.Loc = Location{locStart, locEnd, l.source}

	return varDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#Type
func parseType(l *lexer) (Type, error) {
	var _type Type

	_type, err := parseNamedType(l)

	if _type != nil {
		return _type, nil
	}

	_type, err = parseListType(l)

	if _type != nil {
		return _type, nil
	}

	_type, err = parseNonNullType(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting a type")
	}

	return _type, nil
}

// https://graphql.github.io/graphql-spec/draft/#ListType
func parseListType(l *lexer) (*ListType, error) {
	locStart := l.location().Start

	if !l.tokenEquals(BRACKET_L.string()) {
		return nil, errors.New("Expecting '[' for list type")
	}

	l.get()

	_type, err := parseType(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(BRACKET_R.string()) {
		return nil, errors.New("Expecting ']' for list type")
	}

	l.get()

	listType := &ListType{}

	listType.OfType = _type
	listType.Loc = Location{locStart, l.prevLocation().End, l.source}

	return listType, nil
}

// https://graphql.github.io/graphql-spec/draft/#NonNullType
func parseNonNullType(l *lexer) (*NonNullType, error) {
	locStart := l.location().Start

	var _type Type

	_type, err := parseNamedType(l)

	if err != nil {
		_type, err = parseListType(l)

		if err != nil {
			return nil, err
		}
	}

	if !l.tokenEquals(BANG.string()) {
		return nil, errors.New("Expecting '!' at the end of a non null type")
	}

	l.get()

	nonNull := &NonNullType{}

	nonNull.OfType = _type
	nonNull.Loc = Location{locStart, l.prevLocation().End, l.source}

	return nonNull, nil
}

// https://graphql.github.io/graphql-spec/draft/#Directives
func parseDirectives(l *lexer) (*Directives, error) {
	dirs := &Directives{}

	for l.tokenEquals(AT.string()) {
		dir, err := parseDirective(l)

		if err != nil {
			return nil, err
		}

		*dirs = append(*dirs, *dir)
	}

	if len(*dirs) == 0 {
		return nil, errors.New("Expecting at least one directive")
	}

	return dirs, nil
}

// https://graphql.github.io/graphql-spec/draft/#Directive
func parseDirective(l *lexer) (*Directive, error) {
	locStart := l.location().Start

	if !l.tokenEquals(AT.string()) {
		return nil, errors.New("Expecting '@' for directive")
	} else {
		l.get()

		name, err := parseName(l)

		if err != nil {
			return nil, err
		}

		args, _ := parseArguments(l)

		dir := &Directive{}

		dir.Name = *name
		dir.Arguments = args
		dir.Loc = Location{locStart, l.prevLocation().End, l.source}

		return dir, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#SelectionSet
func parseSelectionSet(l *lexer) (*SelectionSet, error) {
	if !l.tokenEquals(BRACE_L.string()) {
		return nil, errors.New("Expecting '{' for selection set")
	} else {
		l.get()

		selSet := &SelectionSet{}

		for !l.tokenEquals(BRACE_R.string()) {
			sel, err := parseSelection(l)

			if err != nil {
				return nil, err
			}

			*selSet = append(*selSet, sel)
		}

		l.get()

		return selSet, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#Selection
func parseSelection(l *lexer) (Selection, error) {
	var sel Selection

	sel, _ = parseField(l)

	if sel != nil {
		return sel, nil
	}

	sel, _ = parseFragmentSpread(l)

	if sel != nil {
		return sel, nil
	}

	sel, err := parseInlineFragment(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting a selection")
	}

	return sel, nil
}

// https://graphql.github.io/graphql-spec/draft/#Variable
func parseVariable(l *lexer) (*Variable, error) {
	locStart := l.location().Start

	if !l.tokenEquals(DOLLAR.string()) {
		return nil, errors.New("Expecting '$' for varible")
	} else {
		name, err := parseName(l)

		if err != nil {
			return nil, err
		}

		_var := &Variable{}

		_var.Name = *name
		_var.Loc = Location{locStart, name.Location().End, l.source}

		return _var, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#DefaultValue
func parseDefaultValue(l *lexer) (*DefaultValue, error) {
	locStart := l.location().Start

	if !l.tokenEquals(EQUALS.string()) {
		return nil, errors.New("Expecting '=' for default value")
	} else {
		val, err := parseValue(l)

		if err != nil {
			return nil, err
		}

		dVal := &DefaultValue{}

		dVal.Value = val
		dVal.Loc = Location{locStart, l.prevLocation().End, l.source}

		return dVal, nil
	}
}

// ! need to check variable type in order to parse its value
// https://graphql.github.io/graphql-spec/draft/#Value
func parseValue(l *lexer) (Value, error) {
	// need to parse dynamic variables
	//_var, _ := parseVariable(l)

	var val Value

	val, _ = parseIntValue(l)

	if val != nil {
		return val, nil
	}

	val, _ = parseFloatValue(l)

	if val != nil {
		return val, nil
	}

	val, _ = parseStringValue(l)

	if val != nil {
		return val, nil
	}

	val, _ = parseBooleanValue(l)

	if val != nil {
		return val, nil
	}

	val, _ = parseNullValue(l)

	if val != nil {
		return val, nil
	}

	val, _ = parseEnumValue(l)

	if val != nil {
		return val, nil
	}

	val, _ = parseListValue(l)

	if val != nil {
		return val, nil
	}

	val, err := parseObjectValue(l)

	if err != nil {
		return val, errors.Wrap(err, "Expecting a value")
	}

	return val, nil
}

// https://graphql.github.io/graphql-spec/draft/#Arguments
func parseArguments(l *lexer) (*Arguments, error) {
	if !l.tokenEquals(PAREN_L.string()) {
		return nil, errors.New("Expecting '(' for arguments")
	} else {
		l.get()

		args := &Arguments{}

		for !l.tokenEquals(PAREN_R.string()) {
			arg, err := parseArgument(l)

			if err != nil {
				return nil, err
			}

			*args = append(*args, *arg)

		}

		l.get()

		if len(*args) == 0 {
			return nil, errors.New("Expecting at least one argument")
		}

		return args, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#Argument
func parseArgument(l *lexer) (*Argument, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(COLON.string()) {
		return nil, errors.New("Expecting colon after argument name")
	}

	l.get()

	val, err := parseValue(l)

	if err != nil {
		return nil, err
	}

	arg := &Argument{}

	arg.Name = *name
	arg.Value = val
	arg.Loc = Location{name.Location().Start, l.prevLocation().End, l.source}

	return arg, nil
}

// https://graphql.github.io/graphql-spec/draft/#Field
func parseField(l *lexer) (*Field, error) {
	locStart := l.location().Start

	alias, err := parseName(l)

	if err != nil {
		return nil, err
	}

	name := &Name{}

	if l.tokenEquals(COLON.string()) {
		l.get()

		name, err = parseName(l)

		if err != nil {
			return nil, err
		}
	} else {
		*name = *alias

		alias = nil
	}

	args, _ := parseArguments(l)

	dirs, _ := parseDirectives(l)

	selSet, _ := parseSelectionSet(l)

	field := &Field{}

	field.Alias = (*Alias)(alias)
	field.Name = *name
	field.Arguments = args
	field.Directives = dirs
	field.SelectionSet = selSet
	field.Loc = Location{locStart, l.prevLocation().End, l.source}

	return field, nil
}

// https://graphql.github.io/graphql-spec/draft/#FragmentSpread
func parseFragmentSpread(l *lexer) (*FragmentSpread, error) {
	locStart := l.location().Start

	if !l.tokenEquals(SPREAD.string()) {
		return nil, errors.New("Expecting '...' operator for a fragment spread")
	} else {
		l.get()

		fname, err := parseFragmentName(l)

		if err != nil {
			return nil, err
		}

		directives, _ := parseDirectives(l)

		spread := &FragmentSpread{}

		spread.FragmentName = *fname
		spread.Directives = directives
		spread.Loc = Location{locStart, l.prevLocation().End, l.source}

		return spread, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#InlineFragment
func parseInlineFragment(l *lexer) (*InlineFragment, error) {
	locStart := l.location().Start

	if !l.tokenEquals(SPREAD.string()) {
		return nil, errors.New("Expecting '...' for an inline fragment")
	} else {
		l.get()

		typeCon, _ := parseTypeCondition(l)

		directives, _ := parseDirectives(l)

		selSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		inlineFrag := &InlineFragment{}

		inlineFrag.TypeCondition = typeCon
		inlineFrag.Directives = directives
		inlineFrag.SelectionSet = *selSet
		inlineFrag.Loc = Location{locStart, l.prevLocation().End, l.source}

		return inlineFrag, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentName
func parseFragmentName(l *lexer) (*FragmentName, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if name.Value == KW_ON {
		return nil, errors.New("Fragment name cannot be 'on'")
	}

	var fragName *FragmentName

	*fragName = FragmentName(*name)
	fragName.Loc = *name.Location()

	return fragName, nil
}

// https://graphql.github.io/graphql-spec/draft/#TypeCondition
func parseTypeCondition(l *lexer) (*TypeCondition, error) {
	locStart := l.location().Start

	if !l.tokenEquals(KW_ON) {
		return nil, errors.New("Expecting 'on' keyword for a type condition")
	} else {
		namedType, err := parseNamedType(l)

		if err != nil {
			return nil, err
		}

		typeCond := &TypeCondition{}

		typeCond.NamedType = *namedType
		typeCond.Loc = Location{locStart, namedType.Location().End, l.source}

		return typeCond, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#NamedType
func parseNamedType(l *lexer) (*NamedType, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	var namedType *NamedType

	*namedType = NamedType(*name)

	return namedType, nil
}

// ! Check numeric values
// https://graphql.github.io/graphql-spec/draft/#IntValue
func parseIntValue(l *lexer) (*IntValue, error) {
	tok := l.current()

	intVal, err := strconv.ParseInt(tok.Value, 10, 64)

	if err != nil {
		return nil, err
	}

	l.get()

	intValP := &IntValue{}

	intValP.Value = intVal
	intValP.Loc = Location{tok.Start, tok.End, l.source}

	return intValP, nil
}

// https://graphql.github.io/graphql-spec/draft/#FloatValue
func parseFloatValue(l *lexer) (*FloatValue, error) {
	tok := l.current()

	floatVal, err := strconv.ParseFloat(tok.Value, 64)

	if err != nil {
		return nil, err
	}

	l.get()

	floatValP := &FloatValue{}

	floatValP.Value = floatVal
	floatValP.Loc = Location{tok.Start, tok.End, l.source}

	return floatValP, nil
}

// ! Have a discussion about this function
// https://graphql.github.io/graphql-spec/draft/#StringValue
func parseStringValue(l *lexer) (*StringValue, error) {
	tok := l.current()

	sv := &StringValue{}

	sv.Value = tok.Value
	sv.Loc = Location{tok.Start, tok.End, l.source}

	return sv, nil
}

// https://graphql.github.io/graphql-spec/draft/#BooleanValue
func parseBooleanValue(l *lexer) (*BooleanValue, error) {
	tok := l.current()

	boolVal, err := strconv.ParseBool(tok.Value)

	if err != nil {
		return nil, err
	}

	l.get()

	boolValP := &BooleanValue{}

	boolValP.Value = boolVal
	boolValP.Loc = Location{tok.Start, tok.End, l.source}

	return boolValP, nil
}

// ! Figure out what to do with a null value
// https://graphql.github.io/graphql-spec/draft/#NullValue
func parseNullValue(l *lexer) (*NullValue, error) {
	tok := l.current()

	if tok.Value != KW_NULL {
		return nil, errors.New("Expecting 'null' keyword")
	} else {
		l.get()

		null := &NullValue{}
		null.Loc = Location{tok.Start, tok.End, l.source}

		return null, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#EnumValue
func parseEnumValue(l *lexer) (*EnumValue, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	switch name.Value {
	case KW_TRUE, KW_FALSE, KW_NULL:
		return nil, errors.New("Enum value cannot be 'true', 'false' or 'null'")
	default:
		enumVal := &EnumValue{}

		enumVal.Name = *name
		enumVal.Loc = Location{name.Location().Start, name.Location().End, l.source}

		return enumVal, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#ListValue
func parseListValue(l *lexer) (*ListValue, error) {
	locStart := l.location().Start

	if !l.tokenEquals(BRACKET_L.string()) {
		return nil, errors.New("Expecting '[' for a list value")
	} else {
		l.get()

		lstVal := &ListValue{}

		for !l.tokenEquals(BRACKET_R.string()) {
			val, err := parseValue(l)

			if err != nil {
				return nil, err
			}

			lstVal.Values = append(lstVal.Values, val)
		}

		lstVal.Loc = Location{locStart, l.prevLocation().End, l.source}

		return lstVal, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectValue
func parseObjectValue(l *lexer) (*ObjectValue, error) {
	locStart := l.location().Start

	if !l.tokenEquals(BRACE_L.string()) {
		return nil, errors.New("Expecting '{' for an object value")
	} else {
		l.get()

		objVal := &ObjectValue{}

		for !l.tokenEquals(BRACE_R.string()) {
			objField, err := parseObjectField(l)

			if err != nil {
				return nil, err
			}

			objVal.Values = append(objVal.Values, *objField)
		}

		objVal.Loc = Location{locStart, l.prevLocation().End, l.source}

		return objVal, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectField
func parseObjectField(l *lexer) (*ObjectField, error) {
	name, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(COLON.string()) {
		return nil, errors.New("Expecting color after object field name")
	}

	l.get()

	val, err := parseValue(l)

	if err != nil {
		return nil, err
	}

	objField := &ObjectField{}

	objField.Name = *name
	objField.Value = val
	objField.Loc = Location{name.Location().Start, l.prevLocation().End, l.source}

	return objField, nil
}
