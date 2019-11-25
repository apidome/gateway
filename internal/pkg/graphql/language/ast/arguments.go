package ast

import "github.com/omeryahud/caf/internal/pkg/graphql/language/location"

type Arguments []Argument

type Argument struct {
	Name  Name
	Value Value
	Loc   *location.Location
}
