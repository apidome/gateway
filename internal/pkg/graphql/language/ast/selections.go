package ast

import "github.com/omeryahud/caf/internal/pkg/graphql/language/location"

type Selection interface {
	Selections() *SelectionSet
	Location() *location.Location
}

var _ Selection = (*Field)(nil)
var _ Selection = (*FragmentSpread)(nil)
var _ Selection = (*InlineFragment)(nil)

type SelectionSet []Selection

type Field struct {
	Alias        *Alias
	Name         Name
	Arguments    *Arguments
	Directives   *Directives
	SelectionSet *SelectionSet
	Locator
}

func (f *Field) Selections() *SelectionSet {
	return f.SelectionSet
}

type FragmentSpread struct {
	FragmentName FragmentName
	Directives   *Directives
	Locator
}

// ! Find the fragment definition by its name and return its selection set
func (f *FragmentSpread) Selections() *SelectionSet {
	return &SelectionSet{}
}

type InlineFragment struct {
	TypeCondition *TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
	Locator
}

func (i *InlineFragment) Selections() *SelectionSet {
	return &i.SelectionSet
}
