package ast

type Definition interface {
	definition()
}

type TypeSystemExtension interface {
	Definition
}

type TypeSystemDefinition interface {
	Definition
}

type ExecutableDefinition interface {
	Definition
	executableDefinition()
}

type SchemaDefinition struct {
	Directives                   *Directives
	RootOperationTypeDefinitions []RootOperationTypeDefinition
}

type SchemaExtension struct {
	Directives                   *Directives
	RootOperationTypeDefinitions *[]RootOperationTypeDefinition
}

type OperationDefinition struct {
	OperationType       *OperationType
	Name                *name
	VariableDefinitions *VariableDefinitions
	Directives          *Directives
	SelectionSet        SelectionSet
}

type FragmentDefinition struct {
	FragmentName  FragmentName
	TypeCondition TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
}

type VariableDefinition struct {
	Variable     Variable
	Type         Type
	DefaultValue *DefaultValue
	Directives   *Directives
}

type VariableDefinitions []VariableDefinition

type TypeExtension interface {
	TypeSystemExtension
	// Make sure all TypeExtensions implement this.
	typeExtension()
}

type TypeDefinition interface {
	TypeSystemDefinition
	// Make sure all TypeDefinitions implement this.
	typeDefinition()
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

type InputValueDefinition struct {
	Description  *string
	Name         name
	DefaultValue *DefaultValue
	Directives   *Directives
}

type DirectiveDefinition struct {
	Description         *string
	Name                name
	ArgumentsDefinition ArgumentsDefinition
	DirectiveLocations  DirectiveLocations
}

type ArgumentsDefinition []InputValueDefinition

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

type Variable struct {
	Name name
}
