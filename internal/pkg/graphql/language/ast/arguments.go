package ast

type Arguments []Argument

type Argument struct {
	Name  name
	Value Value
}
