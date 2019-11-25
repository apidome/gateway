package ast

import "github.com/omeryahud/caf/internal/pkg/graphql/language/location"

type TypeKind int

const (
	NAMED_TYPE TypeKind = iota + 1
	LIST_TYPE
	NON_NULL_TYPE
)

type Type interface {
	GetTypeKind() TypeKind
}

type NamedType Name

func (nt NamedType) GetTypeKind() TypeKind {
	return NAMED_TYPE
}

type ListType struct {
	OfType Type
	loc    *location.Location
}

func (lt ListType) GetTypeKind() TypeKind {
	return LIST_TYPE
}

type NonNullType struct {
	OfType Type
	loc    *location.Location
}

func (nt NonNullType) GetTypeKind() TypeKind {
	return NON_NULL_TYPE
}

type TypeCondition struct {
	NamedType NamedType
	loc       *location.Location
}

type UnionMemberTypes []NamedType

type ImplementsInterfaces []NamedType
