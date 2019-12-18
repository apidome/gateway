package language

const (
	// Executable Directive Locations
	QUERY               ExecutableDirectiveLocation = "QUERY"
	MUTATION                                        = "MUTATION"
	SUBSCRIPTION                                    = "SUBSCRIPTION"
	FIELD                                           = "FIELD"
	FRAGMENT_DEFINITION                             = "FRAGMENT_DEFINITION"
	FRAGMENT_SPREAD                                 = "FRAGMENT_SPREAD"
	INLINE_FRAGMENT                                 = "INLINE_FRAGMENT"
	VARIABLE_DEFINITION                             = "VARIABLE_DEFINITION"

	// Type System Directive Locations
	SCHEMA                 TypeSystemDirectiveLocation = "SCHEMA"
	SCALAR                                             = "SCALAR"
	OBJECT                                             = "OBJECT"
	FIELD_DEFINITION                                   = "FIELD_DEFINITION"
	ARGUMENT_DEFINITION                                = "ARGUMENT_DEFINITION"
	INTERFACE                                          = "INTERFACE"
	UNION                                              = "UNION"
	ENUM                                               = "ENUM"
	ENUM_VALUE                                         = "ENUM_VALUE"
	INPUT_OBJECT                                       = "INPUT_OBJECT"
	INPUT_FIELD_DEFINITION                             = "INPUT_FIELD_DEFINITION"
)

const (
	OPERATION_QUERY        OperationType = "query"
	OPERATION_MUTATION     OperationType = "mutation"
	OPERATION_SUBSCRIPTION OperationType = "subscription"
)

var executableDirectiveLocations []ExecutableDirectiveLocation = []ExecutableDirectiveLocation{
	QUERY,
	MUTATION,
	SUBSCRIPTION,
	FIELD,
	FRAGMENT_DEFINITION,
	FRAGMENT_SPREAD,
	INLINE_FRAGMENT,
	VARIABLE_DEFINITION,
}

var typeSystemDirectiveLocations []TypeSystemDirectiveLocation = []TypeSystemDirectiveLocation{
	SCHEMA,
	SCALAR,
	OBJECT,
	FIELD_DEFINITION,
	ARGUMENT_DEFINITION,
	INTERFACE,
	UNION,
	ENUM,
	ENUM_VALUE,
	INPUT_OBJECT,
	INPUT_FIELD_DEFINITION,
}

type DirectiveLocation string

type ExecutableDirectiveLocation DirectiveLocation

type TypeSystemDirectiveLocation DirectiveLocation

type Location struct {
	Start  int
	End    int
	Source string
}

type Document struct {
	Definitions []Definition
}

type Arguments []Argument

type Argument struct {
	Name  Name
	Value Value
	Locator
}

type DirectiveLocations []DirectiveLocation

type Directives []Directive

type Directive struct {
	Name      Name
	Arguments *Arguments
	Locator
}

type OperationType string

type LocatorInterface interface {
	Location() *Location
}

type Locator struct {
	Loc Location
}

func (l *Locator) Location() *Location {
	return &l.Loc
}

type Definition interface {
	LocatorInterface
}

type Definitions []Definition

var _ Definition = (*RootOperationTypeDefinition)(nil)
var _ Definition = (*EnumValueDefinition)(nil)
var _ Definition = (*InputValueDefinition)(nil)

/*
	Executable Definitions
*/
type ExecutableDefinition interface {
	Definition
}

var _ ExecutableDefinition = (*OperationDefinition)(nil)
var _ ExecutableDefinition = (*FragmentDefinition)(nil)

type OperationDefinition struct {
	OperationType       OperationType
	Name                *Name
	VariableDefinitions *VariableDefinitions
	SelectionSet        SelectionSet
	Directives          *Directives
	Locator
}

type FragmentDefinition struct {
	FragmentName  FragmentName
	TypeCondition TypeCondition
	SelectionSet  SelectionSet
	Directives    *Directives
	Locator
}

/*
	Type System Definitions
*/

type TypeSystemDefinition interface {
	Definition
}

var _ TypeSystemDefinition = (TypeDefinition)(nil)
var _ TypeSystemDefinition = (*SchemaDefinition)(nil)
var _ TypeSystemDefinition = (*DirectiveDefinition)(nil)

type SchemaDefinition struct {
	Directives                   *Directives
	RootOperationTypeDefinitions []RootOperationTypeDefinition
	Locator
}

type RootOperationTypeDefinition struct {
	OperationType OperationType
	NamedType     NamedType
	Locator
}

type RootOperationTypeDefinitions []RootOperationTypeDefinition

type DirectiveDefinition struct {
	Description         *Description
	Name                Name
	ArgumentsDefinition *ArgumentsDefinition
	DirectiveLocations  DirectiveLocations
	Locator
}

type InputValueDefinition struct {
	Description  *Description
	Name         Name
	Type         Type
	DefaultValue *DefaultValue
	Directives   *Directives
	Locator
}

type ArgumentsDefinition []InputValueDefinition

type TypeDefinition interface {
	TypeSystemDefinition
}

var _ TypeDefinition = (*ScalarTypeDefinition)(nil)
var _ TypeDefinition = (*ObjectTypeDefinition)(nil)
var _ TypeDefinition = (*InterfaceTypeDefinition)(nil)
var _ TypeDefinition = (*UnionTypeDefinition)(nil)
var _ TypeDefinition = (*EnumTypeDefinition)(nil)
var _ TypeDefinition = (*InputObjectTypeDefinition)(nil)

type ScalarTypeDefinition struct {
	Description *Description
	Name        Name
	Directives  *Directives
	Locator
}

type ObjectTypeDefinition struct {
	Description          *Description
	Name                 Name
	ImplementsInterfaces *ImplementsInterfaces
	Directives           *Directives
	FieldsDefinition     *FieldsDefinition
	Locator
}

type InterfaceTypeDefinition struct {
	Description      *Description
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	Locator
}

type UnionTypeDefinition struct {
	Description      *Description
	Name             Name
	Directives       *Directives
	UnionMemberTypes *UnionMemberTypes
	Locator
}

type EnumTypeDefinition struct {
	Description          *Description
	Name                 Name
	Directives           *Directives
	EnumValuesDefinition *EnumValuesDefinition
	Locator
}

type EnumValueDefinition struct {
	Description *Description
	EnumValue   EnumValue
	Directives  *Directives
	Locator
}

type EnumValuesDefinition []EnumValueDefinition

type InputObjectTypeDefinition struct {
	Description           *Description
	Name                  Name
	Directives            *Directives
	InputFieldsDefinition *InputFieldsDefinition
	Locator
}

type InputFieldsDefinition []InputValueDefinition

type VariableDefinition struct {
	Variable     Variable
	Type         Type
	DefaultValue *DefaultValue
	Directives   *Directives
	Locator
}

type VariableDefinitions []VariableDefinition

type FieldDefinition struct {
	Description         *Description
	Name                Name
	ArgumentsDefinition *ArgumentsDefinition
	Type                Type
	Directives          *Directives
	Locator
}

type FieldsDefinition []FieldDefinition

/*
	Type System Extensions
*/

type TypeSystemExtension interface {
	Definition
}

var _ TypeSystemExtension = (*SchemaExtension)(nil)

type SchemaExtension struct {
	Directives                   *Directives
	RootOperationTypeDefinitions *RootOperationTypeDefinitions
	Locator
}

type TypeExtension interface {
	TypeSystemExtension
}

var _ TypeExtension = (*ScalarTypeExtension)(nil)
var _ TypeExtension = (*ObjectTypeExtension)(nil)
var _ TypeExtension = (*InterfaceTypeExtension)(nil)
var _ TypeExtension = (*UnionTypeExtension)(nil)
var _ TypeExtension = (*EnumTypeExtension)(nil)
var _ TypeExtension = (*InputObjectTypeExtension)(nil)

type ScalarTypeExtension struct {
	Name       Name
	Directives Directives
	Locator
}

type ObjectTypeExtension struct {
	Name                 Name
	ImplementsInterfaces *ImplementsInterfaces
	Directives           *Directives
	FieldsDefinition     *FieldsDefinition
	Locator
}

type InterfaceTypeExtension struct {
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	Locator
}

type UnionTypeExtension struct {
	Name             Name
	Directives       *Directives
	UnionMemberTypes *UnionMemberTypes
	Locator
}

type EnumTypeExtension struct {
	Name                 Name
	Directives           *Directives
	EnumValuesDefinition *EnumValuesDefinition
	Locator
}

type InputObjectTypeExtension struct {
	Name                  Name
	Directives            *Directives
	InputFieldsDefinition *InputFieldsDefinition
	Locator
}

type Selection interface {
	LocatorInterface
	Selections() *SelectionSet
}

var _ Selection = (*Field)(nil)
var _ Selection = (*FragmentSpread)(nil)
var _ Selection = (*InlineFragment)(nil)

type SelectionSet []Selection

type Field struct {
	Alias        *Alias
	Name         Name
	Arguments    *Arguments
	Directives   *Directives
	SelectionSet *SelectionSet
	Locator
}

func (f *Field) Selections() *SelectionSet {
	return f.SelectionSet
}

type FragmentSpread struct {
	FragmentName FragmentName
	Directives   *Directives
	Locator
}

// ! Find the fragment definition by its name and return its selection set
func (f *FragmentSpread) Selections() *SelectionSet {
	return &SelectionSet{}
}

type InlineFragment struct {
	TypeCondition *TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
	Locator
}

func (i *InlineFragment) Selections() *SelectionSet {
	return &i.SelectionSet
}

type Type interface {
	LocatorInterface
}

var _ Type = (*NamedType)(nil)
var _ Type = (*ListType)(nil)
var _ Type = (*NonNullType)(nil)

type NamedType Name

type ListType struct {
	Kind   string
	OfType Type
	Locator
}

type NonNullType struct {
	Kind   string
	OfType Type
	Locator
}

type TypeCondition struct {
	NamedType NamedType
	Locator
}

type UnionMemberTypes []NamedType

type ImplementsInterfaces []NamedType

type Name struct {
	Value string
	Locator
}

type Alias Name
type FragmentName Name

type Value interface {
	GetValue() interface{}
	LocatorInterface
}

var _ Value = (*DefaultValue)(nil)
var _ Value = (*ObjectField)(nil)
var _ Value = (*ObjectValue)(nil)
var _ Value = (*ListValue)(nil)
var _ Value = (*IntValue)(nil)
var _ Value = (*FloatValue)(nil)
var _ Value = (*StringValue)(nil)
var _ Value = (*BooleanValue)(nil)
var _ Value = (*EnumValue)(nil)

type DefaultValue struct {
	Value Value
	Locator
}

func (dv DefaultValue) GetValue() interface{} {
	return dv.Value.GetValue()
}

type ObjectField struct {
	Name  Name
	Value Value
	Locator
}

func (of ObjectField) GetValue() interface{} {
	return of.Value.GetValue()
}

type ObjectValue struct {
	Values []ObjectField
	Locator
}

func (ov ObjectValue) GetValue() interface{} {
	return ov
}

type ListValue struct {
	Values []Value
	Locator
}

func (lv ListValue) GetValue() interface{} {
	return lv
}

type IntValue struct {
	Value int64
	Locator
}

func (iv IntValue) GetValue() interface{} {
	return iv
}

type FloatValue struct {
	Value float64
	Locator
}

func (fv FloatValue) GetValue() interface{} {
	return fv
}

type StringValue struct {
	Value string
	Locator
}

func (sv StringValue) GetValue() interface{} {
	return sv
}

type BooleanValue struct {
	Value bool
	Locator
}

func (bv BooleanValue) GetValue() interface{} {
	return bv
}

type EnumValue struct {
	Name Name
	Locator
}

func (ev EnumValue) GetValue() interface{} {
	return ev
}

type NullValue struct {
	Locator
}

func (nv NullValue) GetValue() interface{} {
	return nv
}

type Variable struct {
	Name Name
	Locator
}

type Description StringValue
