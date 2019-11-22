package language

// Documents is a type that holds a parsed graphql Document.
type Document struct {
	Operations []OperationDefinition
	Fragments  []FragmentDefinition
}

// Documents is a type that holds an OperationDefinition.
type OperationDefinition struct {
	SelectionSet
	Type       string
	Name       string
	Variables  []Variable
	Directives []Directive
}

// Documents is a type that holds an FragmentDefinition.
type FragmentDefinition struct {
	Fragment
	Name string
}

type Fragment struct {
	SelectionSet
	TypeCondition string
	Directives    []Directive
}

// Selection is an interface that whatever type that implements it,
// may be included in a SelectionSet (Field, FragmentSpread and
// InlineFragment).
type Selection interface {
	GetFields() []*Field
}

// SelectionSet is a list of Selection instances.
type SelectionSet []Selection

type Field struct {
	SelectionSet
	Alias      string
	Name       string
	Arguments  []Argument
	Directives []Directive
}

type FragmentSpread struct {
	Name       string
	Directives []Directive
}

type InlineFragment struct {
	Fragment
}

type Directive struct {
	Name      string
	Arguments []Argument
}

type Variable struct {
	Name         string
	Type         string
	DefaultValue string
}

type Argument struct {
	Name  string
	Value string
}
