package ast

import "github.com/omeryahud/caf/internal/pkg/graphql/language/location"

type Selection interface {
	GetFields() []Field
}

type SelectionSet []Selection

type Field struct {
	Alias        *Alias
	Name         Name
	Arguments    *Arguments
	Directives   *Directives
	SelectionSet *SelectionSet
	loc          *location.Location
}

func (f Field) GetFields() []Field {
	return []Field{}
}

type FragmentSpread struct {
	FragmentName FragmentName
	Directives   *Directives
	loc          *location.Location
}

func (fs FragmentSpread) GetFields() []Field {
	return []Field{}
}

type InlineFragment struct {
	TypeCondition *TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
	loc           *location.Location
}

func (inf InlineFragment) GetFields() []Field {
	return []Field{}
}
