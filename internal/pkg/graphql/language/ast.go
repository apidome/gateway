package language

const (
	// Executable Directive Locations
	edlQuery              executableDirectiveLocation = "QUERY"
	edlMutation                                       = "MUTATION"
	edlSubscription                                   = "SUBSCRIPTION"
	edlField                                          = "FIELD"
	edlFragmentDefinition                             = "FRAGMENT_DEFINITION"
	edlFragmentSpread                                 = "FRAGMENT_SPREAD"
	edlInlineFragment                                 = "INLINE_FRAGMENT"
	edlVariableDefinition                             = "VARIABLE_DEFINITION"

	// Type System Directive Locations
	tsdlSchema               typeSystemDirectiveLocation = "SCHEMA"
	tsdlScalar                                           = "SCALAR"
	tsdlObject                                           = "OBJECT"
	tsdlFieldDefinition                                  = "FIELD_DEFINITION"
	tsdlArgumentDefinition                               = "ARGUMENT_DEFINITION"
	tsdlInterface                                        = "INTERFACE"
	tsdlUnion                                            = "UNION"
	tsdlEnum                                             = "ENUM"
	tsdlEnumValue                                        = "ENUM_VALUE"
	tsdlInputObject                                      = "INPUT_OBJECT"
	tsdlInputFieldDefinition                             = "INPUT_FIELD_DEFINITION"
)

const (
	operationQuery        operationType = "query"
	operationMutation     operationType = "mutation"
	operationSubscription operationType = "subscription"
)

var executableDirectiveLocations []executableDirectiveLocation = []executableDirectiveLocation{
	edlQuery,
	edlMutation,
	edlSubscription,
	edlField,
	edlFragmentDefinition,
	edlFragmentSpread,
	edlInlineFragment,
	edlVariableDefinition,
}

var typeSystemDirectiveLocations []typeSystemDirectiveLocation = []typeSystemDirectiveLocation{
	tsdlSchema,
	tsdlScalar,
	tsdlObject,
	tsdlFieldDefinition,
	tsdlArgumentDefinition,
	tsdlInterface,
	tsdlUnion,
	tsdlEnum,
	tsdlEnumValue,
	tsdlInputObject,
	tsdlInputFieldDefinition,
}

type directiveLocation string

type executableDirectiveLocation directiveLocation

type typeSystemDirectiveLocation directiveLocation

type location struct {
	Start  int
	End    int
	Source string
}

type document struct {
	Definitions []definition
}

type arguments []argument

type argument struct {
	Name  name
	Value value
	locator
}

type directiveLocations []directiveLocation

type directives []directive

type directive struct {
	Name      name
	Arguments *arguments
	locator
}

type operationType string

type locatorInterface interface {
	Location() *location
}

type locator struct {
	Loc location
}

func (l *locator) Location() *location {
	return &l.Loc
}

type definition interface {
	locatorInterface
}

type definitions []definition

var _ definition = (*rootOperationTypeDefinition)(nil)
var _ definition = (*enumValueDefinition)(nil)
var _ definition = (*inputValueDefinition)(nil)

/*
	Executable Definitions
*/
type executableDefinition interface {
	definition
}

var _ executableDefinition = (*operationDefinition)(nil)
var _ executableDefinition = (*fragmentDefinition)(nil)

type operationDefinition struct {
	OperationType       operationType
	Name                *name
	VariableDefinitions *variableDefinitions
	SelectionSet        selectionSet
	Directives          *directives
	locator
}

type fragmentDefinition struct {
	FragmentName  fragmentName
	TypeCondition typeCondition
	SelectionSet  selectionSet
	Directives    *directives
	locator
}

/*
	Type System Definitions
*/

type typeSystemDefinition interface {
	definition
}

var _ typeSystemDefinition = (typeDefinition)(nil)
var _ typeSystemDefinition = (*schemaDefinition)(nil)
var _ typeSystemDefinition = (*directiveDefinition)(nil)

type schemaDefinition struct {
	Directives                   *directives
	RootOperationTypeDefinitions []rootOperationTypeDefinition
	locator
}

type rootOperationTypeDefinition struct {
	OperationType operationType
	NamedType     namedType
	locator
}

type rootOperationTypeDefinitions []rootOperationTypeDefinition

type directiveDefinition struct {
	Description         *description
	Name                name
	ArgumentsDefinition *argumentsDefinition
	DirectiveLocations  directiveLocations
	locator
}

type inputValueDefinition struct {
	Description  *description
	Name         name
	Type         _type
	DefaultValue *defaultValue
	Directives   *directives
	locator
}

type argumentsDefinition []inputValueDefinition

type typeDefinition interface {
	typeSystemDefinition
}

var _ typeDefinition = (*scalarTypeDefinition)(nil)
var _ typeDefinition = (*objectTypeDefinition)(nil)
var _ typeDefinition = (*interfaceTypeDefinition)(nil)
var _ typeDefinition = (*unionTypeDefinition)(nil)
var _ typeDefinition = (*enumTypeDefinition)(nil)
var _ typeDefinition = (*inputObjectTypeDefinition)(nil)

type scalarTypeDefinition struct {
	Description *description
	Name        name
	Directives  *directives
	locator
}

type objectTypeDefinition struct {
	Description          *description
	Name                 name
	ImplementsInterfaces *implementsInterfaces
	Directives           *directives
	FieldsDefinition     *fieldsDefinition
	locator
}

type interfaceTypeDefinition struct {
	Description      *description
	Name             name
	Directives       *directives
	FieldsDefinition *fieldsDefinition
	locator
}

type unionTypeDefinition struct {
	Description      *description
	Name             name
	Directives       *directives
	UnionMemberTypes *unionMemberTypes
	locator
}

type enumTypeDefinition struct {
	Description          *description
	Name                 name
	Directives           *directives
	EnumValuesDefinition *enumValuesDefinition
	locator
}

type enumValueDefinition struct {
	Description *description
	EnumValue   enumValue
	Directives  *directives
	locator
}

type enumValuesDefinition []enumValueDefinition

type inputObjectTypeDefinition struct {
	Description           *description
	Name                  name
	Directives            *directives
	InputFieldsDefinition *inputFieldsDefinition
	locator
}

type inputFieldsDefinition []inputValueDefinition

type variableDefinition struct {
	Variable     variable
	Type         _type
	DefaultValue *defaultValue
	Directives   *directives
	locator
}

type variableDefinitions []variableDefinition

type fieldDefinition struct {
	Description         *description
	Name                name
	ArgumentsDefinition *argumentsDefinition
	Type                _type
	Directives          *directives
	locator
}

type fieldsDefinition []fieldDefinition

/*
	Type System Extensions
*/

type typeSystemExtension interface {
	definition
}

var _ typeSystemExtension = (*schemaExtension)(nil)

type schemaExtension struct {
	Directives                   *directives
	RootOperationTypeDefinitions *rootOperationTypeDefinitions
	locator
}

type typeExtension interface {
	typeSystemExtension
}

var _ typeExtension = (*scalarTypeExtension)(nil)
var _ typeExtension = (*objectTypeExtension)(nil)
var _ typeExtension = (*interfaceTypeExtension)(nil)
var _ typeExtension = (*unionTypeExtension)(nil)
var _ typeExtension = (*enumTypeExtension)(nil)
var _ typeExtension = (*inputObjectTypeExtension)(nil)

type scalarTypeExtension struct {
	Name       name
	Directives directives
	locator
}

type objectTypeExtension struct {
	Name                 name
	ImplementsInterfaces *implementsInterfaces
	Directives           *directives
	FieldsDefinition     *fieldsDefinition
	locator
}

type interfaceTypeExtension struct {
	Name             name
	Directives       *directives
	FieldsDefinition *fieldsDefinition
	locator
}

type unionTypeExtension struct {
	Name             name
	Directives       *directives
	UnionMemberTypes *unionMemberTypes
	locator
}

type enumTypeExtension struct {
	Name                 name
	Directives           *directives
	EnumValuesDefinition *enumValuesDefinition
	locator
}

type inputObjectTypeExtension struct {
	Name                  name
	Directives            *directives
	InputFieldsDefinition *inputFieldsDefinition
	locator
}

type selection interface {
	locatorInterface
	Selections() *selectionSet
}

var _ selection = (*field)(nil)
var _ selection = (*fragmentSpread)(nil)
var _ selection = (*inlineFragment)(nil)

type selectionSet []selection

type field struct {
	Alias        *alias
	Name         name
	Arguments    *arguments
	Directives   *directives
	SelectionSet *selectionSet
	locator
}

func (f *field) Selections() *selectionSet {
	return f.SelectionSet
}

type fragmentSpread struct {
	FragmentName fragmentName
	Directives   *directives
	locator
}

// ! Find the fragment definition by its name and return its selection set
func (f *fragmentSpread) Selections() *selectionSet {
	return &selectionSet{}
}

type inlineFragment struct {
	TypeCondition *typeCondition
	Directives    *directives
	SelectionSet  selectionSet
	locator
}

func (i *inlineFragment) Selections() *selectionSet {
	return &i.SelectionSet
}

type _type interface {
	locatorInterface
}

var _ _type = (*namedType)(nil)
var _ _type = (*listType)(nil)
var _ _type = (*nonNullType)(nil)

type namedType name

type listType struct {
	Kind   string
	OfType _type
	locator
}

type nonNullType struct {
	Kind   string
	OfType _type
	locator
}

type typeCondition struct {
	NamedType namedType
	locator
}

type unionMemberTypes []namedType

type implementsInterfaces []namedType

type name struct {
	Value string
	locator
}

type alias name
type fragmentName name

type value interface {
	GetValue() interface{}
	locatorInterface
}

var _ value = (*defaultValue)(nil)
var _ value = (*objectField)(nil)
var _ value = (*objectValue)(nil)
var _ value = (*listValue)(nil)
var _ value = (*intValue)(nil)
var _ value = (*floatValue)(nil)
var _ value = (*stringValue)(nil)
var _ value = (*booleanValue)(nil)
var _ value = (*enumValue)(nil)
var _ value = (*variable)(nil)

type defaultValue struct {
	Value value
	locator
}

func (dv defaultValue) GetValue() interface{} {
	return dv.Value.GetValue()
}

type objectField struct {
	Name  name
	Value value
	locator
}

func (of objectField) GetValue() interface{} {
	return of.Value.GetValue()
}

type objectValue struct {
	Values []objectField
	locator
}

func (ov objectValue) GetValue() interface{} {
	return ov
}

type listValue struct {
	Values []value
	locator
}

func (lv listValue) GetValue() interface{} {
	return lv
}

type intValue struct {
	Value int64
	locator
}

func (iv intValue) GetValue() interface{} {
	return iv
}

type floatValue struct {
	Value float64
	locator
}

func (fv floatValue) GetValue() interface{} {
	return fv
}

type stringValue struct {
	Value string
	locator
}

func (sv stringValue) GetValue() interface{} {
	return sv
}

type booleanValue struct {
	Value bool
	locator
}

func (bv booleanValue) GetValue() interface{} {
	return bv
}

type enumValue struct {
	Name name
	locator
}

func (ev enumValue) GetValue() interface{} {
	return ev
}

type nullValue struct {
	locator
}

func (nv nullValue) GetValue() interface{} {
	return nv
}

type variable struct {
	Name name
	locator
}

func (v variable) GetValue() interface{} {
	return v.Name
}

type description stringValue
