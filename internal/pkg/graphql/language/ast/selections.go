package ast

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
}

func (f Field) GetFields() []Field {
	return []Field{}
}

type FragmentSpread struct {
	FragmentName FragmentName
	Directives   *Directives
}

func (fs FragmentSpread) GetFields() []Field {
	return []Field{}
}

type InlineFragment struct {
	TypeCondition *TypeCondition
	Directives    *Directives
	SelectionSet  SelectionSet
}

func (inf InlineFragment) GetFields() []Field {
	return []Field{}
}
