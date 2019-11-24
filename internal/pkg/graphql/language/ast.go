package language

import (
	"flag"
	"github.com/omeryahud/caf/internal/pkg/graphql/language/kinds"
)

type DirectiveLocation int

const (
	// Executable Directive Locations
	QUERY DirectiveLocation = iota + 1
	MUTATION
	SUBSCRIPTION
	FIELD
	FRAGMENT_DEFINITION
	FRAGMENT_SPREAD
	INLINE_FRAGMENT
	VARIABLE_DEFINITION

	// Type System Directive Locations
	SCHEMA
	SCALAR
	OBJECT
	FIELD_DEFINITION
	ARGUMENT_DEFINITION
	INTERFACE
	UNION
	ENUM
	ENUM_VALUE
	INPUT_OBJECT
	INPUT_FIELD_DEFINITION
)

type TypeSystemDirectiveLocation DirectiveLocation

type ExecutableDirectiveLocation DirectiveLocation

type DirectiveLocations []DirectiveLocation

type name string

func parseName(n string) name {
	return name("")
}

type DirectiveDefinition struct {
	Description         *string
	Name                name
	ArgumentsDefinition ArgumentsDefinition
	DirectiveLocations  DirectiveLocations
}

type ArgumentsDefinition []InputValueDefinition

type TypeKind int

const (
	NAMED_TYPE TypeKind = iota + 1
	LIST_TYPE
	NON_NULL_TYPE
)

type Type interface {
	GetTypeKind() TypeKind
}

type NamedType name

func (nt NamedType) GetTypeKind() TypeKind {
	return NAMED_TYPE
}

type ListType struct {
	OfType Type
}

func (lt ListType) GetTypeKind() TypeKind {
	return LIST_TYPE
}

type NonNullType struct {
	OfType Type
}

func (nt NonNullType) GetTypeKind() TypeKind {
	return NON_NULL_TYPE
}

type InputValueDefinition struct {
	Description  *string
	Name         name
	DefaultValue *DefaultValue
	Directives   *Directives
}

type DefaultValue flag.Value

type Value interface {
	Kind() string
	Value() interface{}
}

type Directives []Directive

type Directive struct {
	Name      name
	Arguments *Arguments
}

type Arguments []Argument

type Argument struct {
	Name  name
	Value Value
}

type InputFieldsDefinition []InputValueDefinition

type InputObjectTypeExtension struct {
	Name                  name
	Directives            *Directives
	InputFieldsDefinition *InputFieldsDefinition
}

type InputObjectTypeDefinition struct {
	Description           *string
	Name                  name
	Directives            *Directives
	InputFieldsDefinition *InputFieldsDefinition
}

type EnumTypeExtension struct {
	Name                 name
	Directives           *Directives
	EnumValuesDefinition *EnumValuesDefinition
}

type EnumValuesDefinition []EnumValueDefinition

type EnumValueDefinition struct {
	Description *string
	EnumValue   EnumValue
	Directives  *Directives
}

type EnumValue struct {
	Name name
}

func (ev EnumValue) Kind() string {
	return kinds.EnumValue
}

func (ev EnumValue) Value() interface{} {
	return ev
}

type EnumTypeDefinition struct {
	Description          *string
	Name                 name
	Directives           *Directives
	EnumValuesDefinition *EnumValuesDefinition
}

type UnionTypeExtension struct {
	Name             name
	Directives       *Directives
	UnionMemberTypes *UnionMemberTypes
}

type UnionMemberTypes []NamedType

type UnionTypeDefinition struct {
	Description      *string
	Name             name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
}

type FieldsDefinition []FieldDefinition

type FieldDefinition struct {
	Description         *string
	Name                name
	ArgumentsDefinition *ArgumentsDefinition
	Type                Type
	Directives          *Directives
}

type InterfaceTypeExtension struct {
	Name             name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
}

type InterfaceTypeDefinition struct {
	Description      *string
	Name             name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
}

type ImplementsInterfaces []NamedType

type ObjectTypeExtension struct {
	Name                 name
	ImplementsInterfaces *ImplementsInterfaces
	Directive            *Directives
	FieldsDefinition     *FieldsDefinition
}

type ObjectTypeDefinition struct {
	Description          *string
	Name                 name
	ImplementsInterfaces *ImplementsInterfaces
	Directives           *Directives
	FieldsDefinition     *FieldsDefinition
}

type ScalarTypeExtension struct {
	Name       name
	Directives Directives
}

type ScalarTypeDefinition struct {
	Description *string
	Name        name
	Directives  *Directives
}

type TypeExtension interface {
	// Make sure all TypeExtensions implement this.
	typeExtension()
}

type TypeDefinition interface {
	// Make sure all TypeDefinitions implement this.
	typeDefinition()
}

type RootOperationTypeDefinition struct {
	OperationType OperationType
	NamedType     NamedType
}

type OperationType string

const (
	OPERATION_QUERY        OperationType = "query"
	OPERATION_MUTATION     OperationType = "mutation"
	OPERATION_SUBSCRIPTION OperationType = "subscription"
)

type SchemaExtension struct {
	Directives                  *Directives
	RootOperationTypeDefinition *RootOperationTypeDefinition
}

type TypeSystemExtension interface {
}

type TypeSystemDefinition interface {
}

type Variable struct {
	Name name
}

type VariableDefinition struct {
	Variable     Variable
	Type         Type
	DefaultValue *DefaultValue
	Directives   *Directives
}

type VariableDefinitions []VariableDefinition

type ObjectField struct {
	Name  name
	Value Value
}

type ObjectValue []ObjectField

func (ov ObjectValue) Kind() string {
	return kinds.ObjectValue
}

func (ov ObjectValue) Value() interface{} {
	return ov
}

type ListValue []Value

func (lv ListValue) Kind() string {
	return kinds.ListValue
}

func (lv ListValue) Value() interface{} {
	return lv
}

type IntValue struct {
	value int
}

func (iv IntValue) Kind() string {
	return kinds.IntValue
}

func (iv IntValue) Value() interface{} {
	return iv
}

type FloatValue struct {
	value float64
}

func (fv FloatValue) Kind() string {
	return kinds.FloatValue
}

func (fv FloatValue) Value() interface{} {
	return fv
}

type StringValue struct {
	value string
}

func (sv StringValue) Kind() string {
	return kinds.StringValue
}

func (sv StringValue) Value() interface{} {
	return sv
}

type BooleanValue struct {
	value bool
}

func (bv BooleanValue) Kind() string {
	return kinds.BooleanValue
}

func (bv BooleanValue) Value() interface{} {
	return bv
}

type TypeCondition struct {
	NamedType NamedType
}

type FragmentName name

type FragmentDefinition struct {
	FragmentName  FragmentName
	TypeCondition TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
}

type InlineFragment struct {
	TypeCondition *TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
}

func (inf InlineFragment) GetFields() []Field {
	return []Field{}
}

type FragmentSpread struct {
	FragmentName FragmentName
	Directives   *Directives
}

func (fs FragmentSpread) GetFields() []Field {
	return []Field{}
}

type Alias name

type Field struct {
	Alias        *Alias
	Name         name
	Arguments    *Arguments
	Directives   *Directives
	SelectionSet *SelectionSet
}

func (f Field) GetFields() []Field {
	return []Field{}
}

type Selection interface {
	GetFields() []Field
}

type SelectionSet []Selection

type OperationDefinition struct {
	OperationType       *OperationType
	Name                *name
	VariableDefinitions *VariableDefinitions
	Directives          *Directives
	SelectionSet        SelectionSet
}

type ExecutableDefinition interface {
	Definition
	executableDefinition()
}

type Definition interface {
	definition()
}

type Document struct {
	Definitions []Definition
}
