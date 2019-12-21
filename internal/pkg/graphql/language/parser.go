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

	pDoc := parseDocument(l)

	// recover syntax errors
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	return pDoc, nil
}

// ! Redo error management
// https://graphql.github.io/graphql-spec/draft/#Document
func parseDocument(l *lexer) *document {
	doc := &document{}

	def := parseDefinitions(l)

	doc.Definitions = *def

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

func executableDefinitionExists(l *lexer) bool {
	return false
}

func typeSystemDefinitionExists(l *lexer) bool {
	return false
}

func typeSystemExtensionExists(l *lexer) bool {
	return false
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

func operationDefinitionExists(l *lexer) bool {
	return false
}

func fragmentDefinitionExists(l *lexer) bool {
	return false
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

func schemaDefinitionExists(l *lexer) bool {
	return false
}

func typeDefinitionExists(l *lexer) bool {
	return false
}

func directiveDefinitionExists(l *lexer) bool {
	return false
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
		schDef.Directives = parseDirectives(l)
	}

	if !l.tokenEquals(tokBraceL.string()) {
		panic(errors.New("Missing '{' for a schema definition"))
	}

	l.get()

	schDef.RootOperationTypeDefinitions = *parseRootOperationTypeDefinitions(l)

	if !l.tokenEquals(tokBraceR.string()) {
		panic(errors.New("Missing '}' for schema definition"))
	}

	locEnd := l.current().end

	l.get()

	schDef.Loc = location{locStart, locEnd, l.source}

	return schDef
}

func directivesExist(l *lexer) bool {
	return false
}

// https://graphql.github.io/graphql-spec/draft/#RootOperationTypeDefinition
func parseRootOperationTypeDefinitions(l *lexer) *rootOperationTypeDefinitions {
	rotds := &rootOperationTypeDefinitions{}

	for !l.tokenEquals(tokBraceR.string()) {
		rotd := parseRootOperationTypeDefinition(l)

		*rotds = append(*rotds, *rotd)
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

	rotd.OperationType = *parseOperationType(l)

	if !l.tokenEquals(tokColon.string()) {
		panic(errors.New("Expecting ':' after operation type"))
	}

	l.get()

	rotd.NamedType = *parseNamedType(l)

	rotd.Loc = location{locStart, rotd.NamedType.Location().End, l.source}

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

func scalarTypeDefinitionExists(l *lexer) bool {
	return false
}

func objectTypeDefinitionExists(l *lexer) bool {
	return false
}

func interfaceTypeDefinitionExists(l *lexer) bool {
	return false
}

func unionTypeDefinitionExists(l *lexer) bool {
	return false
}

func enumTypeDefinitionExists(l *lexer) bool {
	return false
}

func inputObjectTypeDefinitionExists(l *lexer) bool {
	return false
}

// https://graphql.github.io/graphql-spec/draft/#ScalarTypeDefinition
func parseScalarTypeDefinition(l *lexer) *scalarTypeDefinition {
	stp := &scalarTypeDefinition{}

	if descriptionExists(l) {
		stp.Description = parseDescription(l)
	}

	if !l.tokenEquals(tsdlScalar) {
		panic(errors.New("Missing 'scalar' keyword for scalar type definition"))
	}

	tok := l.get()

	stp.Name = *parseName(l)

	if directivesExist(l) {
		stp.Directives = parseDirectives(l)
	}

	stp.Loc = location{tok.start, l.prevLocation().End, l.source}

	return stp
}

func descriptionExists(l *lexer) bool {
	return false
}

// https://graphql.github.io/graphql-spec/draft/#Description
func parseDescription(l *lexer) *description {
	strVal := parseStringValue(l)

	return (*description)(strVal)
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeDefinition
func parseObjectTypeDefinition(l *lexer) *objectTypeDefinition {
	otd := &objectTypeDefinition{}

	locStart := l.location().Start

	if descriptionExists(l) {
		otd.Description = parseDescription(l)
	}

	if !l.tokenEquals(kwType) {
		panic(errors.New("Expecting 'type' keyword for object type definition"))
	}

	otd.Name = *parseName(l)

	if implementsInterfacesExists(l) {
		otd.ImplementsInterfaces = parseImplementsInterfaces(l)
	}

	if directivesExist(l) {
		otd.Directives = parseDirectives(l)
	}

	if fieldsDefinitionExists(l) {
		otd.FieldsDefinition = parseFieldsDefinition(l)
	}

	otd.Loc = location{locStart, l.prevLocation().End, l.source}

	return otd
}

func implementsInterfacesExists(l *lexer) bool {
	return false
}

func fieldsDefinitionExists(l *lexer) bool {
	return false
}

// https://graphql.github.io/graphql-spec/draft/#ImplementsInterfaces
func parseImplementsInterfaces(l *lexer) *implementsInterfaces {
	ii := &implementsInterfaces{}

	if !l.tokenEquals(kwImplements) {
		panic(errors.New("Expecting 'implements' keyword"))
	}

	if l.tokenEquals(tokAmp.string()) {
		l.get()
	}

	*ii = append(*ii, *parseNamedType(l))

	for l.tokenEquals(tokAmp.string()) {
		l.get()

		(*ii) = append(*ii, *parseNamedType(l))
	}

	return ii
}

// https://graphql.github.io/graphql-spec/draft/#FieldsDefinition
func parseFieldsDefinition(l *lexer) *fieldsDefinition {
	fd := &fieldsDefinition{}

	if l.tokenEquals(tokBraceL.string()) {
		panic(errors.New("Expecting '{' for fields definition"))
	}

	l.get()
	/// !here
	ret = &fieldsDefinition{}

	for !l.tokenEquals(tokBraceR.string()) {
		var fd *fieldDefinition

		fd, err = parseFieldDefinition(l)

		if err != nil {
			ret = nil
			return
		}

		(*ret) = append(*ret, *fd)
	}

	l.get()

	if len(*ret) == 0 {
		err = errors.New("Expecting at lease one field definition")
		ret = nil
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#FieldsDefinition
func parseFieldDefinition(l *lexer) *fieldDefinition {

	locStart := l.location().Start

	desc, _ := parseDescription(l)

	nam, err := parseName(l)

	if err != nil {
		return
	}

	argsDef, _ := parseArgumentsDefinition(l)

	if !l.tokenEquals(tokColon.string()) {
		err = errors.New("Expecting ':' for field definition")
		return
	}

	l.get()

	_typ, err := parseType(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	ret = &fieldDefinition{}

	ret.Description = desc
	ret.Name = *nam
	ret.ArgumentsDefinition = argsDef
	ret.Type = _typ
	ret.Directives = dirs
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#ArgumentsDefinition
func parseArgumentsDefinition(l *lexer) *argumentsDefinition {

	if !l.tokenEquals(tokParenL.string()) {
		err = errors.New("Expecting '(' for arguments definition")
		return
	}

	l.get()

	ret = &argumentsDefinition{}

	for !l.tokenEquals(tokParenR.string()) {
		var ivDef *inputValueDefinition

		ivDef, err = parseInputValueDefinition(l)

		if err != nil {
			ret = nil
			return
		}

		*ret = append(*ret, *ivDef)
	}

	l.get()

	if len(*ret) == 0 {
		ret = nil
		err = errors.New("Expecting at least one input value definitions")
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#InputValueDefinition
func parseInputValueDefinition(l *lexer) *inputValueDefinition {

	locStart := l.location().Start

	desc, _ := parseDescription(l)

	nam, err := parseName(l)

	if err != nil {
		return
	}

	if !l.tokenEquals(tokColon.string()) {
		err = errors.New("Expecting ':' for input value definition")
		return
	}

	l.get()

	_typ, err := parseType(l)

	if err != nil {
		return
	}

	defVal, _ := parseDefaultValue(l)

	dirs, _ := parseDirectives(l)

	ret = &inputValueDefinition{}

	ret.Description = desc
	ret.Name = *nam
	ret.Type = _typ
	ret.DefaultValue = defVal
	ret.Directives = dirs
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeDefinition
func parseInterfaceTypeDefinition(l *lexer) *interfaceTypeDefinition {

	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(tsdlInterface) {
		err = errors.New("Expecting 'interface' keyword for interface type definition")
		return
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	fds, _ := parseFieldsDefinition(l)

	ret = &interfaceTypeDefinition{}

	ret.Description = desc
	ret.Directives = dirs
	ret.FieldsDefinition = fds
	ret.Name = *nam
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeDefinition
func parseUnionTypeDefinition(l *lexer) *unionTypeDefinition {

	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(tsdlUnion) {
		err = errors.New("Expecting 'union' keyowrd for union type definition")
		return
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	umt, _ := parseUnionMemberTypes(l)

	ret = &unionTypeDefinition{}

	ret.Description = desc
	ret.Name = *nam
	ret.Directives = dirs
	ret.UnionMemberTypes = umt
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#UnionMemberTypes
func parseUnionMemberTypes(l *lexer) *unionMemberTypes {

	if !l.tokenEquals(tokEquals.string()) {
		err = errors.New("Expecting '=' for union member types")
		return
	}

	l.get()

	if l.tokenEquals(tokPipe.string()) {
		l.get()
	}

	nt, err := parseNamedType(l)

	if err != nil {
		return
	}

	ret = &unionMemberTypes{}

	*ret = append(*ret, *nt)

	for l.tokenEquals(tokPipe.string()) {
		l.get()

		nt, err = parseNamedType(l)

		if err != nil {
			return nil, err
		}

		*ret = append(*ret, *nt)
	}

	if len(*ret) == 0 {
		ret = nil
		err = errors.New("Expecting at least one union member type")
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeDefinition
func parseEnumTypeDefinition(l *lexer) *enumTypeDefinition {

	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(tsdlEnum) {
		err = errors.New("Expecting 'enum' keyword for enum type definition")
		return
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	evd, _ := parseEnumValuesDefinition(l)

	ret = &enumTypeDefinition{}

	ret.Description = desc
	ret.Name = *nam
	ret.Directives = dirs
	ret.EnumValuesDefinition = evd
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinition
func parseEnumValuesDefinition(l *lexer) *enumValuesDefinition {

	if !l.tokenEquals(tokBraceL.string()) {
		err = errors.New("Expecting '{' for enum values definition")
		return
	}

	l.get()

	ret = &enumValuesDefinition{}

	for !l.tokenEquals(tokBraceR.string()) {
		var evd *enumValueDefinition

		evd, err = parseEnumValueDefinition(l)

		if err != nil {
			ret = nil
			return
		}

		*ret = append(*ret, *evd)
	}

	l.get()

	if len(*ret) == 0 {
		err = errors.New("Expecting at least one enum value definition")
		ret = nil
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#EnumValuesDefinition
func parseEnumValueDefinition(l *lexer) *enumValueDefinition {

	locStart := l.location().Start

	desc, _ := parseDescription(l)

	ev, err := parseEnumValue(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	ret = &enumValueDefinition{}

	ret.Description = desc
	ret.EnumValue = *ev
	ret.Directives = dirs
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeDefinition
func parseInputObjectTypeDefinition(l *lexer) *inputObjectTypeDefinition {

	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(kwInput) {
		err = errors.New("Expecting 'input' keyword for input object type definition")
		return
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	ifds, _ := parseInputFieldsDefinition(l)

	ret = &inputObjectTypeDefinition{}

	ret.Description = desc
	ret.Directives = dirs
	ret.Name = *nam
	ret.InputFieldsDefinition = ifds
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return ret, nil
}

// https://graphql.github.io/graphql-spec/draft/#InputFieldsDefinition
func parseInputFieldsDefinition(l *lexer) *inputFieldsDefinition {

	if !l.tokenEquals(tokBraceL.string()) {
		err = errors.New("Expecting '{' for input fields definition")
		return
	}

	l.get()

	ret = &inputFieldsDefinition{}

	for !l.tokenEquals(tokBraceR.string()) {
		var ivd *inputValueDefinition

		ivd, err = parseInputValueDefinition(l)

		if err != nil {
			return
		}

		*ret = append(*ret, *ivd)
	}

	l.get()

	if len(*ret) == 0 {
		ret = nil
		err = errors.New("Expecting at least one input field definition")
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveDefinition
func parseDirectiveDefinition(l *lexer) *directiveDefinition {

	locStart := l.location().Start

	desc, _ := parseDescription(l)

	if !l.tokenEquals(kwDirective) {
		err = errors.New("Expecting 'directive' keyword for directive definition")
		return
	}

	l.get()

	if !l.tokenEquals(tokAt.string()) {
		err = errors.New("Expecting '@' for directive definition")
		return
	}

	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	argsDef, _ := parseArgumentsDefinition(l)

	if !l.tokenEquals(kwOn) {
		err = errors.New("Expecting 'on' keyworkd for directive definition")
		return
	}

	l.get()

	dls, err := parseDirectiveLocations(l)

	if err != nil {
		return
	}

	ret = &directiveDefinition{}

	ret.Description = desc
	ret.Name = *nam
	ret.ArgumentsDefinition = argsDef
	ret.DirectiveLocations = *dls
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocations
func parseDirectiveLocations(l *lexer) *directiveLocations {

	ret = &directiveLocations{}

	if l.tokenEquals(tokPipe.string()) {
		l.get()
	}

	dl, err := parseDirectiveLocation(l)

	if err != nil {
		ret = nil
		return
	}

	*ret = append(*ret, *dl)

	for l.tokenEquals(tokPipe.string()) {
		l.get()

		var dl *directiveLocation

		dl, err = parseDirectiveLocation(l)

		if err != nil {
			return
		}

		*ret = append(*ret, *dl)
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#TypeExtension
func parseTypeExtension(l *lexer) typeExtension {

	ret, _ = parseScalarTypeExtension(l)

	if ret != nil {
		return
	}

	ret, _ = parseObjectTypeExtension(l)

	if ret != nil {
		return
	}

	ret, _ = parseInterfaceTypeExtension(l)

	if ret != nil {
		return
	}

	ret, _ = parseUnionTypeExtension(l)

	if ret != nil {
		return
	}

	ret, _ = parseEnumTypeExtension(l)

	if ret != nil {
		return
	}

	ret, _ = parseInputObjectTypeExtension(l)

	if ret != nil {
		return
	}

	err = errors.New("Expecting type extension")
	return
}

// https://graphql.github.io/graphql-spec/draft/#ScalarTypeExtension
func parseScalarTypeExtension(l *lexer) *scalarTypeExtension {

	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwScalar) {
		err = errors.New("Expecting 'extend scalar' keywords for scalar type extension")
		return
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	dirs, err := parseDirectives(l)

	if err != nil {
		return
	}

	ret = &scalarTypeExtension{}

	ret.Name = *nam
	ret.Directives = *dirs
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#ObjectTypeExtension
func parseObjectTypeExtension(l *lexer) *objectTypeExtension {

	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwType) {
		err = errors.New("Expecting 'extend type' keywords for object type extension")
		return
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	ii, _ := parseImplementsInterfaces(l)

	dirs, _ := parseDirectives(l)

	fds, _ := parseFieldsDefinition(l)

	if ii == nil && dirs == nil && fds == nil {
		err = errors.New("Expecting at least one of 'implements interface', 'directives', 'fields definition' for object type extension")
		return
	}

	ret = &objectTypeExtension{}

	ret.Name = *nam
	ret.ImplementsInterfaces = ii
	ret.Directives = dirs
	ret.FieldsDefinition = fds
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#InterfaceTypeExtension
func parseInterfaceTypeExtension(l *lexer) *interfaceTypeExtension {

	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwInterface) {
		err = errors.New("Expecting 'extend interface' keywords for interface type extension")
		return
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	fds, _ := parseFieldsDefinition(l)

	if dirs == nil && fds == nil {
		err = errors.New("Expecting at least one of 'directives', 'fields definition' for interface type extension")
		return
	}

	ret = &interfaceTypeExtension{}

	ret.Name = *nam
	ret.Directives = dirs
	ret.FieldsDefinition = fds
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#UnionTypeExtension
func parseUnionTypeExtension(l *lexer) *unionTypeExtension {

	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, tsdlUnion) {
		err = errors.New("Expecting 'extend union' keywords for union type extension")
		return
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	umt, _ := parseUnionMemberTypes(l)

	if dirs == nil && umt == nil {
		err = errors.New("Expecting at  least one of 'directives', 'union member types' for union type extension")
		return
	}

	ret = &unionTypeExtension{}

	ret.Name = *nam
	ret.Directives = dirs
	ret.UnionMemberTypes = umt
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#EnumTypeExtension
func parseEnumTypeExtension(l *lexer) *enumTypeExtension {

	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, tsdlEnum) {
		err = errors.New("Expecting 'extend enum' keywords for enum type extension")
		return
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	evd, _ := parseEnumValuesDefinition(l)

	if dirs == nil && evd == nil {
		err = errors.New("Expecting at least one of 'directives', 'enum values definition' for enum type extension")
		return
	}

	ret = &enumTypeExtension{}

	ret.Name = *nam
	ret.Directives = dirs
	ret.EnumValuesDefinition = evd
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#DirectiveLocation
func parseDirectiveLocation(l *lexer) *directiveLocation {

	edl, err := parseExecutableDirectiveLocation(l)

	if err == nil {
		ret = (*directiveLocation)(edl)
		return
	}

	tsdl, err := parseTypeSystemDirectiveLocation(l)

	if err != nil {
		err = errors.Wrap(err, "Expecting a directive location")
		ret = nil
		return
	}

	ret = (*directiveLocation)(tsdl)

	return
}

// https://graphql.github.io/graphql-spec/draft/#ExecutableDirectiveLocation
func parseExecutableDirectiveLocation(l *lexer) *executableDirectiveLocation {

	tok := l.current()

	for i := range executableDirectiveLocations {
		if string(executableDirectiveLocations[i]) == tok.value {
			l.get()

			ret = &executableDirectiveLocations[i]

			return
		}
	}

	err = errors.New("Expecting executable directive location")
	return
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemDirectiveLocation
func parseTypeSystemDirectiveLocation(l *lexer) *typeSystemDirectiveLocation {

	tok := l.current()

	for i := range typeSystemDirectiveLocations {
		if string(typeSystemDirectiveLocations[i]) == tok.value {
			l.get()

			ret = &typeSystemDirectiveLocations[i]

			return
		}
	}

	err = errors.New("Expecting type systen directive location")
	return
}

// https://graphql.github.io/graphql-spec/draft/#TypeSystemExtension
func parseTypeSystemExtension(l *lexer) typeSystemExtension {

	ret, err = parseSchemaExtension(l)

	if err != nil {
		return
	}

	ret, err = parseTypeExtension(l)

	if err != nil {
		err = errors.Wrap(err, "Expecting type system extension")
		ret = nil
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#SchemaExtension
func parseSchemaExtension(l *lexer) *schemaExtension {

	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwSchema) {
		err = errors.New("Expecting 'extend schema' keywords for schema extension")
		return
	}

	l.get()
	l.get()

	dirs, _ := parseDirectives(l)

	if !l.tokenEquals(tokBraceL.string()) {
		err = errors.New("Expecting '{' for schema extension")
		return
	}

	l.get()

	rotds, _ := parseRootOperationTypeDefinitions(l)

	if !l.tokenEquals(tokBraceR.string()) {
		err = errors.New("Expecting '}' for schema extension")
		return
	}

	l.get()

	if dirs == nil && rotds == nil {
		err = errors.New("Expecting directives or root operation type definitions for schema extension")
		return
	}

	ret = &schemaExtension{}

	ret.Directives = dirs
	ret.RootOperationTypeDefinitions = rotds
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#InputObjectTypeExtension
func parseInputObjectTypeExtension(l *lexer) *inputObjectTypeExtension {

	locStart := l.location().Start

	if !l.tokenEquals(kwExtend, kwInput) {
		err = errors.New("Expecting 'extend' keyword for input object type extension")
		return
	}

	l.get()
	l.get()

	nam, err := parseName(l)

	if err != nil {
		return
	}

	dirs, _ := parseDirectives(l)

	idfs, _ := parseInputFieldsDefinition(l)

	if dirs == nil && idfs == nil {
		err = errors.New("Expecting at lease one of 'directives', 'input fields definition' fo input object type extension")
		return
	}

	ret = &inputObjectTypeExtension{}

	ret.Name = *nam
	ret.Directives = dirs
	ret.InputFieldsDefinition = idfs
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#OperationDefinition
func parseOperationDefinition(l *lexer) *operationDefinition {

	locStart := l.location().Start

	// Shorthand query
	// https://graphql.github.io/graphql-spec/draft/#sec-Language.Operations.Query-shorthand
	if l.tokenEquals(tokBraceL.string()) {
		var shorthandQuery *selectionSet

		shorthandQuery, err = parseSelectionSet(l)

		if err != nil {
			return
		}

		ret = &operationDefinition{}

		ret.OperationType = kwQuery
		ret.SelectionSet = *shorthandQuery
		ret.Loc = location{locStart, l.prevLocation().End, l.source}

		return
	} else if !l.tokenEquals(kwQuery) &&
		!l.tokenEquals(kwMutation) &&
		!l.tokenEquals(kwSubscription) {
		err = errors.New("Expecting one of 'query', 'mutation', 'subscription' for operation definition")
		return
	} else {
		tok := l.get()

		opType := tok.value

		nam, _ := parseName(l)

		varDef, _ := parseVariableDefinitions(l)

		directives, _ := parseDirectives(l)

		var selSet *selectionSet

		selSet, err = parseSelectionSet(l)

		if err != nil {
			return
		}

		ret = &operationDefinition{}

		ret.OperationType = operationType(opType)
		ret.Name = nam
		ret.VariableDefinitions = varDef
		ret.Directives = directives
		ret.SelectionSet = *selSet
		ret.Loc = location{locStart, tok.end, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentDefinition
func parseFragmentDefinition(l *lexer) *fragmentDefinition {

	locStart := l.location().Start

	if !l.tokenEquals(kwFragment) {
		err = errors.New("Expecting fragment keyword")
		return
	} else {
		l.get()

		var nam *fragmentName

		nam, err = parseFragmentName(l)

		if err != nil {
			return
		}

		if nam.Value == kwOn {
			err = errors.New("Fragment name cannot be 'on'")
			return
		}

		var typeCond *typeCondition

		typeCond, err = parseTypeCondition(l)

		if err != nil {
			return
		}

		directives, _ := parseDirectives(l)

		var selectionSet *selectionSet

		selectionSet, err = parseSelectionSet(l)

		if err != nil {
			return
		}

		ret = &fragmentDefinition{}

		ret.FragmentName = *nam
		ret.TypeCondition = *typeCond
		ret.Directives = directives
		ret.SelectionSet = *selectionSet
		ret.Loc = location{locStart, l.prevLocation().End, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#Name
func parseName(l *lexer) *name {

	tok := l.current()

	pattern := "^[_A-Za-z][_0-9A-Za-z]*$"

	// If the current token is not a Name, return nil
	if tok.kind != tokName {
		err = errors.New("Not a name")
		return
	}

	// Check if the given name matches the regex provided by graphql spec at
	// https://graphql.github.io/graphql-spec/draft/#Name
	match, err := regexp.MatchString(pattern, tok.value)
	if err != nil {
		err = errors.Wrap(err, "failed to parse name: ")
		return
	}

	// If the name does not match the requirements, return an error.
	if !match {
		err = errors.New("invalid name - " + tok.value)
		return
	}

	l.get()

	ret = &name{}

	// Populate the Name struct.
	ret.Value = tok.value
	ret.Loc.Start = tok.start
	ret.Loc.End = tok.end
	ret.Loc.Source = l.source

	// Return the AST Name object.
	return
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinition
func parseVariableDefinitions(l *lexer) *variableDefinitions {

	if !l.tokenEquals(tokParenL.string()) {
		err = errors.New("Expecting '(' opener for variable definitions")
		return
	} else {
		l.get()

		ret = &variableDefinitions{}

		for !l.tokenEquals(tokParenR.string()) {
			var varDef *variableDefinition

			varDef, err = parseVariableDefinition(l)

			if err != nil {
				return
			}

			*ret = append(*ret, *varDef)
		}

		l.get()

		if len(*ret) == 0 {
			err = errors.New("Expecting at least one variable definition")
			return
		}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#VariableDefinition
func parseVariableDefinition(l *lexer) *variableDefinition {

	locStart := l.location().Start

	_var, err := parseVariable(l)

	if err != nil {
		return
	}

	if !l.tokenEquals(tokColon.string()) {
		err = errors.New("Expecting a colon after variable name")
		return
	}

	_typ, err := parseType(l)

	if err != nil {
		return
	}

	locEnd := _typ.Location().End

	defVal, _ := parseDefaultValue(l)

	directives, _ := parseDirectives(l)

	ret = &variableDefinition{}

	ret.Variable = *_var
	ret.Type = _typ
	ret.DefaultValue = defVal
	ret.Directives = directives
	ret.Loc = location{locStart, locEnd, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#Type
func parseType(l *lexer) _type {

	ret, err = parseNamedType(l)

	if err == nil {
		return
	}

	ret, err = parseListType(l)

	if err == nil {
		return
	}

	ret, err = parseNonNullType(l)

	if err != nil {
		err = errors.Wrap(err, "Expecting a type")
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#ListType
func parseListType(l *lexer) *listType {

	locStart := l.location().Start

	if !l.tokenEquals(tokBracketL.string()) {
		err = errors.New("Expecting '[' for list type")
		return
	}

	l.get()

	_typ, err := parseType(l)

	if err != nil {
		return
	}

	if !l.tokenEquals(tokBracketR.string()) {
		err = errors.New("Expecting ']' for list type")
		return
	}

	l.get()

	ret = &listType{}

	ret.OfType = _typ
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#NonNullType
func parseNonNullType(l *lexer) *nonNullType {

	locStart := l.location().Start

	var _typ _type

	_typ, err = parseNamedType(l)

	if err != nil {
		_typ, err = parseListType(l)

		if err != nil {
			return
		}
	}

	if !l.tokenEquals(tokBang.string()) {
		err = errors.New("Expecting '!' at the end of a non null type")
		return
	}

	l.get()

	ret = &nonNullType{}

	ret.OfType = _typ
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#Directives
func parseDirectives(l *lexer) *directives {

	ret = &directives{}

	for l.tokenEquals(tokAt.string()) {
		var dir *directive

		dir, err = parseDirective(l)

		if err != nil {
			ret = nil
			return
		}

		*ret = append(*ret, *dir)
	}

	if len(*ret) == 0 {
		err = errors.New("Expecting at least one directive")
		ret = nil
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#Directive
func parseDirective(l *lexer) *directive {

	locStart := l.location().Start

	if !l.tokenEquals(tokAt.string()) {
		err = errors.New("Expecting '@' for directive")
		return
	} else {
		l.get()

		var nam *name

		nam, err = parseName(l)

		if err != nil {
			return
		}

		args, _ := parseArguments(l)

		ret = &directive{}

		ret.Name = *nam
		ret.Arguments = args
		ret.Loc = location{locStart, l.prevLocation().End, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#SelectionSet
func parseSelectionSet(l *lexer) *selectionSet {

	if !l.tokenEquals(tokBraceL.string()) {
		err = errors.New("Expecting '{' for selection set")
		return
	} else {
		l.get()

		ret = &selectionSet{}

		for !l.tokenEquals(tokBraceR.string()) {
			var sel selection

			sel, err = parseSelection(l)

			if err != nil {
				ret = nil
				return
			}

			*ret = append(*ret, sel)
		}

		l.get()

		if len(*ret) == 0 {
			ret = nil
			err = errors.New("Expecting at least one selection")
			return
		}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#Selection
func parseSelection(l *lexer) selection {

	ret, err = parseField(l)

	if err == nil {
		return
	}

	ret, err = parseFragmentSpread(l)

	if err == nil {
		return
	}

	ret, err = parseInlineFragment(l)

	if err != nil {
		err = errors.Wrap(err, "Expecting a selection")
		ret = nil
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#Variable
func parseVariable(l *lexer) *variable {

	locStart := l.location().Start

	if !l.tokenEquals(tokDollar.string()) {
		err = errors.New("Expecting '$' for varible")
		return
	} else {
		var nam *name

		nam, err = parseName(l)

		if err != nil {
			return
		}

		ret = &variable{}

		ret.Name = *nam
		ret.Loc = location{locStart, nam.Location().End, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#DefaultValue
func parseDefaultValue(l *lexer) *defaultValue {

	locStart := l.location().Start

	if !l.tokenEquals(tokEquals.string()) {
		err = errors.New("Expecting '=' for default value")
		return
	} else {
		var val value

		val, err = parseValue(l)

		if err != nil {
			return
		}

		ret = &defaultValue{}

		ret.Value = val
		ret.Loc = location{locStart, l.prevLocation().End, l.source}

		return
	}
}

// ! need to check variable type in order to parse its value
// https://graphql.github.io/graphql-spec/draft/#Value
func parseValue(l *lexer) value {

	// need to parse dynamic variables
	//_var, _ := parseVariable(l)

	ret, err = parseIntValue(l)

	if err == nil {
		return
	}

	ret, err = parseFloatValue(l)

	if err == nil {
		return
	}

	ret, err = parseStringValue(l)

	if err == nil {
		return
	}

	ret, err = parseBooleanValue(l)

	if err == nil {
		return
	}

	ret, err = parseNullValue(l)

	if err == nil {
		return
	}

	ret, err = parseEnumValue(l)

	if err == nil {
		return
	}

	ret, err = parseListValue(l)

	if err == nil {
		return
	}

	ret, err = parseObjectValue(l)

	if err != nil {
		err = errors.Wrap(err, "Expecting a value")
		ret = nil
		return
	}

	return
}

// https://graphql.github.io/graphql-spec/draft/#Arguments
func parseArguments(l *lexer) *arguments {

	if !l.tokenEquals(tokParenL.string()) {
		err = errors.New("Expecting '(' for arguments")
		return
	} else {
		l.get()

		ret = &arguments{}

		for !l.tokenEquals(tokParenR.string()) {
			var arg *argument

			arg, err = parseArgument(l)

			if err != nil {
				ret = nil
				return
			}

			*ret = append(*ret, *arg)

		}

		l.get()

		if len(*ret) == 0 {
			err = errors.New("Expecting at least one argument")
			ret = nil
			return
		}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#Argument
func parseArgument(l *lexer) *argument {

	nam, err := parseName(l)

	if err != nil {
		return
	}

	if !l.tokenEquals(tokColon.string()) {
		err = errors.New("Expecting colon after argument name")
		return
	}

	l.get()

	val, err := parseValue(l)

	if err != nil {
		return
	}

	ret = &argument{}

	ret.Name = *nam
	ret.Value = val
	ret.Loc = location{nam.Location().Start, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#Field
func parseField(l *lexer) *field {

	locStart := l.location().Start

	alia, err := parseName(l)

	if err != nil {
		return
	}

	nam := &name{}

	if l.tokenEquals(tokColon.string()) {
		l.get()

		nam, err = parseName(l)

		if err != nil {
			return
		}
	} else {
		*nam = *alia

		alia = nil
	}

	args, _ := parseArguments(l)

	dirs, _ := parseDirectives(l)

	selSet, _ := parseSelectionSet(l)

	ret = &field{}

	ret.Alias = (*alias)(alia)
	ret.Name = *nam
	ret.Arguments = args
	ret.Directives = dirs
	ret.SelectionSet = selSet
	ret.Loc = location{locStart, l.prevLocation().End, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#FragmentSpread
func parseFragmentSpread(l *lexer) *fragmentSpread {

	locStart := l.location().Start

	if !l.tokenEquals(tokSpread.string()) {
		err = errors.New("Expecting '...' operator for a fragment spread")
		return
	} else {
		l.get()

		var fname *fragmentName

		fname, err = parseFragmentName(l)

		if err != nil {
			return
		}

		directives, _ := parseDirectives(l)

		ret = &fragmentSpread{}

		ret.FragmentName = *fname
		ret.Directives = directives
		ret.Loc = location{locStart, l.prevLocation().End, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#InlineFragment
func parseInlineFragment(l *lexer) *inlineFragment {

	locStart := l.location().Start

	if !l.tokenEquals(tokSpread.string()) {
		err = errors.New("Expecting '...' for an inline fragment")
		return
	} else {
		l.get()

		typeCon, _ := parseTypeCondition(l)

		directives, _ := parseDirectives(l)

		var selSet *selectionSet

		selSet, err = parseSelectionSet(l)

		if err != nil {
			return
		}

		ret = &inlineFragment{}

		ret.TypeCondition = typeCon
		ret.Directives = directives
		ret.SelectionSet = *selSet
		ret.Loc = location{locStart, l.prevLocation().End, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#FragmentName
func parseFragmentName(l *lexer) *fragmentName {

	nam, err := parseName(l)

	if err != nil {
		return
	}

	if nam.Value == kwOn {
		err = errors.New("Fragment name cannot be 'on'")
		return
	}

	ret = &fragmentName{}
	*ret = fragmentName(*nam)
	ret.Loc = *nam.Location()

	return
}

// https://graphql.github.io/graphql-spec/draft/#TypeCondition
func parseTypeCondition(l *lexer) *typeCondition {

	locStart := l.location().Start

	if !l.tokenEquals(kwOn) {
		err = errors.New("Expecting 'on' keyword for a type condition")
		return
	} else {
		var namedTyp *namedType

		namedTyp, err = parseNamedType(l)

		if err != nil {
			return
		}

		ret = &typeCondition{}

		ret.NamedType = *namedTyp
		ret.Loc = location{locStart, namedTyp.Location().End, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#NamedType
func parseNamedType(l *lexer) *namedType {

	nam, err := parseName(l)

	if err != nil {
		return
	}

	ret = new(namedType)

	*ret = namedType(*nam)

	return
}

// ! Check numeric values
// https://graphql.github.io/graphql-spec/draft/#IntValue
func parseIntValue(l *lexer) *intValue {

	tok := l.current()

	intVal, err := strconv.ParseInt(tok.value, 10, 64)

	if err != nil {
		return
	}

	l.get()

	ret = &intValue{}

	ret.Value = intVal
	ret.Loc = location{tok.start, tok.end, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#FloatValue
func parseFloatValue(l *lexer) *floatValue {

	tok := l.current()

	floatVal, err := strconv.ParseFloat(tok.value, 64)

	if err != nil {
		return
	}

	l.get()

	ret = &floatValue{}

	ret.Value = floatVal
	ret.Loc = location{tok.start, tok.end, l.source}

	return
}

// ! Have a discussion about this function
// https://graphql.github.io/graphql-spec/draft/#StringValue
func parseStringValue(l *lexer) *stringValue {

	tok := l.current()

	ret = &stringValue{}

	ret.Value = tok.value
	ret.Loc = location{tok.start, tok.end, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#BooleanValue
func parseBooleanValue(l *lexer) *booleanValue {

	tok := l.current()

	boolVal, err := strconv.ParseBool(tok.value)

	if err != nil {
		return
	}

	l.get()

	ret = &booleanValue{}

	ret.Value = boolVal
	ret.Loc = location{tok.start, tok.end, l.source}

	return
}

// https://graphql.github.io/graphql-spec/draft/#NullValue
func parseNullValue(l *lexer) *nullValue {

	tok := l.current()

	if tok.value != kwNull {
		err = errors.New("Expecting 'null' keyword")
		return
	} else {
		l.get()

		ret = &nullValue{}
		ret.Loc = location{tok.start, tok.end, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#EnumValue
func parseEnumValue(l *lexer) *enumValue {

	nam, err := parseName(l)

	if err != nil {
		return
	}

	switch nam.Value {
	case kwTrue, kwFalse, kwNull:
		err = errors.New("Enum value cannot be 'true', 'false' or 'null'")
		return
	default:
		ret = &enumValue{}

		ret.Name = *nam
		ret.Loc = location{nam.Location().Start, nam.Location().End, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#ListValue
func parseListValue(l *lexer) *listValue {

	locStart := l.location().Start

	if !l.tokenEquals(tokBracketL.string()) {
		return nil, errors.New("Expecting '[' for a list value")
	} else {
		l.get()

		ret = &listValue{}

		for !l.tokenEquals(tokBracketR.string()) {
			var val value

			val, err = parseValue(l)

			if err != nil {
				return
			}

			ret.Values = append(ret.Values, val)
		}

		ret.Loc = location{locStart, l.prevLocation().End, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectValue
func parseObjectValue(l *lexer) *objectValue {

	locStart := l.location().Start

	if !l.tokenEquals(tokBraceL.string()) {
		err = errors.New("Expecting '{' for an object value")
		return
	} else {
		l.get()

		ret = &objectValue{}

		for !l.tokenEquals(tokBraceR.string()) {
			var objField *objectField

			objField, err = parseObjectField(l)

			if err != nil {
				return
			}

			ret.Values = append(ret.Values, *objField)
		}

		ret.Loc = location{locStart, l.prevLocation().End, l.source}

		return
	}
}

// https://graphql.github.io/graphql-spec/draft/#ObjectField
func parseObjectField(l *lexer) *objectField {

	nam, err := parseName(l)

	if err != nil {
		return
	}

	if !l.tokenEquals(tokColon.string()) {
		err = errors.New("Expecting color after object field name")
		return
	}

	l.get()

	val, err := parseValue(l)

	if err != nil {
		return
	}

	ret = &objectField{}

	ret.Name = *nam
	ret.Value = val
	ret.Loc = location{nam.Location().Start, l.prevLocation().End, l.source}

	return
}
