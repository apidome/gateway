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

type description stringValue

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
	definition() definition
}

type definitions []definition

var _ definition = (executableDefinition)(nil)
var _ definition = (typeSystemDefinition)(nil)
var _ definition = (typeSystemExtension)(nil)

/*
	Executable Definitions
*/
type executableDefinition interface {
	definition
	executableDefinition() executableDefinition
	GetName() *name
	GetDirectives() *directives
	GetSelectionSet() selectionSet
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

func (o *operationDefinition) definition() definition {
	return o
}

func (o *operationDefinition) executableDefinition() executableDefinition {
	return o
}

func (o *operationDefinition) GetName() *name {
	return o.Name
}

func (o *operationDefinition) GetDirectives() *directives {
	return o.Directives
}

func (o *operationDefinition) GetSelectionSet() selectionSet {
	return o.SelectionSet
}

type fragmentDefinition struct {
	FragmentName  name
	TypeCondition typeCondition
	SelectionSet  selectionSet
	Directives    *directives
	locator
}

func (f *fragmentDefinition) definition() definition {
	return f
}

func (f *fragmentDefinition) executableDefinition() executableDefinition {
	return f
}

func (f *fragmentDefinition) GetName() *name {
	return &f.FragmentName
}

func (f *fragmentDefinition) GetDirectives() *directives {
	return f.Directives
}

func (f *fragmentDefinition) GetSelectionSet() selectionSet {
	return f.SelectionSet
}

/*
	Type System Definitions
*/

type typeSystemDefinition interface {
	definition
	typeSystemDefinition() typeSystemDefinition
}

var _ typeSystemDefinition = (typeDefinition)(nil)
var _ typeSystemDefinition = (*schemaDefinition)(nil)
var _ typeSystemDefinition = (*directiveDefinition)(nil)

type schemaDefinition struct {
	Directives                   *directives
	RootOperationTypeDefinitions []rootOperationTypeDefinition
	locator
}

func (s *schemaDefinition) definition() definition {
	return s
}

func (s *schemaDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
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

func (d *directiveDefinition) definition() definition {
	return d
}

func (d *directiveDefinition) typeSystemDefinition() typeSystemDefinition {
	return d
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
	typeDefinition() typeDefinition
	GetDescription() *description
	GetName() name
	GetDirectives() *directives
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

func (s *scalarTypeDefinition) definition() definition {
	return s
}

func (s *scalarTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *scalarTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (s *scalarTypeDefinition) GetDescription() *description {
	return s.Description
}

func (s *scalarTypeDefinition) GetName() name {
	return s.Name
}

func (s *scalarTypeDefinition) GetDirectives() *directives {
	return s.Directives
}

type objectTypeDefinition struct {
	Description          *description
	Name                 name
	ImplementsInterfaces *implementsInterfaces
	Directives           *directives
	FieldsDefinition     *fieldsDefinition
	locator
}

func (s *objectTypeDefinition) definition() definition {
	return s
}

func (s *objectTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *objectTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (o *objectTypeDefinition) GetDescription() *description {
	return o.Description
}

func (o *objectTypeDefinition) GetName() name {
	return o.Name
}

func (o *objectTypeDefinition) GetDirectives() *directives {
	return o.Directives
}

type interfaceTypeDefinition struct {
	Description      *description
	Name             name
	Directives       *directives
	FieldsDefinition *fieldsDefinition
	locator
}

func (s *interfaceTypeDefinition) definition() definition {
	return s
}

func (s *interfaceTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *interfaceTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (i *interfaceTypeDefinition) GetDescription() *description {
	return i.Description
}

func (i *interfaceTypeDefinition) GetName() name {
	return i.Name
}

func (i *interfaceTypeDefinition) GetDirectives() *directives {
	return i.Directives
}

type unionTypeDefinition struct {
	Description      *description
	Name             name
	Directives       *directives
	UnionMemberTypes *unionMemberTypes
	locator
}

func (s *unionTypeDefinition) definition() definition {
	return s
}

func (s *unionTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *unionTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (u *unionTypeDefinition) GetDescription() *description {
	return u.Description
}

func (u *unionTypeDefinition) GetName() name {
	return u.Name
}

func (u *unionTypeDefinition) GetDirectives() *directives {
	return u.Directives
}

type enumTypeDefinition struct {
	Description          *description
	Name                 name
	Directives           *directives
	EnumValuesDefinition *enumValuesDefinition
	locator
}

func (s *enumTypeDefinition) definition() definition {
	return s
}

func (s *enumTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *enumTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (e *enumTypeDefinition) GetDescription() *description {
	return e.Description
}

func (e *enumTypeDefinition) GetName() name {
	return e.Name
}

func (e *enumTypeDefinition) GetDirectives() *directives {
	return e.Directives
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

func (s *inputObjectTypeDefinition) definition() definition {
	return s
}

func (s *inputObjectTypeDefinition) typeSystemDefinition() typeSystemDefinition {
	return s
}

func (s *inputObjectTypeDefinition) typeDefinition() typeDefinition {
	return s
}

func (i *inputObjectTypeDefinition) GetDescription() *description {
	return i.Description
}

func (i *inputObjectTypeDefinition) GetName() name {
	return i.Name
}

func (i *inputObjectTypeDefinition) GetDirectives() *directives {
	return i.Directives
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
	typeSystemExtension() typeSystemExtension
	GetDirectives() *directives
}

var _ typeSystemExtension = (*schemaExtension)(nil)
var _ typeSystemExtension = (typeExtension)(nil)

type schemaExtension struct {
	Directives                   *directives
	RootOperationTypeDefinitions *rootOperationTypeDefinitions
	locator
}

func (s *schemaExtension) definition() definition {
	return s
}

func (s *schemaExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *schemaExtension) GetDirectives() *directives {
	return s.Directives
}

type typeExtension interface {
	typeSystemExtension
	typeExtension() typeExtension
	GetName() name
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

func (s *scalarTypeExtension) definition() definition {
	return s
}

func (s *scalarTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *scalarTypeExtension) typeExtension() typeExtension {
	return s
}

func (s *scalarTypeExtension) GetName() name {
	return s.Name
}

func (s *scalarTypeExtension) GetDirectives() *directives {
	return &s.Directives
}

type objectTypeExtension struct {
	Name                 name
	ImplementsInterfaces *implementsInterfaces
	Directives           *directives
	FieldsDefinition     *fieldsDefinition
	locator
}

func (s *objectTypeExtension) definition() definition {
	return s
}

func (s *objectTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *objectTypeExtension) typeExtension() typeExtension {
	return s
}

func (o *objectTypeExtension) GetName() name {
	return o.Name
}

func (o *objectTypeExtension) GetDirectives() *directives {
	return o.Directives
}

type interfaceTypeExtension struct {
	Name             name
	Directives       *directives
	FieldsDefinition *fieldsDefinition
	locator
}

func (s *interfaceTypeExtension) definition() definition {
	return s
}

func (s *interfaceTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *interfaceTypeExtension) typeExtension() typeExtension {
	return s
}

func (i *interfaceTypeExtension) GetName() name {
	return i.Name
}

func (i *interfaceTypeExtension) GetDirectives() *directives {
	return i.Directives
}

type unionTypeExtension struct {
	Name             name
	Directives       *directives
	UnionMemberTypes *unionMemberTypes
	locator
}

func (s *unionTypeExtension) definition() definition {
	return s
}

func (s *unionTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *unionTypeExtension) typeExtension() typeExtension {
	return s
}

func (u *unionTypeExtension) GetName() name {
	return u.Name
}

func (u *unionTypeExtension) GetDirectives() *directives {
	return u.Directives
}

type enumTypeExtension struct {
	Name                 name
	Directives           *directives
	EnumValuesDefinition *enumValuesDefinition
	locator
}

func (s *enumTypeExtension) definition() definition {
	return s
}

func (s *enumTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *enumTypeExtension) typeExtension() typeExtension {
	return s
}

func (e *enumTypeExtension) GetName() name {
	return e.Name
}

func (e *enumTypeExtension) GetDirectives() *directives {
	return e.Directives
}

type inputObjectTypeExtension struct {
	Name                  name
	Directives            *directives
	InputFieldsDefinition *inputFieldsDefinition
	locator
}

func (s *inputObjectTypeExtension) definition() definition {
	return s
}

func (s *inputObjectTypeExtension) typeSystemExtension() typeSystemExtension {
	return s
}

func (s *inputObjectTypeExtension) typeExtension() typeExtension {
	return s
}

func (i *inputObjectTypeExtension) GetName() name {
	return i.Name
}

func (i *inputObjectTypeExtension) GetDirectives() *directives {
	return i.Directives
}

type selection interface {
	locatorInterface
	selection() selection
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

func (f *field) selection() selection {
	return f
}

func (f *field) Selections() *selectionSet {
	return f.SelectionSet
}

type fragmentSpread struct {
	FragmentName name
	Directives   *directives
	locator
}

func (f *fragmentSpread) selection() selection {
	return f
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

func (i *inlineFragment) selection() selection {
	return i
}

func (i *inlineFragment) Selections() *selectionSet {
	return &i.SelectionSet
}

type _type interface {
	_type() _type
	locatorInterface
}

var _ _type = (*namedType)(nil)
var _ _type = (*listType)(nil)
var _ _type = (*nonNullType)(nil)

type namedType name

func (n *namedType) _type() _type {
	return n
}

type listType struct {
	Kind   string
	OfType _type
	locator
}

func (n *listType) _type() _type {
	return n
}

type nonNullType struct {
	Kind   string
	OfType _type
	locator
}

func (n *nonNullType) _type() _type {
	return n
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

type value interface {
	locatorInterface
	value() value
	GetValue() interface{}
}

var _ value = (*objectField)(nil)
var _ value = (*objectValue)(nil)
var _ value = (*listValue)(nil)
var _ value = (*intValue)(nil)
var _ value = (*floatValue)(nil)
var _ value = (*stringValue)(nil)
var _ value = (*booleanValue)(nil)
var _ value = (*enumValue)(nil)
var _ value = (*variable)(nil)

type objectField struct {
	Name  name
	Value value
	locator
}

func (d *objectField) value() value {
	return d
}

func (of objectField) GetValue() interface{} {
	return of.Value.GetValue()
}

type objectValue struct {
	Values []objectField
	locator
}

func (d *objectValue) value() value {
	return d
}

func (ov objectValue) GetValue() interface{} {
	return ov
}

type listValue struct {
	Values []value
	locator
}

func (d *listValue) value() value {
	return d
}

func (lv listValue) GetValue() interface{} {
	return lv
}

type intValue struct {
	Value int64
	locator
}

func (d *intValue) value() value {
	return d
}

func (iv intValue) GetValue() interface{} {
	return iv
}

type floatValue struct {
	Value float64
	locator
}

func (d *floatValue) value() value {
	return d
}

func (fv floatValue) GetValue() interface{} {
	return fv
}

type stringValue struct {
	Value string
	locator
}

func (d *stringValue) value() value {
	return d
}

func (sv stringValue) GetValue() interface{} {
	return sv
}

type booleanValue struct {
	Value bool
	locator
}

func (d *booleanValue) value() value {
	return d
}

func (bv booleanValue) GetValue() interface{} {
	return bv
}

type enumValue struct {
	Name name
	locator
}

func (d *enumValue) value() value {
	return d
}

func (ev enumValue) GetValue() interface{} {
	return ev
}

type nullValue struct {
	locator
}

func (d *nullValue) value() value {
	return d
}

func (nv nullValue) GetValue() interface{} {
	return nv
}

type variable struct {
	Name name
	locator
}

func (d *variable) value() value {
	return d
}

func (v variable) GetValue() interface{} {
	return v.Name
}

type defaultValue struct {
	Value value
	locator
}

func (dv defaultValue) GetValue() interface{} {
	return dv.Value.GetValue()
}
