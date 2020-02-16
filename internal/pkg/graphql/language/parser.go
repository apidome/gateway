package language

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func isNilInterface(i interface{}) bool {
	return reflect.ValueOf(i).IsNil()
}

func Parse(schema *document, doc string) (ret *document, err error) {
	l, err := newlexer(doc)

	//recover syntax errors
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	ret = parseDocument(l)

	validateDocument(schema, ret)

	return
}

func ParseSchema(doc string) (ret *document, err error) {
	l, err := newlexer(doc)

	//recover syntax errors
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	ret = parseDocument(l)

	return
}

// https://graphql.github.io/graphql-spec/draft/#Document
func parseDocument(l *lexer) *document {
	doc := &document{}

	def := parseDefinitions(l)

	doc.definitions = *def

	return doc
}

// https://graphql.github.io/graphql-spec/draft/#Definition
func parseDefinitions(l *lexer) *definitions {
	defs := &definitions{}

	for !l.tokenEquals(tokEOF.string()) {
		def := parseDefinition(l)

		*defs = append(*defs, def)
	}

	if len(*defs) == 0 {
		panic(errors.New("No definitions found in document"))
	}

	return defs
}

// https://graphql.github.io/graphql-spec/draft/#Definition
func parseDefinition(l *lexer) definition {
	if executableDefinitionExists(l) {
		return parseExecutableDefinition(l)
	}

	if typeSystemDefinitionExists(l) {
		return parseTypeSystemDefinition(l)
	}

	if typeSystemExtensionExists(l) {
		return parseTypeSystemExtension(l)
	}

	panic(errors.New("Expecting one of 'executable definition', 'type system definition', 'type system extension'"))
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDefinition
func parseExecutableDefinition(l *lexer) executableDefinition {
	if operationDefinitionExists(l) {
		return parseOperationDefinition(l)
	}

	if fragmentDefinitionExists(l) {
		return parseFragmentDefinition(l)
	}

	panic(errors.New("Expecting one of 'operation definition', 'fragment definition'"))
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDefinition
func parseTypeSystemDefinition(l *lexer) typeSystemDefinition {
	if schemaDefinitionExists(l) {
		return parseSchemaDefinition(l)
	}

	if typeDefinitionExists(l) {
		return parseTypeDefinition(l)
	}

	if directiveDefinitionExists(l) {
		return parseDirectiveDefinition(l)
	}

	panic(errors.New("Expecting one of 'schema definition', 'type definition', 'directive definition'"))
}

// https://graphql.github.io/graphql-spec/draft/#SchemaDefinition
func parseSchemaDefinition(l *lexer) *schemaDefinition {
	schDef := &schemaDefinition{}

	locStart := l.current().start

	if !l.tokenEquals(kwSchema) {
		panic(errors.New("Missing 'schema' keyword for a schema definition"))
	}

	l.get()

	if directivesExist(l) {
		schDef.directives = parseDirectives(l)
	}

	if !l.tokenEquals(tokBraceL.string()) {
		panic(errors.New("Missing '{' for a schema definition"))
	}

	l.get()

	schDef.rootOperationTypeDefinitions = *parseRootOperationTypeDefinitions(l)

	if !l.tokenEquals(tokBraceR.string()) {
		panic(errors.New("Missing '}' for schema definition"))
	}

	locEnd := l.get().end

	schDef.loc = location{locStart, locEnd, l.source}

	return schDef
}

// https://graphql.github.io/graphql-spec/draft/#RootOperationTypeDefinition
func parseRootOperationTypeDefinitions(l *lexer) *rootOperationTypeDefinitions {
	rotds := &rootOperationTypeDefinitions{}

	for !l.tokenEquals(tokBraceR.string()) {
		rotd := parseRootOperationTypeDefinition(l)

		*rotds = append(*rotds, rotd)
	}

	l.get()

	if len(*rotds) == 0 {
		panic(errors.New("Expecting at least one root operation type definition"))
	}

	return rotds
}

// https://graphql.github.io/graphql-spec/draft/#RootOperationTypeDefinition
func parseRootOperationTypeDefinition(l *lexer) *rootOperationTypeDefinition {
	rotd := &rootOperationTypeDefinition{}

	locStart := l.current().start

	rotd.operationType = *parseOperationType(l)

	if !l.tokenEquals(tokColon.string()) {
		panic(errors.New("Expecting ':' after operation type"))
	}

	l.get()

	rotd.namedType = *parseNamedType(l)

	rotd.loc = location{locStart, rotd.namedType.Location().End(), l.source}

	return rotd
}

// https://graphql.github.io/graphql-spec/draft/#OperationType
func parseOperationType(l *lexer) *operationType {
	opType := new(operationType)

	tok := l.current()

	if tok.value != string(operationMutation) &&
		tok.value != string(operationQuery) &&
		tok.value != string(operationSubscription) {
		panic(errors.New("Expecting 'query', 'mutation' or 'subscription' as operation type"))
	}

	*opType = (operationType)(tok.value)

	return opType
}

// https://graphql.github.io/graphql-spec/draft/#TypeDefinition
func parseTypeDefinition(l *lexer) typeDefinition {
	if scalarTypeDefinitionExists(l) {
		return parseScalarTypeDefinition(l)
	}

	if objectTypeDefinitionExists(l) {
		return parseObjectTypeDefinition(l)
	}

	if interfaceTypeDefinitionExists(l) {
		return parseInterfaceTypeDefinition(l)
	}

	if unionTypeDefinitionExists(l) {
		return parseUnionTypeDefinition(l)
	}

	if enumTypeDefinitionExists(l) {
		return parseEnumTypeDefinition(l)
	}

	if inputObjectTypeDefinitionExists(l) {
		return parseInputObjectTypeDefinition(l)
	}

	panic(errors.New("Expecting a type definition"))
}

// https://graphql.github.io/graphql-spec/draft/#ScalarTypeDefinition
func parseScalarTypeDefinition(l *lexer) *scalarTypeDefinition {
	stp := &scalarTypeDefinition{}

	if descriptionExists(l) {
		stp.description = parseDescription(l)
	}

	if !l.tokenEquals(tsdlScalar) {
		panic(errors.New("Missing 'scalar' keyword for scalar type definition"))
	}

	tok := l.get()

	stp.name = *parseName(l)

	if directivesExist(l) {
		stp.directives = parseDirectives(l)
	}

	stp.loc = location{tok.start, l.prevLocation().End(), l.source}

	return stp
}

// https://graphql.github.io/graphql-spec/draft/#Description
func parseDescription(l *lexer) *description {
	strVal := parseStringValue(l)

	return (*description)(strVal)
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeDefinition
func parseObjectTypeDefinition(l *lexer) *objectTypeDefinition {
	otd := &objectTypeDefinition{}

	locStart := l.location().Start()

	if descriptionExists(l) {
		otd.description = parseDescription(l)
	}

	if !l.tokenEquals(kwType) {
		panic(errors.New("Expecting 'type' keyword for object type definition"))
	}

	l.get()

	otd.name = *parseName(l)

	if implementsInterfacesExists(l) {
		otd.implementsInterfaces = parseImplementsInterfaces(l)
	}

	if directivesExist(l) {
		otd.directives = parseDirectives(l)
	}

	if fieldsDefinitionExists(l) {
		otd.fieldsDefinition = parseFieldsDefinition(l)
	}

	otd.loc = location{locStart, l.prevLocation().End(), l.source}

	return otd
}

// https://graphql.github.io/graphql-spec/draft/#ImplementsInterfaces
func parseImplementsInterfaces(l *lexer) *implementsInterfaces {
	ii := &implementsInterfaces{}

	if !l.tokenEquals(kwImplements) {
		panic(errors.New("Expecting 'implements' keyword"))
	}

	l.get()

	if l.tokenEquals(tokAmp.string()) {
		l.get()
	}

	*ii = append(*ii, parseNamedType(l))

	for l.tokenEquals(tokAmp.string()) {
		l.get()

		(*ii) = append(*ii, parseNamedType(l))
	}

	return ii
}

// https://graphql.github.io/graphql-spec/draft/#FieldsDefinition
func parseFieldsDefinition(l *lexer) *fieldsDefinition {
	fds := &fieldsDefinition{}

	if !l.tokenEquals(tokBraceL.string()) {
		panic(errors.New("Expecting '{' for fields definition"))
	}

	l.get()

	for !l.tokenEquals(tokBraceR.string()) {
		(*fds) = append(*fds, parseFieldDefinition(l))
	}

	l.get()

	if len(*fds) == 0 {
		panic(errors.New("Expecting at lease one field definition"))
	}

	return fds
}

// https://graphql.github.io/graphql-spec/draft/#FieldsDefinition
func parseFieldDefinition(l *lexer) *fieldDefinition {
	fd := &fieldDefinition{}

	locStart := l.location().Start()

	if descriptionExists(l) {
		fd.description = parseDescription(l)
	}

	fd.name = *parseName(l)

	if argumentsDefinitionExist(l) {
		fd.argumentsDefinition = parseArgumentsDefinition(l)
	}

	if !l.tokenEquals(tokColon.string()) {
		panic(errors.New("Expecting ':' for field definition"))
	}

	l.get()

	fd._type = parseType(l)

	if directivesExist(l) {
		fd.directives = parseDirectives(l)
	}

	fd.loc = location{locStart, l.prevLocation().End(), l.source}

	return fd
}

// https://graphql.github.io/graphql-spec/draft/#ArgumentsDefinition
func parseArgumentsDefinition(l *lexer) *argumentsDefinition {
	argsDef := &argumentsDefinition{}

	if !l.tokenEquals(tokParenL.string()) {
		panic(errors.New("Expecting '(' for arguments definition"))
	}

	l.get()

	for !l.tokenEquals(tokParenR.string()) {
		*argsDef = append(*argsDef, parseInputValueDefinition(l))
	}

	l.get()

	if len(*argsDef) == 0 {
		panic(errors.New("Expecting at least one input value definitions"))
	}

	return argsDef
}

// https://graphql.github.io/graphql-spec/draft/#InputValueDefinition
func parseInputValueDefinition(l *lexer) *inputValueDefinition {
	ivd := &inputValueDefinition{}

	locStart := l.location().Start()

	if descriptionExists(l) {
		ivd.description = parseDescription(l)
	}

	ivd.name = *parseName(l)

	if !l.tokenEquals(tokColon.string()) {
		panic(errors.New("Expecting ':' for input value definition"))
	}

	l.get()

	ivd._type = parseType(l)

	if defaultValueExists(l) {
		ivd.defaultValue = parseDefaultValue(l)
	}

	if directivesExist(l) {
		ivd.directives = parseDirectives(l)
	}

	ivd.loc = location{locStart, l.prevLocation().End(), l.source}

	return ivd
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeDefinition
func parseInterfaceTypeDefinition(l *lexer) *interfaceTypeDefinition {
	itd := &interfaceTypeDefinition{}

	locStart := l.location().Start()

	if descriptionExists(l) {
		itd.description = parseDescription(l)
	}

	if !l.tokenEquals(kwInterface) {
		panic(errors.New("Expecting 'interface' keyword for interface type definition"))
	}

	l.get()

	itd.name = *parseName(l)

	if directivesExist(l) {
		itd.directives = parseDirectives(l)
	}

	if fieldsDefinitionExists(l) {
		itd.fieldsDefinition = parseFieldsDefinition(l)
	}

	itd.loc = location{locStart, l.prevLocation().End(), l.source}

	return itd
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeDefinition
func parseUnionTypeDefinition(l *lexer) *unionTypeDefinition {
	utd := &unionTypeDefinition{}

	locStart := l.location().Start()

	if descriptionExists(l) {
		utd.description = parseDescription(l)
	}

	if !l.tokenEquals(kwUnion) {
		panic(errors.New("Expecting 'union' keyowrd for union type definition"))
	}

	l.get()

	utd.name = *parseName(l)

	if directivesExist(l) {
		utd.directives = parseDirectives(l)
	}

	if unionMemberTypesExist(l) {
		utd.unionMemberTypes = parseUnionMemberTypes(l)
	}

	utd.loc = location{locStart, l.prevLocation().End(), l.source}

	return utd
}

// https://graphql.github.io/graphql-spec/draft/#UnionMemberTypes
func parseUnionMemberTypes(l *lexer) *unionMemberTypes {
	umt := &unionMemberTypes{}

	if !l.tokenEquals(tokEquals.string()) {
		panic(errors.New("Expecting '=' for union member types"))
	}

	l.get()

	if l.tokenEquals(tokPipe.string()) {
		l.get()
	}

	*umt = append(*umt, parseNamedType(l))

	for l.tokenEquals(tokPipe.string()) {
		l.get()

		*umt = append(*umt, parseNamedType(l))
	}

	if len(*umt) == 0 {
		panic(errors.New("Expecting at least one union member type"))
	}

	return umt
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeDefinition
func parseEnumTypeDefinition(l *lexer) *enumTypeDefinition {
	etd := &enumTypeDefinition{}

	locStart := l.location().Start()

	if descriptionExists(l) {
		etd.description = parseDescription(l)
	}

	if !l.tokenEquals(kwEnum) {
		panic(errors.New("Expecting 'enum' keyword for enum type definition"))
	}

	l.get()

	etd.name = *parseName(l)

	if directivesExist(l) {
		etd.directives = parseDirectives(l)
	}

	if enumValuesDefinitionExist(l) {
		etd.enumValuesDefinition = parseEnumValuesDefinition(l)
	}

	etd.loc = location{locStart, l.prevLocation().End(), l.source}

	return etd
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinition
func parseEnumValuesDefinition(l *lexer) *enumValuesDefinition {
	evd := &enumValuesDefinition{}

	if !l.tokenEquals(tokBraceL.string()) {
		panic(errors.New("Expecting '{' for enum values definition"))
	}

	l.get()

	for !l.tokenEquals(tokBraceR.string()) {
		*evd = append(*evd, parseEnumValueDefinition(l))
	}

	l.get()

	if len(*evd) == 0 {
		panic(errors.New("Expecting at least one enum value definition"))
	}

	return evd
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinitionvd
func parseEnumValueDefinition(l *lexer) *enumValueDefinition {
	evd := &enumValueDefinition{}

	locStart := l.location().Start()

	if descriptionExists(l) {
		evd.description = parseDescription(l)
	}

	evd.enumValue = *parseEnumValue(l)

	if directivesExist(l) {
		evd.directives = parseDirectives(l)
	}

	evd.loc = location{locStart, l.prevLocation().End(), l.source}

	return evd
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeDefinition
func parseInputObjectTypeDefinition(l *lexer) *inputObjectTypeDefinition {
	iotd := &inputObjectTypeDefinition{}

	locStart := l.location().Start()

	if descriptionExists(l) {
		iotd.description = parseDescription(l)
	}

	if !l.tokenEquals(kwInput) {
		panic(errors.New("Expecting 'input' keyword for input object type definition"))
	}

	l.get()

	iotd.name = *parseName(l)

	if directivesExist(l) {
		iotd.directives = parseDirectives(l)
	}

	if inputFieldsDefinitionExists(l) {
		iotd.inputFieldsDefinition = parseInputFieldsDefinition(l)
	}

	iotd.loc = location{locStart, l.prevLocation().End(), l.source}

	return iotd
}

// https://graphql.github.io/graphql-spec/draft/#InputFieldsDefinition
func parseInputFieldsDefinition(l *lexer) *inputFieldsDefinition {
	ifd := &inputFieldsDefinition{}

	if !l.tokenEquals(tokBraceL.string()) {
		panic(errors.New("Expecting '{' for input fields definition"))
	}

	l.get()

	for !l.tokenEquals(tokBraceR.string()) {
		*ifd = append(*ifd, parseInputValueDefinition(l))
	}

	l.get()

	if len(*ifd) == 0 {
		panic(errors.New("Expecting at least one input field definition"))
	}

	return ifd
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveDefinition
func parseDirectiveDefinition(l *lexer) *directiveDefinition {
	dd := &directiveDefinition{}

	locStart := l.location().Start()

	if descriptionExists(l) {
		dd.description = parseDescription(l)
	}

	if !l.tokenEquals(kwDirective) {
		panic(errors.New("Expecting 'directive' keyword for directive definition"))
	}

	l.get()

	if !l.tokenEquals(tokAt.string()) {
		panic(errors.New("Expecting '@' for directive definition"))
	}

	l.get()

	dd.name = *parseName(l)

	if argumentsDefinitionExist(l) {
		dd.argumentsDefinition = parseArgumentsDefinition(l)
	}

	if l.tokenEquals(kwRepeatable) {
		l.get()
	}

	if !l.tokenEquals(kwOn) {
		panic(errors.New("Expecting 'on' keyworkd for directive definition"))
	}

	l.get()

	dd.directiveLocations = *parseDirectiveLocations(l)

	dd.loc = location{locStart, l.prevLocation().End(), l.source}

	return dd
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocations
func parseDirectiveLocations(l *lexer) *directiveLocations {
	dls := &directiveLocations{}

	if l.tokenEquals(tokPipe.string()) {
		l.get()
	}

	*dls = append(*dls, *parseDirectiveLocation(l))

	for l.tokenEquals(tokPipe.string()) {
		l.get()

		*dls = append(*dls, *parseDirectiveLocation(l))
	}

	return dls
}

// https://graphql.github.io/graphql-spec/draft/#TypeExtension
func parseTypeExtension(l *lexer) typeExtension {
	if !l.tokenEquals(kwExtend) {
		panic(errors.New("Expecting 'extend' keyword"))
	}

	l.get()

	if scalarTypeDefinitionExists(l) {
		return parseScalarTypeExtension(l)
	}

	if objectTypeDefinitionExists(l) {
		return parseObjectTypeExtension(l)
	}

	if interfaceTypeExtensionExists(l) {
		return parseInterfaceTypeExtension(l)
	}

	if unionTypeExtensionExists(l) {
		return parseUnionTypeExtension(l)
	}

	if enumTypeExtensionExists(l) {
		return parseEnumTypeExtension(l)
	}

	if inputObjectTypeExtensionExists(l) {
		return parseInputObjectTypeExtension(l)
	}

	panic(errors.New("Expecting type extension"))
}

// https://graphql.github.io/graphql-spec/draft/#ScalarTypeExtension
func parseScalarTypeExtension(l *lexer) *scalarTypeExtension {
	ste := &scalarTypeExtension{}

	locStart := l.location().Start()

	if !l.tokenEquals(kwScalar) {
		panic(errors.New("Expecting 'extend scalar' keywords for scalar type extension"))
	}

	l.get()

	ste.name = *parseName(l)
	ste.directives = *parseDirectives(l)
	ste.loc = location{locStart, l.prevLocation().End(), l.source}

	return ste
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeExtension
func parseObjectTypeExtension(l *lexer) *objectTypeExtension {
	ote := &objectTypeExtension{}

	locStart := l.location().Start()

	if !l.tokenEquals(kwType) {
		panic(errors.New("Expecting 'extend type' keywords for object type extension"))
	}

	l.get()

	ote.name = *parseName(l)

	if implementsInterfacesExists(l) {
		ote.implementsInterfaces = parseImplementsInterfaces(l)
	}

	if directivesExist(l) {
		ote.directives = parseDirectives(l)
	}

	if fieldsDefinitionExists(l) {
		ote.fieldsDefinition = parseFieldsDefinition(l)
	}

	if ote.ImplementsInterfaces == nil &&
		ote.Directives == nil &&
		ote.FieldsDefinition == nil {
		panic(errors.New("Expecting at least one of 'implements interface', 'directives', 'fields definition' for object type extension"))
	}

	ote.loc = location{locStart, l.prevLocation().End(), l.source}

	return ote
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeExtension
func parseInterfaceTypeExtension(l *lexer) *interfaceTypeExtension {
	ite := &interfaceTypeExtension{}

	locStart := l.location().Start()

	if !l.tokenEquals(kwInterface) {
		panic(errors.New("Expecting 'extend interface' keywords for interface type extension"))
	}

	l.get()

	ite.name = *parseName(l)

	if directivesExist(l) {
		ite.directives = parseDirectives(l)
	}

	if fieldsDefinitionExists(l) {
		ite.fieldsDefinition = parseFieldsDefinition(l)
	}

	if ite.Directives == nil && ite.FieldsDefinition == nil {
		panic(errors.New("Expecting at least one of 'directives', 'fields definition' for interface type extension"))
	}

	ite.loc = location{locStart, l.prevLocation().End(), l.source}

	return ite
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeExtension
func parseUnionTypeExtension(l *lexer) *unionTypeExtension {
	ute := &unionTypeExtension{}

	locStart := l.location().Start()

	if !l.tokenEquals(tsdlUnion) {
		panic(errors.New("Expecting 'extend union' keywords for union type extension"))
	}

	l.get()

	ute.name = *parseName(l)

	if directivesExist(l) {
		ute.directives = parseDirectives(l)
	}

	if unionMemberTypesExist(l) {
		ute.unionMemberTypes = parseUnionMemberTypes(l)
	}

	if ute.Directives == nil && ute.UnionMemberTypes == nil {
		panic(errors.New("Expecting at  least one of 'directives', 'union member types' for union type extension"))
	}

	ute.loc = location{locStart, l.prevLocation().End(), l.source}

	return ute
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeExtension
func parseEnumTypeExtension(l *lexer) *enumTypeExtension {
	ete := &enumTypeExtension{}

	locStart := l.location().Start()

	if !l.tokenEquals(tsdlEnum) {
		panic(errors.New("Expecting 'extend enum' keywords for enum type extension"))
	}

	l.get()

	ete.name = *parseName(l)

	if directivesExist(l) {
		ete.directives = parseDirectives(l)
	}

	if enumValuesDefinitionExist(l) {
		ete.enumValuesDefinition = parseEnumValuesDefinition(l)
	}

	if ete.Directives == nil && ete.EnumValuesDefinition == nil {
		panic(errors.New("Expecting at least one of 'directives', 'enum values definition' for enum type extension"))
	}

	ete.loc = location{locStart, l.prevLocation().End(), l.source}

	return ete
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocation
func parseDirectiveLocation(l *lexer) *directiveLocation {
	if executableDirectiveLocationExists(l) {
		return (*directiveLocation)(parseExecutableDirectiveLocation(l))
	}

	if typeSystemDirectiveLocationExists(l) {
		return (*directiveLocation)(parseTypeSystemDirectiveLocation(l))
	}

	panic(errors.New("Expecting a directive location"))
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDirectiveLocation
func parseExecutableDirectiveLocation(l *lexer) *executableDirectiveLocation {
	tok := l.current()

	for i := range executableDirectiveLocations {
		if string(executableDirectiveLocations[i]) == tok.value {
			l.get()

			return &executableDirectiveLocations[i]
		}
	}

	panic(errors.New("Expecting executable directive location"))
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDirectiveLocation
func parseTypeSystemDirectiveLocation(l *lexer) *typeSystemDirectiveLocation {
	tok := l.current()

	for i := range typeSystemDirectiveLocations {
		if string(typeSystemDirectiveLocations[i]) == tok.value {
			l.get()

			return &typeSystemDirectiveLocations[i]
		}
	}

	panic(errors.New("Expecting type systen directive location"))
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemExtension
func parseTypeSystemExtension(l *lexer) typeSystemExtension {
	if schemaExtensionExists(l) {
		return parseSchemaExtension(l)
	}

	if typeExtensionExists(l) {
		return parseTypeExtension(l)
	}

	panic(errors.New("Expecting type system extension"))
}

// https://graphql.github.io/graphql-spec/draft/#SchemaExtension
func parseSchemaExtension(l *lexer) *schemaExtension {
	se := &schemaExtension{}

	locStart := l.location().Start()

	if !l.tokenEquals(kwExtend, kwSchema) {
		panic(errors.New("Expecting 'extend schema' keywords for schema extension"))
	}

	l.get()
	l.get()

	if directivesExist(l) {
		se.directives = parseDirectives(l)
	}

	if !l.tokenEquals(tokBraceL.string()) {
		panic(errors.New("Expecting '{' for schema extension"))
	}

	l.get()

	if rootOperationTypeDefinitionsExist(l) {
		se.rootOperationTypeDefinitions = parseRootOperationTypeDefinitions(l)
	}

	l.get()

	if se.Directives == nil && se.RootOperationTypeDefinitions == nil {
		panic(errors.New("Expecting directives or root operation type definitions for schema extension"))
	}

	se.loc = location{locStart, l.prevLocation().End(), l.source}

	return se
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeExtension
func parseInputObjectTypeExtension(l *lexer) *inputObjectTypeExtension {
	iote := &inputObjectTypeExtension{}

	locStart := l.location().Start()

	if !l.tokenEquals(kwInput) {
		panic(errors.New("Expecting 'extend' keyword for input object type extension"))
	}

	l.get()

	iote.name = *parseName(l)

	if directivesExist(l) {
		iote.directives = parseDirectives(l)
	}

	if inputFieldsDefinitionExists(l) {
		iote.inputFieldsDefinition = parseInputFieldsDefinition(l)
	}

	if iote.Directives == nil && iote.InputFieldsDefinition == nil {
		panic(errors.New("Expecting at lease one of 'directives', 'input fields definition' fo input object type extension"))
	}

	iote.loc = location{locStart, l.prevLocation().End(), l.source}

	return iote
}

// https://graphql.github.io/graphql-spec/draft/#OperationDefinition
func parseOperationDefinition(l *lexer) *operationDefinition {

	locStart := l.location().Start()

	// Shorthand query
	// https://graphql.github.io/graphql-spec/draft/#sec-Language.Operations.Query-shorthand
	if l.tokenEquals(tokBraceL.string()) {
		od := &operationDefinition{}
		od.operationType = kwQuery
		od.selectionSet = *parseSelectionSet(l)
		od.loc = location{locStart, l.prevLocation().End(), l.source}

		return od
	} else if !l.tokenEquals(kwQuery) &&
		!l.tokenEquals(kwMutation) &&
		!l.tokenEquals(kwSubscription) {
		panic(errors.New("Expecting one of 'query', 'mutation', 'subscription' for operation definition"))
	} else {
		od := &operationDefinition{}

		tok := l.get()

		opType := tok.value

		if nameExists(l) {
			od.name = parseName(l)
		}

		if variableDefinitionsExist(l) {
			od.variableDefinitions = parseVariableDefinitions(l)
		}

		if directivesExist(l) {
			od.directives = parseDirectives(l)
		}

		od.selectionSet = *parseSelectionSet(l)

		od.operationType = operationType(opType)
		od.loc = location{locStart, tok.end, l.source}

		return od
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentDefinition
func parseFragmentDefinition(l *lexer) *fragmentDefinition {

	locStart := l.location().Start()

	if !l.tokenEquals(kwFragment) {
		panic(errors.New("Expecting fragment keyword"))
	} else {
		fd := &fragmentDefinition{}

		l.get()

		fd.fragmentName = *parseFragmentName(l)

		if fd.fragmentName.Value() == kwOn {
			panic(errors.New("Fragment name cannot be 'on'"))
		}

		fd.typeCondition = *parseTypeCondition(l)

		if directivesExist(l) {
			fd.directives = parseDirectives(l)
		}

		fd.selectionSet = *parseSelectionSet(l)
		fd.loc = location{locStart, l.prevLocation().End(), l.source}

		return fd
	}
}

// https://graphql.github.io/graphql-spec/draft/#Name
func parseName(l *lexer) *name {
	tok := l.current()

	pattern := "^[_A-Za-z][_0-9A-Za-z]*$"

	// If the current token is not a Name, return nil
	if tok.kind != tokName {
		panic(errors.New("Not a name"))
	}

	// Check if the given name matches the regex provided by graphql spec at
	// https://graphql.github.io/graphql-spec/draft/#Name
	match, err := regexp.MatchString(pattern, tok.value)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse name: "))
	}

	// If the name does not match the requirements, return an error.
	if !match {
		panic(errors.New("invalid name - " + tok.value))
	}

	l.get()

	nm := &name{}

	// Populate the Name struct.
	nm.value = tok.value
	nm.loc.start = tok.start
	nm.loc.end = tok.end
	nm.loc.source = l.source

	// Return the AST Name object.
	return nm
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinition
func parseVariableDefinitions(l *lexer) *variableDefinitions {
	if !l.tokenEquals(tokParenL.string()) {
		panic(errors.New("Expecting '(' opener for variable definitions"))
	} else {
		l.get()

		vd := &variableDefinitions{}

		for !l.tokenEquals(tokParenR.string()) {
			*vd = append(*vd, parseVariableDefinition(l))
		}

		l.get()

		if len(*vd) == 0 {
			panic(errors.New("Expecting at least one variable definition"))
		}

		return vd
	}
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinition
func parseVariableDefinition(l *lexer) *variableDefinition {
	vd := &variableDefinition{}

	locStart := l.location().Start()

	vd.variable = *parseVariable(l)

	if !l.tokenEquals(tokColon.string()) {
		panic(errors.New("Expecting a colon after variable name"))
	}

	l.get()

	vd._type = parseType(l)

	if defaultValueExists(l) {
		vd.defaultValue = parseDefaultValue(l)
	}

	if directivesExist(l) {
		vd.directives = parseDirectives(l)
	}

	vd.loc = location{locStart, l.prevLocation().End(), l.source}

	return vd
}

// https://graphql.github.io/graphql-spec/draft/#Type
func parseType(l *lexer) _type {
	var retType _type

	if listTypeExists(l) {
		retType = parseListType(l)
	}

	if namedTypeExists(l) {
		retType = parseNamedType(l)
	}

	if nonNullTypeExists(l) {
		retType = parseNonNullType(l, retType)
	}

	return retType
}

// https://graphql.github.io/graphql-spec/draft/#ListType
func parseListType(l *lexer) *listType {
	lt := &listType{}

	locStart := l.location().Start()

	if !l.tokenEquals(tokBracketL.string()) {
		panic(errors.New("Expecting '[' for list type"))
	}

	l.get()

	lt.OfType = parseType(l)

	if !l.tokenEquals(tokBracketR.string()) {
		panic(errors.New("Expecting ']' for list type"))
	}

	l.get()

	lt.loc = location{locStart, l.prevLocation().End(), l.source}

	return lt
}

// https://graphql.github.io/graphql-spec/draft/#NonNullType
func parseNonNullType(l *lexer, t _type) *nonNullType {
	nnt := &nonNullType{}

	nnt.ofType = t
	nnt.loc = location{t.Location().Start(), l.location().End(), l.source}

	l.get()

	return nnt
}

// https://graphql.github.io/graphql-spec/draft/#Directives
func parseDirectives(l *lexer) *directives {

	dirs := &directives{}

	for l.tokenEquals(tokAt.string()) {
		*dirs = append(*dirs, parseDirective(l))
	}

	if len(*dirs) == 0 {
		panic(errors.New("Expecting at least one directive"))
	}

	return dirs
}

// https://graphql.github.io/graphql-spec/draft/#Directive
func parseDirective(l *lexer) *directive {
	dir := &directive{}

	locStart := l.location().Start()

	if !l.tokenEquals(tokAt.string()) {
		panic(errors.New("Expecting '@' for directive"))
	} else {
		l.get()

		dir.name = *parseName(l)

		if argumentsExist(l) {
			dir.arguments = parseArguments(l)
		}

		dir.loc = location{locStart, l.prevLocation().End(), l.source}

		return dir
	}
}

// https://graphql.github.io/graphql-spec/draft/#SelectionSet
func parseSelectionSet(l *lexer) *selectionSet {
	ss := &selectionSet{}

	if !l.tokenEquals(tokBraceL.string()) {
		panic(errors.New("Expecting '{' for selection set"))
	} else {
		l.get()

		for !l.tokenEquals(tokBraceR.string()) {
			*ss = append(*ss, parseSelection(l))
		}

		l.get()

		if len(*ss) == 0 {
			panic(errors.New("Expecting at least one selection"))
		}

		return ss
	}
}

// https://graphql.github.io/graphql-spec/draft/#Selection
func parseSelection(l *lexer) selection {
	if fieldExists(l) {
		return parseField(l)
	}

	if inlineFragmentExists(l) {
		return parseInlineFragment(l)
	}

	if fragmentSpreadExists(l) {
		return parseFragmentSpread(l)
	}

	panic(errors.New("Expecting a selection"))
}

// https://graphql.github.io/graphql-spec/draft/#Variable
func parseVariable(l *lexer) *variable {
	v := &variable{}

	locStart := l.location().Start()

	if !l.tokenEquals(tokDollar.string()) {
		panic(errors.New("Expecting '$' for varible"))
	} else {
		l.get()

		v.name = *parseName(l)
		v.loc = location{locStart, l.prevLocation().End(), l.source}

		return v
	}
}

// https://graphql.github.io/graphql-spec/draft/#DefaultValue
func parseDefaultValue(l *lexer) *defaultValue {
	dv := &defaultValue{}

	locStart := l.location().Start()

	if !l.tokenEquals(tokEquals.string()) {
		panic(errors.New("Expecting '=' for default value"))
	} else {
		l.get()

		dv.value = parseValue(l)
		dv.loc = location{locStart, l.prevLocation().End(), l.source}

		return dv
	}
}

// https://graphql.github.io/graphql-spec/draft/#Value
func parseValue(l *lexer) value {
	if variableExists(l) {
		return parseVariable(l)
	}

	if intValueExists(l) {
		return parseIntValue(l)
	}

	if floatValueExists(l) {
		return parseFloatValue(l)
	}

	if stringValueExists(l) {
		return parseStringValue(l)
	}

	if booleanValueExists(l) {
		return parseBooleanValue(l)
	}

	if nullValueExists(l) {
		return parseNullValue(l)
	}

	if enumValueExists(l) {
		return parseEnumValue(l)
	}

	if listValueExists(l) {
		return parseListValue(l)
	}

	if objectValueExists(l) {
		return parseObjectValue(l)
	}

	panic(errors.New("No valid value found."))
}

// https://graphql.github.io/graphql-spec/draft/#IntValue
func intValueExists(l *lexer) bool {
	_, err := strconv.ParseInt(l.current().value, 10, 64)
	return err == nil
}

// https://graphql.github.io/graphql-spec/draft/#FloatValue
func floatValueExists(l *lexer) bool {
	_, err := strconv.ParseFloat(l.current().value, 64)
	return err == nil
}

// https://graphql.github.io/graphql-spec/draft/#StringValue
func stringValueExists(l *lexer) bool {
	return singleQuotesStringValueExists(l) || blockStringExists(l)
}

// https://graphql.github.io/graphql-spec/draft/#BooleanValue
func booleanValueExists(l *lexer) bool {
	_, err := strconv.ParseBool(l.current().value)
	return err == nil
}

// https://graphql.github.io/graphql-spec/draft/#NullValue
func nullValueExists(l *lexer) bool {
	return l.current().value == kwNull
}

// https://graphql.github.io/graphql-spec/draft/#EnumValue
func enumValueExists(l *lexer) bool {
	str := l.current().value

	return nameExists(l) && str != kwTrue && str != kwFalse && str != kwNull
}

// https://graphql.github.io/graphql-spec/draft/#ListValue
func listValueExists(l *lexer) bool {
	return l.tokenEquals(tokBracketL.string())
}

// https://graphql.github.io/graphql-spec/draft/#ObjectValue
func objectValueExists(l *lexer) bool {
	return l.tokenEquals(tokBraceL.string())
}

// https://graphql.github.io/graphql-spec/draft/#Variable
func variableExists(l *lexer) bool {
	return string(l.current().value[0]) == tokDollar.string()
}

// https://graphql.github.io/graphql-spec/draft/#Arguments
func parseArguments(l *lexer) *arguments {
	args := &arguments{}

	if !l.tokenEquals(tokParenL.string()) {
		panic(errors.New("Expecting '(' for arguments"))
	} else {
		l.get()

		for !l.tokenEquals(tokParenR.string()) {
			*args = append(*args, parseArgument(l))
		}

		l.get()

		if len(*args) == 0 {
			panic(errors.New("Expecting at least one argument"))
		}

		return args
	}
}

// https://graphql.github.io/graphql-spec/draft/#Argument
func parseArgument(l *lexer) *argument {
	arg := &argument{}

	arg.name = *parseName(l)

	if !l.tokenEquals(tokColon.string()) {
		panic(errors.New("Expecting colon after argument name"))
	}

	l.get()

	arg.value = parseValue(l)

	arg.loc = location{arg.name.Location().Start(), l.prevLocation().End(), l.source}

	return arg
}

// https://graphql.github.io/graphql-spec/draft/#Field
func parseField(l *lexer) *field {
	f := &field{}

	locStart := l.location().Start()

	if aliasExists(l) {
		f.alias = parseAlias(l)
	}

	f.name = *parseName(l)

	if argumentsExist(l) {
		f.arguments = parseArguments(l)
	}

	if directivesExist(l) {
		f.directives = parseDirectives(l)
	}

	if selectionSetExists(l) {
		f.selectionSet = parseSelectionSet(l)
	}

	f.loc = location{locStart, l.prevLocation().End(), l.source}

	return f
}

// https://graphql.github.io/graphql-spec/draft/#FragmentSpread
func parseFragmentSpread(l *lexer) *fragmentSpread {
	fs := &fragmentSpread{}

	locStart := l.location().Start()

	if !l.tokenEquals(tokSpread.string()) {
		panic(errors.New("Expecting '...' operator for a fragment spread"))
	} else {
		l.get()

		fs.fragmentName = *parseFragmentName(l)

		if directivesExist(l) {
			fs.directives = parseDirectives(l)
		}

		fs.loc = location{locStart, l.prevLocation().End(), l.source}

		return fs
	}
}

// https://graphql.github.io/graphql-spec/draft/#InlineFragment
func parseInlineFragment(l *lexer) *inlineFragment {
	inf := &inlineFragment{}

	locStart := l.location().Start()

	if !l.tokenEquals(tokSpread.string()) {
		panic(errors.New("Expecting '...' for an inline fragment"))
	} else {
		l.get()

		if typeConditionExists(l) {
			inf.typeCondition = parseTypeCondition(l)
		}

		if directivesExist(l) {
			inf.directives = parseDirectives(l)
		}

		inf.selectionSet = *parseSelectionSet(l)
		inf.loc = location{locStart, l.prevLocation().End(), l.source}

		return inf
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentName
func parseFragmentName(l *lexer) *name {
	nam := parseName(l)

	if nam.Value() == kwOn {
		panic(errors.New("Fragment name cannot be 'on'"))
	}

	fn := &name{}
	*fn = name(*nam)
	fn.loc = *nam.Location()

	return fn
}

// https://graphql.github.io/graphql-spec/draft/#TypeCondition
func parseTypeCondition(l *lexer) *typeCondition {
	tc := &typeCondition{}

	locStart := l.location().Start()

	if !l.tokenEquals(kwOn) {
		panic(errors.New("Expecting 'on' keyword for a type condition"))
	} else {
		l.get()

		tc.namedType = *parseNamedType(l)
		tc.loc = location{locStart, l.prevLocation().End(), l.source}

		return tc
	}
}

// https://graphql.github.io/graphql-spec/draft/#NamedType
func parseNamedType(l *lexer) *namedType {
	nam := parseName(l)

	nt := new(namedType)

	nt.name = *nam

	return nt
}

// https://graphql.github.io/graphql-spec/draft/#IntValue
func parseIntValue(l *lexer) *intValue {
	iv := &intValue{}

	tok := l.current()

	intVal, err := strconv.ParseInt(tok.value, 10, 64)

	if err != nil {
		panic(errors.Wrap(err, "Failed parsing int value"))
	}

	l.get()

	iv._value = intVal
	iv.loc = location{tok.start, tok.end, l.source}

	return iv
}

// https://graphql.github.io/graphql-spec/draft/#FloatValue
func parseFloatValue(l *lexer) *floatValue {
	fv := &floatValue{}

	tok := l.current()

	floatVal, err := strconv.ParseFloat(tok.value, 64)

	if err != nil {
		panic(errors.Wrap(err, "Failed parsing flot value"))
	}

	l.get()

	fv._value = floatVal
	fv.loc = location{tok.start, tok.end, l.source}

	return fv
}

// https://graphql.github.io/graphql-spec/draft/#StringValue
func parseStringValue(l *lexer) *stringValue {
	tok := l.current()

	sv := &stringValue{}

	if singleQuotesStringValueExists(l) {
		sv._value = *parseSingleQuotesStringValue(l)
	}

	if blockStringExists(l) {
		sv._value = *parseBlockString(l)
	}

	l.get()
	sv._value = tok.value
	sv.loc = location{tok.start, tok.end, l.source}

	return sv
}

// https://graphql.github.io/graphql-spec/draft/#StringValue
func parseSingleQuotesStringValue(l *lexer) *string {
	strVal := new(string)

	*strVal = l.current().value[1 : len(l.current().value)-1]

	if len(*strVal) == 0 {
		return strVal
	}

	if !validateSourceText(*strVal) {
		panic(errors.New("Unsupported characters in a string value"))
	}

	str := *strVal

	for i, _ := range str {
		if str[i] == '\\' {
			if i+1 >= len(str) {
				panic(errors.New("Backslashes are not allowed in a string value"))
			} else {
				if str[i+1] == 'u' {
					if i+5 >= len(str) {
						panic(errors.New("Invalid escaped unicode character in string"))
					} else {
						if !validateEscapedUnicode(str[i+1 : i+6]) {
							panic(errors.New(("Invalid escaped unicode character in string")))
						}
					}
				} else {
					switch str[i+1] {
					case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
					default:
						panic(errors.New("Invalid escaped character in string"))
					}
				}
			}
		}
	}

	return strVal
}

// https://graphql.github.io/graphql-spec/draft/#StringValue
func parseBlockString(l *lexer) *string {
	strVal := new(string)

	*strVal = l.current().value[3 : len(l.current().value)-3]

	str := *strVal

	if validateSourceText(str) &&
		!strings.Contains(str, "\"\"\"") &&
		!strings.Contains(str, "\\\"\"\"") {
		return strVal
	}

	if str == "\\\"\"\"" {
		return strVal
	}

	panic(errors.New("Invalid characters in block string"))
}

// https://graphql.github.io/graphql-spec/draft/#SourceCharacter
func validateSourceText(str string) bool {
	reg, err := regexp.Compile("[\u0009\u000A\u000D\u0020-\uFFFF]*")

	if err != nil {
		panic(err)
	}

	return reg.MatchString(str)
}

func validateEscapedUnicode(str string) bool {
	if str[0] != 'u' {
		return false
	}

	str = str[1:len(str)]

	reg, err := regexp.Compile("/[0-9A-Fa-f]{4}/")

	if err != nil {
		panic(err)
	}

	return reg.MatchString(str)
}

// https://graphql.github.io/graphql-spec/draft/#BooleanValue
func parseBooleanValue(l *lexer) *booleanValue {
	bv := &booleanValue{}

	tok := l.current()

	boolVal, err := strconv.ParseBool(tok.value)

	if err != nil {
		panic(errors.Wrap(err, "Failed parsing bool value"))
	}

	l.get()

	bv._value = boolVal
	bv.loc = location{tok.start, tok.end, l.source}

	return bv
}

// https://graphql.github.io/graphql-spec/draft/#NullValue
func parseNullValue(l *lexer) *nullValue {
	nv := &nullValue{}

	tok := l.current()

	if tok.value != kwNull {
		panic(errors.New("Expecting 'null' keyword"))
	} else {
		l.get()

		nv.loc = location{tok.start, tok.end, l.source}

		return nv
	}
}

// https://graphql.github.io/graphql-spec/draft/#EnumValue
func parseEnumValue(l *lexer) *enumValue {
	nam := parseName(l)

	switch nam.Value() {
	case kwTrue, kwFalse, kwNull:
		panic(errors.New("Enum value cannot be 'true', 'false' or 'null'"))
	default:
		nv := &enumValue{}

		nv.name = *nam
		nv.loc = location{nam.Location().Start(), nam.Location().End(), l.source}

		return nv
	}
}

// https://graphql.github.io/graphql-spec/draft/#ListValue
func parseListValue(l *lexer) *listValue {
	lv := &listValue{}

	locStart := l.location().Start()

	if !l.tokenEquals(tokBracketL.string()) {
		panic(errors.New("Expecting '[' for a list value"))
	} else {
		l.get()

		for !l.tokenEquals(tokBracketR.string()) {
			lv.values = append(lv.values, parseValue(l))
		}

		l.get()

		lv.loc = location{locStart, l.prevLocation().End(), l.source}

		return lv
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectValue
func parseObjectValue(l *lexer) *objectValue {
	ov := &objectValue{}

	locStart := l.location().Start()

	if !l.tokenEquals(tokBraceL.string()) {
		panic(errors.New("Expecting '{' for an object value"))
	} else {
		l.get()

		for !l.tokenEquals(tokBraceR.string()) {
			ov.values = append(ov.values, *parseObjectField(l))
		}

		l.get()

		ov.loc = location{locStart, l.prevLocation().End(), l.source}

		return ov
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectField
func parseObjectField(l *lexer) *objectField {
	of := &objectField{}

	of.name = *parseName(l)

	if !l.tokenEquals(tokColon.string()) {
		panic(errors.New("Expecting color after object field name"))
	}

	l.get()

	of._value = parseValue(l)
	of.loc = location{of.name.Location().Start(), l.prevLocation().End(), l.source}

	return of
}

// https://graphql.github.io/graphql-spec/draft/#Alias
func parseAlias(l *lexer) *alias {
	a := &alias{}
	a.value = parseName(l).Value()

	if !l.tokenEquals(tokColon.string()) {
		panic(errors.New("Expecting colon after alias name"))
	}

	l.get()
	return a
}

// https://graphql.github.io/graphql-spec/draft/#TypeCondition
func typeConditionExists(l *lexer) bool {
	return l.tokenEquals(kwOn)
}

// https://graphql.github.io/graphql-spec/draft/#Alias
func aliasExists(l *lexer) bool {
	return l.tokens[l.currentTokenIndex+1].value == tokColon.string()
}

// https://graphql.github.io/graphql-spec/draft/#SelectionSet
func selectionSetExists(l *lexer) bool {
	return l.tokenEquals(tokBraceL.string())
}

// https://graphql.github.io/graphql-spec/draft/#Field
func fieldExists(l *lexer) bool {
	return nameExists(l)
}

// https://graphql.github.io/graphql-spec/draft/#FragmentSpread
func fragmentSpreadExists(l *lexer) bool {
	return l.tokenEquals(tokSpread.string())
}

// https://graphql.github.io/graphql-spec/draft/#InlineFragment
func inlineFragmentExists(l *lexer) bool {
	if !l.tokenEquals(tokSpread.string()) {
		return false
	}

	tmpLexer := *l
	tmpLexer.get()

	if typeConditionExists(&tmpLexer) {
		parseTypeCondition(&tmpLexer)
	}

	if directivesExist(&tmpLexer) {
		parseDirectives(&tmpLexer)
	}

	return selectionSetExists(&tmpLexer)
}

// https://graphql.github.io/graphql-spec/draft/#Arguments
func argumentsExist(l *lexer) bool {
	return l.tokenEquals(tokParenL.string())
}

// https://graphql.github.io/graphql-spec/draft/#NamedType
func namedTypeExists(l *lexer) bool {
	return nameExists(l)
}

// https://graphql.github.io/graphql-spec/draft/#ListType
func listTypeExists(l *lexer) bool {
	return l.tokenEquals(tokBracketL.string())
}

// https://graphql.github.io/graphql-spec/draft/#NonNullType
func nonNullTypeExists(l *lexer) bool {
	return l.tokenEquals(tokBang.string())
	//return (namedTypeExists(l) || listTypeExists(l)) && l.tokens[l.currentTokenIndex+1].value == tokBang.string()
}

// https://graphql.github.io/graphql-spec/draft/#Name
func nameExists(l *lexer) bool {
	tok := l.current()

	pattern := "^[_A-Za-z][_0-9A-Za-z]*$"

	// If the current token is not a Name, return nil
	if tok.kind != tokName {
		return false
	}

	// Check if the given name matches the regex provided by graphql spec at
	// https://graphql.github.io/graphql-spec/draft/#Name
	match, err := regexp.MatchString(pattern, tok.value)
	if err != nil {
		return false
	}

	// If the name does not match the requirements, return an error.
	if !match {
		return false
	}

	return true
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinitions
func variableDefinitionsExist(l *lexer) bool {
	return l.tokenEquals(tokParenL.string())
}

// https://graphql.github.io/graphql-spec/draft/#RootOperationTypeDefinition
func rootOperationTypeDefinitionsExist(l *lexer) bool {
	return operationTypeExists(l)
}

// https://graphql.github.io/graphql-spec/draft/#OperationType
func operationTypeExists(l *lexer) bool {
	return l.tokenEquals(kwQuery) || l.tokenEquals(kwMutation) || l.tokenEquals(kwSubscription)
}

// https://graphql.github.io/graphql-spec/draft/#TypeExtension
func typeExtensionExists(l *lexer) bool {
	return (scalarTypeExtensionExists(l) ||
		objectTypeExtensionExists(l) ||
		interfaceTypeExtensionExists(l) ||
		unionTypeExtensionExists(l) ||
		enumTypeExtensionExists(l) ||
		inputObjectTypeExtensionExist(l))
}

// http://spec.graphql.org/draft/#InputObjectTypeExtension
func inputObjectTypeExtensionExist(l *lexer) bool {
	return l.tokenEquals(kwExtend, kwInput)
}

// https://graphql.github.io/graphql-spec/draft/#SchemaExtension
func schemaExtensionExists(l *lexer) bool {
	return l.tokenEquals(kwExtend, kwSchema)
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDirectiveLocation
func executableDirectiveLocationExists(l *lexer) bool {
	for _, edl := range executableDirectiveLocations {
		if l.tokenEquals((string)(edl)) {
			return true
		}
	}

	return false
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDirectiveLocation
func typeSystemDirectiveLocationExists(l *lexer) bool {
	for _, tsdl := range typeSystemDirectiveLocations {
		if l.tokenEquals((string)(tsdl)) {
			return true
		}
	}

	return false
}

// https://graphql.github.io/graphql-spec/draft/#ScalarTypeExtension
func scalarTypeExtensionExists(l *lexer) bool {
	return l.tokenEquals(kwExtend, kwScalar)
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeExtension
func objectTypeExtensionExists(l *lexer) bool {
	return l.tokenEquals(kwExtend, kwType)
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeExtension
func interfaceTypeExtensionExists(l *lexer) bool {
	return l.tokenEquals(kwExtend, kwInterface)
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeExtension
func unionTypeExtensionExists(l *lexer) bool {
	return l.tokenEquals(kwExtend, kwUnion)
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeExtension
func enumTypeExtensionExists(l *lexer) bool {
	return l.tokenEquals(kwExtend, kwEnum)
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeExtension
func inputObjectTypeExtensionExists(l *lexer) bool {
	return l.tokenEquals(kwExtend, kwInput)
}

// https://graphql.github.io/graphql-spec/draft/#InputFieldsDefinition
func inputFieldsDefinitionExists(l *lexer) bool {
	return l.tokenEquals(tokParenL.string())
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinition
func enumValuesDefinitionExist(l *lexer) bool {
	return l.tokenEquals(tokBracketL.string())
}

// https://graphql.github.io/graphql-spec/draft/#UnionMemberTypes
func unionMemberTypesExist(l *lexer) bool {
	return l.tokenEquals(tokEquals.string())
}

// https://graphql.github.io/graphql-spec/draft/#DefaultValue
func defaultValueExists(l *lexer) bool {
	return l.tokenEquals(tokEquals.string())
}

// https://graphql.github.io/graphql-spec/draft/#ArgumentsDefinition
func argumentsDefinitionExist(l *lexer) bool {
	return l.tokenEquals(tokParenL.string())
}

// https://graphql.github.io/graphql-spec/draft/#ImplementsInterfaces
func implementsInterfacesExists(l *lexer) bool {
	return l.tokenEquals(kwImplements)
}

// https://graphql.github.io/graphql-spec/draft/#FieldsDefinition
func fieldsDefinitionExists(l *lexer) bool {
	return l.tokenEquals(tokBraceL.string())
}

// https://graphql.github.io/graphql-spec/draft/#Description
func descriptionExists(l *lexer) bool {
	return stringValueExists(l)
}

// https://graphql.github.io/graphql-spec/draft/#ScalarTypeDefinition
func scalarTypeDefinitionExists(l *lexer) bool {
	return descriptionExists(l) || l.tokenEquals(kwScalar)
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeDefinition
func objectTypeDefinitionExists(l *lexer) bool {
	return descriptionExists(l) || l.tokenEquals(kwType)
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeDefinition
func interfaceTypeDefinitionExists(l *lexer) bool {
	return descriptionExists(l) || l.tokenEquals(kwInterface)
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeDefinition
func unionTypeDefinitionExists(l *lexer) bool {
	return descriptionExists(l) || l.tokenEquals(kwUnion)
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeDefinition
func enumTypeDefinitionExists(l *lexer) bool {
	return descriptionExists(l) || l.tokenEquals(kwEnum)
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeDefinition
func inputObjectTypeDefinitionExists(l *lexer) bool {
	return descriptionExists(l) || l.tokenEquals(kwInput)
}

// https://graphql.github.io/graphql-spec/draft/#Directives
func directivesExist(l *lexer) bool {
	return l.tokenEquals(tokAt.string())
}

// https://graphql.github.io/graphql-spec/draft/#SchemaDefinition
func schemaDefinitionExists(l *lexer) bool {
	return l.tokenEquals(kwSchema)
}

// https://graphql.github.io/graphql-spec/draft/#TypeDefinition
func typeDefinitionExists(l *lexer) bool {
	return scalarTypeDefinitionExists(l) || objectTypeDefinitionExists(l) || interfaceTypeDefinitionExists(l) ||
		unionTypeDefinitionExists(l) || enumTypeDefinitionExists(l) || inputObjectTypeDefinitionExists(l)
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveDefinition
func directiveDefinitionExists(l *lexer) bool {
	return descriptionExists(l) || l.tokenEquals(kwDirective)
}

// https://graphql.github.io/graphql-spec/draft/#OperationDefinition
func operationDefinitionExists(l *lexer) bool {
	return operationTypeExists(l) || selectionSetExists(l)
}

// https://graphql.github.io/graphql-spec/draft/#FragmentDefinition
func fragmentDefinitionExists(l *lexer) bool {
	return l.tokenEquals(kwFragment)
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDefinition
func executableDefinitionExists(l *lexer) bool {
	return operationDefinitionExists(l) || fragmentDefinitionExists(l)
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDefinition
func typeSystemDefinitionExists(l *lexer) bool {
	return schemaDefinitionExists(l) || typeDefinitionExists(l) || directiveDefinitionExists(l)
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemExtension
func typeSystemExtensionExists(l *lexer) bool {
	return schemaExtensionExists(l) || typeExtensionExists(l)
}

func singleQuotesStringValueExists(l *lexer) bool {
	return l.current().kind == tokString
}

func blockStringExists(l *lexer) bool {
	return l.current().kind == tokBlockString
}

func getUnexpected(l *lexer) string {
	return l.current().value
}
