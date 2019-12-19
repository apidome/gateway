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

func Parse(doc string) (*document, error) {
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
func parseDocument(l *lexer) (*document, error) {
	def, err := parseDefinitions(l)

	if err != nil {
		return nil, err
	}

	doc := &document{}

	doc.Definitions = *def

	return doc, nil
}

// https://graphql.github.io/graphql-spec/draft/#Definition
func parseDefinitions(l *lexer) (*definitions, error) {
	defs := &definitions{}

	for !l.tokenEquals(tokEOF.string()) {
		var def definition

		def, err := parseDefinition(l)

		if err != nil {
			return nil, err
		}

		*defs = append(*defs, def)
	}

	if len(*defs) == 0 {
		return nil, errors.New("No definitions found in document")
	}

	return defs, nil
}

// https://graphql.github.io/graphql-spec/draft/#Definition
func parseDefinition(l *lexer) (definition, error) {
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
func parseExecutableDefinition(l *lexer) (executableDefinition, error) {
	var execDef executableDefinition

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
func parseTypeSystemDefinition(l *lexer) (typeSystemDefinition, error) {
	var def typeSystemDefinition

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
func parseSchemaDefinition(l *lexer) (*schemaDefinition, error) {
	locStart := l.current().start

	if !l.tokenEquals(kwSchema) {
		return nil, errors.New("Missing 'schema' keyword for a schema definition")
	}

	l.get()

	dirs, _ := parseDirectives(l)

	if !l.tokenEquals(tokBraceL.string()) {
		return nil, errors.New("Missing '{' for a schema definition")
	}

	l.get()

	rOtd, err := parseRootOperationTypeDefinitions(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(tokBraceR.string()) {
		return nil, errors.New("Missing '}' for schema definition")
	}

	locEnd := l.current().end

	l.get()

	schDef := &schemaDefinition{}

	schDef.Directives = dirs
	schDef.RootOperationTypeDefinitions = *rOtd
	schDef.Loc = location{locStart, locEnd, l.source}

	return schDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#RootOperationTypeDefinition
func parseRootOperationTypeDefinitions(l *lexer) (*rootOperationTypeDefinitions, error) {
	rotds := &rootOperationTypeDefinitions{}

	for !l.tokenEquals(tokBraceR.string()) {
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
func parseRootOperationTypeDefinition(l *lexer) (*rootOperationTypeDefinition, error) {
	locStart := l.current().start

	opType, err := parseOperationType(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(tokColon.string()) {
		return nil, errors.New("Expecting ':' after operation type")
	}

	l.get()

	namedTyp, err := parseNamedType(l)

	if err != nil {
		return nil, err
	}

	rotd := &rootOperationTypeDefinition{}

	rotd.OperationType = *opType
	rotd.NamedType = *namedTyp
	rotd.Loc = location{locStart, namedTyp.Location().End, l.source}

	return rotd, nil
}

// https://graphql.github.io/graphql-spec/draft/#OperationType
func parseOperationType(l *lexer) (*operationType, error) {
	tok := l.current()

	if tok.value != string(operationMutation) &&
		tok.value != string(operationQuery) &&
		tok.value != string(operationSubscription) {
		return nil,
			errors.New("Expecting 'query', 'mutation' or 'subscription' as operation type")
	}

	opType := new(operationType)

	*opType = (operationType)(tok.value)

	return opType, nil
}

// https://graphql.github.io/graphql-spec/draft/#TypeDefinition
func parseTypeDefinition(l *lexer) (typeDefinition, error) {
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
func parseScalarTypeDefinition(l *lexer) (*scalarTypeDefinition, error) {
	desc, _ := parseDescription(l)

	if !l.tokenEquals(tsdlScalar) {
		return nil, errors.New("Missing 'scalar' keyword for scalar type definition")
	}

	tok := l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	scalarTd := &scalarTypeDefinition{}

	scalarTd.Description = desc
	scalarTd.Name = *nam
	scalarTd.Directives = dirs
	scalarTd.Loc = location{tok.start, l.prevLocation().End, l.source}

	return scalarTd, nil
}

// https://graphql.github.io/graphql-spec/draft/#Description
func parseDescription(l *lexer) (*description, error) {
	strVal, err := parseStringValue(l)

	if err != nil {
		return nil, err
	}

	return (*description)(strVal), nil
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeDefinition
func parseObjectTypeDefinition(l *lexer) (*objectTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(kwType) {
		return nil, errors.New("Expecting 'type' keyword for object type definition")
	}

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	ii, _ := parseImplementsInterfaces(l)

	dirs, _ := parseDirectives(l)

	fd, _ := parseFieldsDefinition(l)

	objTd := &objectTypeDefinition{}

	objTd.Description = desc
	objTd.Directives = dirs
	objTd.FieldsDefinition = fd
	objTd.ImplementsInterfaces = ii
	objTd.Name = *nam
	objTd.Loc = location{locStart, l.prevLocation().End, l.source}

	return objTd, nil
}

// https://graphql.github.io/graphql-spec/draft/#ImplementsInterfaces
func parseImplementsInterfaces(l *lexer) (*implementsInterfaces, error) {
	if !l.tokenEquals(kwImplements) {
		return nil, errors.New("Expecting 'implements' keyword")
	}

	if l.tokenEquals(tokAmp.string()) {
		l.get()
	}

	nt, err := parseNamedType(l)

	if err != nil {
		return nil, err
	}

	ii := &implementsInterfaces{}

	(*ii) = append(*ii, *nt)

	for l.tokenEquals(tokAmp.string()) {
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
func parseFieldsDefinition(l *lexer) (*fieldsDefinition, error) {
	if l.tokenEquals(tokBraceL.string()) {
		return nil, errors.New("Expecting '{' for fields definition")
	}

	l.get()

	fds := &fieldsDefinition{}

	for !l.tokenEquals(tokBraceR.string()) {
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
func parseFieldDefinition(l *lexer) (*fieldDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	argsDef, _ := parseArgumentsDefinition(l)

	if !l.tokenEquals(tokColon.string()) {
		return nil, errors.New("Expecting ':' for field definition")
	}

	l.get()

	_typ, err := parseType(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	fd := &fieldDefinition{}

	fd.Description = desc
	fd.Name = *nam
	fd.ArgumentsDefinition = argsDef
	fd.Type = _typ
	fd.Directives = dirs
	fd.Loc = location{locStart, l.prevLocation().End, l.source}

	return fd, nil
}

// https://graphql.github.io/graphql-spec/draft/#ArgumentsDefinition
func parseArgumentsDefinition(l *lexer) (*argumentsDefinition, error) {
	if !l.tokenEquals(tokParenL.string()) {
		return nil, errors.New("Expecting '(' for arguments definition")
	}

	l.get()

	argsDef := &argumentsDefinition{}

	for !l.tokenEquals(tokParenR.string()) {
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
func parseInputValueDefinition(l *lexer) (*inputValueDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(tokColon.string()) {
		return nil, errors.New("Expecting ':' for input value definition")
	}

	l.get()

	_typ, err := parseType(l)

	if err != nil {
		return nil, err
	}

	defVal, _ := parseDefaultValue(l)

	dirs, _ := parseDirectives(l)

	ivDef := &inputValueDefinition{}

	ivDef.Description = desc
	ivDef.Name = *nam
	ivDef.Type = _typ
	ivDef.DefaultValue = defVal
	ivDef.Directives = dirs
	ivDef.Loc = location{locStart, l.prevLocation().End, l.source}

	return ivDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeDefinition
func parseInterfaceTypeDefinition(l *lexer) (*interfaceTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(tsdlInterface) {
		return nil, errors.New("Expecting 'interface' keyword for interface type definition")
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	fds, _ := parseFieldsDefinition(l)

	itd := &interfaceTypeDefinition{}

	itd.Description = desc
	itd.Directives = dirs
	itd.FieldsDefinition = fds
	itd.Name = *nam
	itd.Loc = location{locStart, l.prevLocation().End, l.source}

	return itd, nil
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeDefinition
func parseUnionTypeDefinition(l *lexer) (*unionTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(tsdlUnion) {
		return nil, errors.New("Expecting 'union' keyowrd for union type definition")
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	umt, _ := parseUnionMemberTypes(l)

	utd := &unionTypeDefinition{}

	utd.Description = desc
	utd.Name = *nam
	utd.Directives = dirs
	utd.UnionMemberTypes = umt
	utd.Loc = location{locStart, l.prevLocation().End, l.source}

	return utd, nil
}

// https://graphql.github.io/graphql-spec/draft/#UnionMemberTypes
func parseUnionMemberTypes(l *lexer) (*unionMemberTypes, error) {
	if !l.tokenEquals(tokEquals.string()) {
		return nil, errors.New("Expecting '=' for union member types")
	}

	l.get()

	if l.tokenEquals(tokPipe.string()) {
		l.get()
	}

	nt, err := parseNamedType(l)

	if err != nil {
		return nil, err
	}

	umt := &unionMemberTypes{}

	*umt = append(*umt, *nt)

	for l.tokenEquals(tokPipe.string()) {
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
func parseEnumTypeDefinition(l *lexer) (*enumTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(tsdlEnum) {
		return nil, errors.New("Expecting 'enum' keyword for enum type definition")
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	evd, _ := parseEnumValuesDefinition(l)

	etd := &enumTypeDefinition{}

	etd.Description = desc
	etd.Name = *nam
	etd.Directives = dirs
	etd.EnumValuesDefinition = evd
	etd.Loc = location{locStart, l.prevLocation().End, l.source}

	return etd, nil
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinition
func parseEnumValuesDefinition(l *lexer) (*enumValuesDefinition, error) {
	if !l.tokenEquals(tokBraceL.string()) {
		return nil, errors.New("Expecting '{' for enum values definition")
	}

	l.get()

	evds := &enumValuesDefinition{}

	for !l.tokenEquals(tokBraceR.string()) {
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
func parseEnumValueDefinition(l *lexer) (*enumValueDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	ev, err := parseEnumValue(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	evd := &enumValueDefinition{}

	evd.Description = desc
	evd.EnumValue = *ev
	evd.Directives = dirs
	ev.Loc = location{locStart, l.prevLocation().End, l.source}

	return evd, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeDefinition
func parseInputObjectTypeDefinition(l *lexer) (*inputObjectTypeDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(kwInput) {
		return nil, errors.New("Expecting 'input' keyword for input object type definition")
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	ifds, _ := parseInputFieldsDefinition(l)

	iotd := &inputObjectTypeDefinition{}

	iotd.Description = desc
	iotd.Directives = dirs
	iotd.Name = *nam
	iotd.InputFieldsDefinition = ifds
	iotd.Loc = location{locStart, l.prevLocation().End, l.source}

	return iotd, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputFieldsDefinition
func parseInputFieldsDefinition(l *lexer) (*inputFieldsDefinition, error) {
	if !l.tokenEquals(tokBraceL.string()) {
		return nil, errors.New("Expecting '{' for input fields definition")
	}

	l.get()

	ifds := &inputFieldsDefinition{}

	for !l.tokenEquals(tokBraceR.string()) {
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
func parseDirectiveDefinition(l *lexer) (*directiveDefinition, error) {
	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(kwDirective) {
		return nil, errors.New("Expecting 'directive' keyword for directive definition")
	}

	l.get()

	if !l.tokenEquals(tokAt.string()) {
		return nil, errors.New("Expecting '@' for directive definition")
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	argsDef, _ := parseArgumentsDefinition(l)

	if !l.tokenEquals(kwOn) {
		return nil, errors.New("Expecting 'on' keyworkd for directive definition")
	}

	l.get()

	dls, err := parseDirectiveLocations(l)

	if err != nil {
		return nil, err
	}

	df := &directiveDefinition{}

	df.Description = desc
	df.Name = *nam
	df.ArgumentsDefinition = argsDef
	df.DirectiveLocations = *dls
	df.Loc = location{locStart, l.prevLocation().End, l.source}

	return df, nil
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocations
func parseDirectiveLocations(l *lexer) (*directiveLocations, error) {
	dls := &directiveLocations{}

	if l.tokenEquals(tokPipe.string()) {
		l.get()
	}

	dl, err := parseDirectiveLocation(l)

	if err != nil {
		return nil, err
	}

	*dls = append(*dls, *dl)

	for l.tokenEquals(tokPipe.string()) {
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
func parseTypeExtension(l *lexer) (typeExtension, error) {
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
func parseScalarTypeExtension(l *lexer) (*scalarTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwScalar) {
		return nil, errors.New("Expecting 'extend scalar' keywords for scalar type extension")
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, err := parseDirectives(l)

	if err != nil {
		return nil, err
	}

	ste := &scalarTypeExtension{}

	ste.Name = *nam
	ste.Directives = *dirs
	ste.Loc = location{locStart, l.prevLocation().End, l.source}

	return ste, nil
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeExtension
func parseObjectTypeExtension(l *lexer) (*objectTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwType) {
		return nil, errors.New("Expecting 'extend type' keywords for object type extension")
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	ii, _ := parseImplementsInterfaces(l)

	dirs, _ := parseDirectives(l)

	fds, _ := parseFieldsDefinition(l)

	if ii == nil && dirs == nil && fds == nil {
		return nil, errors.New("Expecting at least one of 'implements interface', 'directives', 'fields definition' for object type extension")
	}

	ote := &objectTypeExtension{}

	ote.Name = *nam
	ote.ImplementsInterfaces = ii
	ote.Directives = dirs
	ote.FieldsDefinition = fds
	ote.Loc = location{locStart, l.prevLocation().End, l.source}

	return ote, nil
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeExtension
func parseInterfaceTypeExtension(l *lexer) (*interfaceTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwInterface) {
		return nil, errors.New("Expecting 'extend interface' keywords for interface type extension")
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	fds, _ := parseFieldsDefinition(l)

	if dirs == nil && fds == nil {
		return nil, errors.New("Expecting at least one of 'directives', 'fields definition' for interface type extension")
	}

	ite := &interfaceTypeExtension{}

	ite.Name = *nam
	ite.Directives = dirs
	ite.FieldsDefinition = fds
	ite.Loc = location{locStart, l.prevLocation().End, l.source}

	return ite, nil
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeExtension
func parseUnionTypeExtension(l *lexer) (*unionTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, tsdlUnion) {
		return nil, errors.New("Expecting 'extend union' keywords for union type extension")
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	umt, _ := parseUnionMemberTypes(l)

	if dirs == nil && umt == nil {
		return nil, errors.New("Expecting at  least one of 'directives', 'union member types' for union type extension")
	}

	ute := &unionTypeExtension{}

	ute.Name = *nam
	ute.Directives = dirs
	ute.UnionMemberTypes = umt
	ute.Loc = location{locStart, l.prevLocation().End, l.source}

	return ute, nil
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeExtension
func parseEnumTypeExtension(l *lexer) (*enumTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, tsdlEnum) {
		return nil, errors.New("Expecting 'extend enum' keywords for enum type extension")
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	evd, _ := parseEnumValuesDefinition(l)

	if dirs == nil && evd == nil {
		return nil, errors.New("Expecting at least one of 'directives', 'enum values definition' for enum type extension")
	}

	ete := &enumTypeExtension{}

	ete.Name = *nam
	ete.Directives = dirs
	ete.EnumValuesDefinition = evd
	ete.Loc = location{locStart, l.prevLocation().End, l.source}

	return ete, nil
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocation
func parseDirectiveLocation(l *lexer) (*directiveLocation, error) {
	edl, err := parseExecutableDirectiveLocation(l)

	if edl != nil {
		return (*directiveLocation)(edl), nil
	}

	tsdl, err := parseTypeSystemDirectiveLocation(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting a directive location")
	}

	return (*directiveLocation)(tsdl), nil
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDirectiveLocation
func parseExecutableDirectiveLocation(l *lexer) (*executableDirectiveLocation, error) {
	tok := l.current()

	for i := range executableDirectiveLocations {
		if string(executableDirectiveLocations[i]) == tok.value {
			l.get()

			edl := executableDirectiveLocations[i]

			return &edl, nil
		}
	}

	return nil, errors.New("Expecting executable directive location")
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDirectiveLocation
func parseTypeSystemDirectiveLocation(l *lexer) (*typeSystemDirectiveLocation, error) {
	tok := l.current()

	for i := range typeSystemDirectiveLocations {
		if string(typeSystemDirectiveLocations[i]) == tok.value {
			l.get()

			tsdl := typeSystemDirectiveLocations[i]

			return &tsdl, nil
		}
	}

	return nil, errors.New("Expecting type systen directive location")
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemExtension
func parseTypeSystemExtension(l *lexer) (typeSystemExtension, error) {
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
func parseSchemaExtension(l *lexer) (*schemaExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwSchema) {
		return nil, errors.New("Expecting 'extend schema' keywords for schema extension")
	}

	l.get()
	l.get()

	dirs, _ := parseDirectives(l)

	if !l.tokenEquals(tokBraceL.string()) {
		return nil, errors.New("Expecting '{' for schema extension")
	}

	l.get()

	rotds, _ := parseRootOperationTypeDefinitions(l)

	if !l.tokenEquals(tokBraceR.string()) {
		return nil, errors.New("Expecting '}' for schema extension")
	}

	l.get()

	if dirs == nil && rotds == nil {
		return nil, errors.New("Expecting directives or root operation type definitions for schema extension")
	}

	se := &schemaExtension{}

	se.Directives = dirs
	se.RootOperationTypeDefinitions = rotds
	se.Loc = location{locStart, l.prevLocation().End, l.source}

	return se, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeExtension
func parseInputObjectTypeExtension(l *lexer) (*inputObjectTypeExtension, error) {
	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwInput) {
		return nil, errors.New("Expecting 'extend' keyword for input object type extension")
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	dirs, _ := parseDirectives(l)

	idfs, _ := parseInputFieldsDefinition(l)

	if dirs == nil && idfs == nil {
		return nil, errors.New("Expecting at lease one of 'directives', 'input fields definition' fo input object type extension")
	}

	iote := &inputObjectTypeExtension{}

	iote.Name = *nam
	iote.Directives = dirs
	iote.InputFieldsDefinition = idfs
	iote.Loc = location{locStart, l.prevLocation().End, l.source}

	return iote, nil
}

// https://graphql.github.io/graphql-spec/draft/#OperationDefinition
func parseOperationDefinition(l *lexer) (*operationDefinition, error) {
	locStart := l.location().Start

	// Shorthand query
	// https://graphql.github.io/graphql-spec/draft/#sec-Language.Operations.Query-shorthand
	if l.tokenEquals(tokBraceL.string()) {

		shorthandQuery, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		opDef := &operationDefinition{}

		opDef.OperationType = kwQuery
		opDef.SelectionSet = *shorthandQuery
		opDef.Loc = location{locStart, l.prevLocation().End, l.source}

		return opDef, nil
	} else if !l.tokenEquals(kwQuery) &&
		!l.tokenEquals(kwMutation) &&
		!l.tokenEquals(kwSubscription) {
		return nil, errors.New("Expecting one of 'query', 'mutation', 'subscription' for operation definition")
	} else {
		tok := l.get()

		opType := tok.value

		nam, _ := parseName(l)

		varDef, _ := parseVariableDefinitions(l)

		directives, _ := parseDirectives(l)

		selSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		opDefinition := &operationDefinition{}

		opDefinition.OperationType = operationType(opType)
		opDefinition.Name = nam
		opDefinition.VariableDefinitions = varDef
		opDefinition.Directives = directives
		opDefinition.SelectionSet = *selSet
		opDefinition.Loc = location{locStart, tok.end, l.source}

		return opDefinition, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentDefinition
func parseFragmentDefinition(l *lexer) (*fragmentDefinition, error) {
	locStart := l.location().Start

	if !l.tokenEquals(kwFragment) {
		return nil, errors.New("Expecting fragment keyword")
	} else {
		l.get()

		nam, err := parseFragmentName(l)

		if err != nil {
			return nil, err
		}

		if nam.Value == kwOn {
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

		fragDef := &fragmentDefinition{}

		fragDef.FragmentName = *nam
		fragDef.TypeCondition = *typeCond
		fragDef.Directives = directives
		fragDef.SelectionSet = *selectionSet
		fragDef.Loc = location{locStart, l.prevLocation().End, l.source}

		return fragDef, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#Name
func parseName(l *lexer) (*name, error) {
	tok := l.current()

	pattern := "^[_A-Za-z][_0-9A-Za-z]*$"

	// If the current token is not a Name, return nil
	if tok.kind != tokName {
		return nil, errors.New("Not a name")
	}

	// Check if the given name matches the regex provided by graphql spec at
	// https://graphql.github.io/graphql-spec/draft/#Name
	match, err := regexp.MatchString(pattern, tok.value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse name: ")
	}

	// If the name does not match the requirements, return an error.
	if !match {
		return nil, errors.New("invalid name - " + tok.value)
	}

	l.get()

	nam := &name{}

	// Populate the Name struct.
	nam.Value = tok.value
	nam.Loc.Start = tok.start
	nam.Loc.End = tok.end
	nam.Loc.Source = l.source

	// Return the AST Name object.
	return nam, nil
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinition
func parseVariableDefinitions(l *lexer) (*variableDefinitions, error) {
	if !l.tokenEquals(tokParenL.string()) {
		return nil, errors.New("Expecting '(' opener for variable definitions")
	} else {
		l.get()

		varDefs := &variableDefinitions{}

		for !l.tokenEquals(tokParenR.string()) {
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
func parseVariableDefinition(l *lexer) (*variableDefinition, error) {
	locStart := l.location().Start

	_var, err := parseVariable(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(tokColon.string()) {
		return nil, errors.New("Expecting a colon after variable name")
	}

	_typ, err := parseType(l)

	if err != nil {
		return nil, err
	}

	locEnd := _typ.Location().End

	defVal, _ := parseDefaultValue(l)

	directives, _ := parseDirectives(l)

	varDef := &variableDefinition{}

	varDef.Variable = *_var
	varDef.Type = _typ
	varDef.DefaultValue = defVal
	varDef.Directives = directives
	varDef.Loc = location{locStart, locEnd, l.source}

	return varDef, nil
}

// https://graphql.github.io/graphql-spec/draft/#Type
func parseType(l *lexer) (_type, error) {
	var _typ _type

	_typ, err := parseNamedType(l)

	if _typ != nil {
		return _typ, nil
	}

	_typ, err = parseListType(l)

	if _typ != nil {
		return _typ, nil
	}

	_typ, err = parseNonNullType(l)

	if err != nil {
		return nil, errors.Wrap(err, "Expecting a type")
	}

	return _typ, nil
}

// https://graphql.github.io/graphql-spec/draft/#ListType
func parseListType(l *lexer) (*listType, error) {
	locStart := l.location().Start

	if !l.tokenEquals(tokBracketL.string()) {
		return nil, errors.New("Expecting '[' for list type")
	}

	l.get()

	_typ, err := parseType(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(tokBracketR.string()) {
		return nil, errors.New("Expecting ']' for list type")
	}

	l.get()

	listTyp := &listType{}

	listTyp.OfType = _typ
	listTyp.Loc = location{locStart, l.prevLocation().End, l.source}

	return listTyp, nil
}

// https://graphql.github.io/graphql-spec/draft/#NonNullType
func parseNonNullType(l *lexer) (*nonNullType, error) {
	locStart := l.location().Start

	var _typ _type

	_typ, err := parseNamedType(l)

	if err != nil {
		_typ, err = parseListType(l)

		if err != nil {
			return nil, err
		}
	}

	if !l.tokenEquals(tokBang.string()) {
		return nil, errors.New("Expecting '!' at the end of a non null type")
	}

	l.get()

	nonNull := &nonNullType{}

	nonNull.OfType = _typ
	nonNull.Loc = location{locStart, l.prevLocation().End, l.source}

	return nonNull, nil
}

// https://graphql.github.io/graphql-spec/draft/#Directives
func parseDirectives(l *lexer) (*directives, error) {
	dirs := &directives{}

	for l.tokenEquals(tokAt.string()) {
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
func parseDirective(l *lexer) (*directive, error) {
	locStart := l.location().Start

	if !l.tokenEquals(tokAt.string()) {
		return nil, errors.New("Expecting '@' for directive")
	} else {
		l.get()

		nam, err := parseName(l)

		if err != nil {
			return nil, err
		}

		args, _ := parseArguments(l)

		dir := &directive{}

		dir.Name = *nam
		dir.Arguments = args
		dir.Loc = location{locStart, l.prevLocation().End, l.source}

		return dir, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#SelectionSet
func parseSelectionSet(l *lexer) (*selectionSet, error) {
	if !l.tokenEquals(tokBraceL.string()) {
		return nil, errors.New("Expecting '{' for selection set")
	} else {
		l.get()

		selSet := &selectionSet{}

		for !l.tokenEquals(tokBraceR.string()) {
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
func parseSelection(l *lexer) (selection, error) {
	var sel selection

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
func parseVariable(l *lexer) (*variable, error) {
	locStart := l.location().Start

	if !l.tokenEquals(tokDollar.string()) {
		return nil, errors.New("Expecting '$' for varible")
	} else {
		nam, err := parseName(l)

		if err != nil {
			return nil, err
		}

		_var := &variable{}

		_var.Name = *nam
		_var.Loc = location{locStart, nam.Location().End, l.source}

		return _var, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#DefaultValue
func parseDefaultValue(l *lexer) (*defaultValue, error) {
	locStart := l.location().Start

	if !l.tokenEquals(tokEquals.string()) {
		return nil, errors.New("Expecting '=' for default value")
	} else {
		val, err := parseValue(l)

		if err != nil {
			return nil, err
		}

		dVal := &defaultValue{}

		dVal.Value = val
		dVal.Loc = location{locStart, l.prevLocation().End, l.source}

		return dVal, nil
	}
}

// ! need to check variable type in order to parse its value
// https://graphql.github.io/graphql-spec/draft/#Value
func parseValue(l *lexer) (value, error) {
	// need to parse dynamic variables
	//_var, _ := parseVariable(l)

	var val value

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
func parseArguments(l *lexer) (*arguments, error) {
	if !l.tokenEquals(tokParenL.string()) {
		return nil, errors.New("Expecting '(' for arguments")
	} else {
		l.get()

		args := &arguments{}

		for !l.tokenEquals(tokParenR.string()) {
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
func parseArgument(l *lexer) (*argument, error) {
	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(tokColon.string()) {
		return nil, errors.New("Expecting colon after argument name")
	}

	l.get()

	val, err := parseValue(l)

	if err != nil {
		return nil, err
	}

	arg := &argument{}

	arg.Name = *nam
	arg.Value = val
	arg.Loc = location{nam.Location().Start, l.prevLocation().End, l.source}

	return arg, nil
}

// https://graphql.github.io/graphql-spec/draft/#Field
func parseField(l *lexer) (*field, error) {
	locStart := l.location().Start

	alia, err := parseName(l)

	if err != nil {
		return nil, err
	}

	nam := &name{}

	if l.tokenEquals(tokColon.string()) {
		l.get()

		nam, err = parseName(l)

		if err != nil {
			return nil, err
		}
	} else {
		*nam = *alia

		alia = nil
	}

	args, _ := parseArguments(l)

	dirs, _ := parseDirectives(l)

	selSet, _ := parseSelectionSet(l)

	fld := &field{}

	fld.Alias = (*alias)(alia)
	fld.Name = *nam
	fld.Arguments = args
	fld.Directives = dirs
	fld.SelectionSet = selSet
	fld.Loc = location{locStart, l.prevLocation().End, l.source}

	return fld, nil
}

// https://graphql.github.io/graphql-spec/draft/#FragmentSpread
func parseFragmentSpread(l *lexer) (*fragmentSpread, error) {
	locStart := l.location().Start

	if !l.tokenEquals(tokSpread.string()) {
		return nil, errors.New("Expecting '...' operator for a fragment spread")
	} else {
		l.get()

		fname, err := parseFragmentName(l)

		if err != nil {
			return nil, err
		}

		directives, _ := parseDirectives(l)

		spread := &fragmentSpread{}

		spread.FragmentName = *fname
		spread.Directives = directives
		spread.Loc = location{locStart, l.prevLocation().End, l.source}

		return spread, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#InlineFragment
func parseInlineFragment(l *lexer) (*inlineFragment, error) {
	locStart := l.location().Start

	if !l.tokenEquals(tokSpread.string()) {
		return nil, errors.New("Expecting '...' for an inline fragment")
	} else {
		l.get()

		typeCon, _ := parseTypeCondition(l)

		directives, _ := parseDirectives(l)

		selSet, err := parseSelectionSet(l)

		if err != nil {
			return nil, err
		}

		inlineFrag := &inlineFragment{}

		inlineFrag.TypeCondition = typeCon
		inlineFrag.Directives = directives
		inlineFrag.SelectionSet = *selSet
		inlineFrag.Loc = location{locStart, l.prevLocation().End, l.source}

		return inlineFrag, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentName
func parseFragmentName(l *lexer) (*fragmentName, error) {
	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if nam.Value == kwOn {
		return nil, errors.New("Fragment name cannot be 'on'")
	}

	var fragName *fragmentName

	*fragName = fragmentName(*nam)
	fragName.Loc = *nam.Location()

	return fragName, nil
}

// https://graphql.github.io/graphql-spec/draft/#TypeCondition
func parseTypeCondition(l *lexer) (*typeCondition, error) {
	locStart := l.location().Start

	if !l.tokenEquals(kwOn) {
		return nil, errors.New("Expecting 'on' keyword for a type condition")
	} else {
		namedTyp, err := parseNamedType(l)

		if err != nil {
			return nil, err
		}

		typeCond := &typeCondition{}

		typeCond.NamedType = *namedTyp
		typeCond.Loc = location{locStart, namedTyp.Location().End, l.source}

		return typeCond, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#NamedType
func parseNamedType(l *lexer) (*namedType, error) {
	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	var namedTyp *namedType

	*namedTyp = namedType(*nam)

	return namedTyp, nil
}

// ! Check numeric values
// https://graphql.github.io/graphql-spec/draft/#IntValue
func parseIntValue(l *lexer) (*intValue, error) {
	tok := l.current()

	intVal, err := strconv.ParseInt(tok.value, 10, 64)

	if err != nil {
		return nil, err
	}

	l.get()

	intValP := &intValue{}

	intValP.Value = intVal
	intValP.Loc = location{tok.start, tok.end, l.source}

	return intValP, nil
}

// https://graphql.github.io/graphql-spec/draft/#FloatValue
func parseFloatValue(l *lexer) (*floatValue, error) {
	tok := l.current()

	floatVal, err := strconv.ParseFloat(tok.value, 64)

	if err != nil {
		return nil, err
	}

	l.get()

	floatValP := &floatValue{}

	floatValP.Value = floatVal
	floatValP.Loc = location{tok.start, tok.end, l.source}

	return floatValP, nil
}

// ! Have a discussion about this function
// https://graphql.github.io/graphql-spec/draft/#StringValue
func parseStringValue(l *lexer) (*stringValue, error) {
	tok := l.current()

	sv := &stringValue{}

	sv.Value = tok.value
	sv.Loc = location{tok.start, tok.end, l.source}

	return sv, nil
}

// https://graphql.github.io/graphql-spec/draft/#BooleanValue
func parseBooleanValue(l *lexer) (*booleanValue, error) {
	tok := l.current()

	boolVal, err := strconv.ParseBool(tok.value)

	if err != nil {
		return nil, err
	}

	l.get()

	boolValP := &booleanValue{}

	boolValP.Value = boolVal
	boolValP.Loc = location{tok.start, tok.end, l.source}

	return boolValP, nil
}

// https://graphql.github.io/graphql-spec/draft/#NullValue
func parseNullValue(l *lexer) (*nullValue, error) {
	tok := l.current()

	if tok.value != kwNull {
		return nil, errors.New("Expecting 'null' keyword")
	} else {
		l.get()

		null := &nullValue{}
		null.Loc = location{tok.start, tok.end, l.source}

		return null, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#EnumValue
func parseEnumValue(l *lexer) (*enumValue, error) {
	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	switch nam.Value {
	case kwTrue, kwFalse, kwNull:
		return nil, errors.New("Enum value cannot be 'true', 'false' or 'null'")
	default:
		enumVal := &enumValue{}

		enumVal.Name = *nam
		enumVal.Loc = location{nam.Location().Start, nam.Location().End, l.source}

		return enumVal, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#ListValue
func parseListValue(l *lexer) (*listValue, error) {
	locStart := l.location().Start

	if !l.tokenEquals(tokBracketL.string()) {
		return nil, errors.New("Expecting '[' for a list value")
	} else {
		l.get()

		lstVal := &listValue{}

		for !l.tokenEquals(tokBracketR.string()) {
			val, err := parseValue(l)

			if err != nil {
				return nil, err
			}

			lstVal.Values = append(lstVal.Values, val)
		}

		lstVal.Loc = location{locStart, l.prevLocation().End, l.source}

		return lstVal, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectValue
func parseObjectValue(l *lexer) (*objectValue, error) {
	locStart := l.location().Start

	if !l.tokenEquals(tokBraceL.string()) {
		return nil, errors.New("Expecting '{' for an object value")
	} else {
		l.get()

		objVal := &objectValue{}

		for !l.tokenEquals(tokBraceR.string()) {
			objField, err := parseObjectField(l)

			if err != nil {
				return nil, err
			}

			objVal.Values = append(objVal.Values, *objField)
		}

		objVal.Loc = location{locStart, l.prevLocation().End, l.source}

		return objVal, nil
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectField
func parseObjectField(l *lexer) (*objectField, error) {
	nam, err := parseName(l)

	if err != nil {
		return nil, err
	}

	if !l.tokenEquals(tokColon.string()) {
		return nil, errors.New("Expecting color after object field name")
	}

	l.get()

	val, err := parseValue(l)

	if err != nil {
		return nil, err
	}

	objField := &objectField{}

	objField.Name = *nam
	objField.Value = val
	objField.Loc = location{nam.Location().Start, l.prevLocation().End, l.source}

	return objField, nil
}
