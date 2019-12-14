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

type Definition interface {
	GetKind() string
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
	Kind                string
	OperationType       OperationType
	Name                *Name
	VariableDefinitions *VariableDefinitions
	Directives          *Directives
	SelectionSet        SelectionSet
	Loc                 *location.Location
}

func (od OperationDefinition) GetKind() string {
	return od.Kind
}

type FragmentDefinition struct {
	Kind          string
	FragmentName  FragmentName
	TypeCondition TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
	Loc           *location.Location
}

func (fd FragmentDefinition) GetKind() string {
	return fd.Kind
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
	Kind                         string
	Directives                   *Directives
	RootOperationTypeDefinitions []RootOperationTypeDefinition
	Loc                          *location.Location
}

func (sd SchemaDefinition) GetKind() string {
	return sd.Kind
}

type RootOperationTypeDefinition struct {
	Kind          string
	OperationType OperationType
	NamedType     NamedType
	Loc           *location.Location
}

func (rod RootOperationTypeDefinition) GetKind() string {
	return rod.Kind
}

type DirectiveDefinition struct {
	Kind                string
	Description         *string
	Name                Name
	ArgumentsDefinition ArgumentsDefinition
	DirectiveLocations  DirectiveLocations
	Loc                 *location.Location
}

func (dd DirectiveDefinition) GetKind() string {
	return dd.Kind
}

type InputValueDefinition struct {
	Kind         string
	Description  *string
	Name         Name
	DefaultValue *DefaultValue
	Directives   *Directives
	Loc          *location.Location
}

func (ivd InputValueDefinition) GetKind() string {
	return ivd.Kind
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
	Kind        string
	Description *string
	Name        Name
	Directives  *Directives
	Loc         *location.Location
}

func (sd ScalarTypeDefinition) GetKind() string {
	return sd.Kind
}

type ObjectTypeDefinition struct {
	Kind                 string
	Description          *string
	Name                 Name
	ImplementsInterfaces *ImplementsInterfaces
	Directives           *Directives
	FieldsDefinition     *FieldsDefinition
	Loc                  *location.Location
}

func (od ObjectTypeDefinition) GetKind() string {
	return od.Kind
}

type InterfaceTypeDefinition struct {
	Kind             string
	Description      *string
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	Loc              *location.Location
}

func (id InterfaceTypeDefinition) GetKind() string {
	return id.Kind
}

type UnionTypeDefinition struct {
	Kind             string
	Description      *string
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	Loc              *location.Location
}

func (ud UnionTypeDefinition) GetKind() string {
	return ud.Kind
}

type EnumTypeDefinition struct {
	Kind                 string
	Description          *string
	Name                 Name
	Directives           *Directives
	EnumValuesDefinition *EnumValuesDefinition
	Loc                  *location.Location
}

func (ed EnumTypeDefinition) GetKind() string {
	return ed.Kind
}

type EnumValueDefinition struct {
	Kind        string
	Description *string
	EnumValue   EnumValue
	Directives  *Directives
	Loc         *location.Location
}

func (evd EnumValueDefinition) GetKind() string {
	return evd.Kind
}

type EnumValuesDefinition []EnumValueDefinition

type InputObjectTypeDefinition struct {
	Kind                  string
	Description           *string
	Name                  Name
	Directives            *Directives
	InputFieldsDefinition *InputFieldsDefinition
	Loc                   *location.Location
}

func (iod InputObjectTypeDefinition) GetKind() string {
	return iod.Kind
}

type InputFieldsDefinition []InputValueDefinition

type VariableDefinition struct {
	Kind         string
	Variable     Variable
	Type         Type
	DefaultValue *DefaultValue
	Directives   *Directives
	Loc          *location.Location
}

type VariableDefinitions []VariableDefinition

type FieldDefinition struct {
	Kind                string
	Description         *string
	Name                Name
	ArgumentsDefinition *ArgumentsDefinition
	Type                Type
	Directives          *Directives
	Loc                 *location.Location
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
	Kind                         string
	Directives                   *Directives
	RootOperationTypeDefinitions *[]RootOperationTypeDefinition
	Loc                          *location.Location
}

func (se SchemaExtension) GetKind() string {
	return se.Kind
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
	Kind       string
	Name       Name
	Directives Directives
	Loc        *location.Location
}

func (se ScalarTypeExtension) GetKind() string {
	return se.Kind
}

type ObjectTypeExtension struct {
	Kind                 string
	Name                 Name
	ImplementsInterfaces *ImplementsInterfaces
	Directive            *Directives
	FieldsDefinition     *FieldsDefinition
	Loc                  *location.Location
}

func (oe ObjectTypeExtension) GetKind() string {
	return oe.Kind
}

type InterfaceTypeExtension struct {
	Kind             string
	Name             Name
	Directives       *Directives
	FieldsDefinition *FieldsDefinition
	Loc              *location.Location
}

func (ie InterfaceTypeExtension) GetKind() string {
	return ie.Kind
}

type UnionTypeExtension struct {
	Kind             string
	Name             Name
	Directives       *Directives
	UnionMemberTypes *UnionMemberTypes
	Loc              *location.Location
}

func (ue UnionTypeExtension) GetKind() string {
	return ue.Kind
}

type EnumTypeExtension struct {
	Kind                 string
	Name                 Name
	Directives           *Directives
	EnumValuesDefinition *EnumValuesDefinition
	Loc                  *location.Location
}

func (ee EnumTypeExtension) GetKind() string {
	return ee.Kind
}

type InputObjectTypeExtension struct {
	Kind                  string
	Name                  Name
	Directives            *Directives
	InputFieldsDefinition *InputFieldsDefinition
	Loc                   *location.Location
}

func (ioe InputObjectTypeExtension) GetKind() string {
	return ioe.Kind
}
