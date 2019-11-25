package ast

import "github.com/omeryahud/caf/internal/pkg/graphql/language/location"

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
	loc                          *location.Location
}

type SchemaExtension struct {
	Directives                   *Directives
	RootOperationTypeDefinitions *[]RootOperationTypeDefinition
	loc                          *location.Location
}

type OperationDefinition struct {
	OperationType       *OperationType
	Name                *Name
	VariableDefinitions *VariableDefinitions
	Directives          *Directives
	SelectionSet        SelectionSet
	loc                 *location.Location
}

type FragmentDefinition struct {
	FragmentName  FragmentName
	TypeCondition TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
	loc           *location.Location
}

type VariableDefinition struct {
	Variable     Variable
	Type         Type
	DefaultValue *DefaultValue
	Directives   *Directives
	loc          *location.Location
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
	Name                 Name
	Directives           *Directives
	EnumValuesDefinition *EnumValuesDefinition
	loc                  *location.Location
}

type UnionTypeExtension struct {
	Name             Name
	Directives       *Directives
	UnionMemberTypes *UnionMemberTypes
	loc              *location.Location
}

type UnionMemberTypes []NamedType

type UnionTypeDefinition struct {
	Description      *string
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	loc              *location.Location
}

type FieldsDefinition []FieldDefinition

type FieldDefinition struct {
	Description         *string
	Name                Name
	ArgumentsDefinition *ArgumentsDefinition
	Type                Type
	Directives          *Directives
	loc                 *location.Location
}

type InterfaceTypeExtension struct {
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	loc              *location.Location
}

type InterfaceTypeDefinition struct {
	Description      *string
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	loc              *location.Location
}

type ImplementsInterfaces []NamedType

type ObjectTypeExtension struct {
	Name                 Name
	ImplementsInterfaces *ImplementsInterfaces
	Directive            *Directives
	FieldsDefinition     *FieldsDefinition
	loc                  *location.Location
}

type ObjectTypeDefinition struct {
	Description          *string
	Name                 Name
	ImplementsInterfaces *ImplementsInterfaces
	Directives           *Directives
	FieldsDefinition     *FieldsDefinition
	loc                  *location.Location
}

type ScalarTypeExtension struct {
	Name       Name
	Directives Directives
	loc        *location.Location
}

type ScalarTypeDefinition struct {
	Description *string
	Name        Name
	Directives  *Directives
	loc         *location.Location
}

type InputFieldsDefinition []InputValueDefinition

type InputObjectTypeExtension struct {
	Name                  Name
	Directives            *Directives
	InputFieldsDefinition *InputFieldsDefinition
	loc                   *location.Location
}

type InputObjectTypeDefinition struct {
	Description           *string
	Name                  Name
	Directives            *Directives
	InputFieldsDefinition *InputFieldsDefinition
	loc                   *location.Location
}

type EnumTypeExtension struct {
	Name                 Name
	Directives           *Directives
	EnumValuesDefinition *EnumValuesDefinition
	loc                  *location.Location
}

type EnumValuesDefinition []EnumValueDefinition

type EnumValueDefinition struct {
	Description *string
	EnumValue   EnumValue
	Directives  *Directives
	loc         *location.Location
}

type InputValueDefinition struct {
	Description  *string
	Name         Name
	DefaultValue *DefaultValue
	Directives   *Directives
	loc          *location.Location
}

type DirectiveDefinition struct {
	Description         *string
	Name                Name
	ArgumentsDefinition ArgumentsDefinition
	DirectiveLocations  DirectiveLocations
	loc                 *location.Location
}

type ArgumentsDefinition []InputValueDefinition

type RootOperationTypeDefinition struct {
	OperationType OperationType
	NamedType     NamedType
	loc           *location.Location
}

type OperationType string

const (
	OPERATION_QUERY        OperationType = "query"
	OPERATION_MUTATION     OperationType = "mutation"
	OPERATION_SUBSCRIPTION OperationType = "subscription"
)

type Variable struct {
	Name Name
	loc  *location.Location
}
