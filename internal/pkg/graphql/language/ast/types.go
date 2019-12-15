package ast

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
