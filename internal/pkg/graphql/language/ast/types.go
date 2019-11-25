package ast

import "github.com/omeryahud/caf/internal/pkg/graphql/language/location"

type TypeKind int

type Type interface {
	GetKind() string
}

type NamedType struct {
	Kind string
	Name Name
}

func (nt NamedType) GetKind() string {
	return nt.Kind
}

type ListType struct {
	Kind   string
	OfType Type
	Loc    *location.Location
}

func (lt ListType) GetKind() string {
	return lt.Kind
}

type NonNullType struct {
	Kind   string
	OfType Type
	Loc    *location.Location
}

func (nt NonNullType) GetKind() string {
	return nt.Kind
}

type TypeCondition struct {
	NamedType NamedType
	Loc       *location.Location
}

type UnionMemberTypes []NamedType

type ImplementsInterfaces []NamedType
