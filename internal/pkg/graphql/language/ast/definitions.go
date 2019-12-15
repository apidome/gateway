package ast

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/location"
)

type OperationType string

const (
	OPERATION_QUERY        OperationType = "query"
	OPERATION_MUTATION     OperationType = "mutation"
	OPERATION_SUBSCRIPTION OperationType = "subscription"
)

type LocatorInterface interface {
	Location() *location.Location
}

type Locator struct {
	Loc location.Location
}

func (l *Locator) Location() *location.Location {
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
	Directives          *Directives
	SelectionSet        SelectionSet
	Locator
}

type FragmentDefinition struct {
	FragmentName  FragmentName
	TypeCondition TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
	Locator
}

/*
	Type System Definitions
*/

type TypeSystemDefinition interface {
	Definition
}

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

type DirectiveDefinition struct {
	Description         *string
	Name                Name
	ArgumentsDefinition ArgumentsDefinition
	DirectiveLocations  DirectiveLocations
	Locator
}

type InputValueDefinition struct {
	Description  *string
	Name         Name
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
	Description *string
	Name        Name
	Directives  *Directives
	Locator
}

type ObjectTypeDefinition struct {
	Description          *string
	Name                 Name
	ImplementsInterfaces *ImplementsInterfaces
	Directives           *Directives
	FieldsDefinition     *FieldsDefinition
	Locator
}

type InterfaceTypeDefinition struct {
	Description      *string
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	Locator
}

type UnionTypeDefinition struct {
	Description      *string
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	Locator
}

type EnumTypeDefinition struct {
	Description          *string
	Name                 Name
	Directives           *Directives
	EnumValueDefinitions *EnumValueDefinitions
	Locator
}

type EnumValueDefinition struct {
	Description *string
	EnumValue   EnumValue
	Directives  *Directives
	Locator
}

type EnumValueDefinitions []EnumValueDefinition

type InputObjectTypeDefinition struct {
	Description           *string
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
	Description         *string
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
	RootOperationTypeDefinitions *[]RootOperationTypeDefinition
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
	Directive            *Directives
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
	EnumValueDefinitions *EnumValueDefinitions
	Locator
}

type InputObjectTypeExtension struct {
	Name                  Name
	Directives            *Directives
	InputFieldsDefinition *InputFieldsDefinition
	Locator
}
