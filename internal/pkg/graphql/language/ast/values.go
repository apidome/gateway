package ast

import "github.com/omeryahud/caf/internal/pkg/graphql/language/kinds"

type name string

func parseName(n string) name {
	return name("")
}

type Alias name

type FragmentName name

type Value interface {
	Kind() string
	Value() interface{}
}

type DefaultValue Value

type ObjectField struct {
	Name  name
	Value Value
}

type ObjectValue []ObjectField

func (ov ObjectValue) Kind() string {
	return kinds.ObjectValue
}

func (ov ObjectValue) Value() interface{} {
	return ov
}

type ListValue []Value

func (lv ListValue) Kind() string {
	return kinds.ListValue
}

func (lv ListValue) Value() interface{} {
	return lv
}

type IntValue struct {
	value int
}

func (iv IntValue) Kind() string {
	return kinds.IntValue
}

func (iv IntValue) Value() interface{} {
	return iv
}

type FloatValue struct {
	value float64
}

func (fv FloatValue) Kind() string {
	return kinds.FloatValue
}

func (fv FloatValue) Value() interface{} {
	return fv
}

type StringValue struct {
	value string
}

func (sv StringValue) Kind() string {
	return kinds.StringValue
}

func (sv StringValue) Value() interface{} {
	return sv
}

type BooleanValue struct {
	value bool
}

func (bv BooleanValue) Kind() string {
	return kinds.BooleanValue
}

func (bv BooleanValue) Value() interface{} {
	return bv
}

type EnumValue struct {
	Name name
}

func (ev EnumValue) Kind() string {
	return kinds.EnumValue
}

func (ev EnumValue) Value() interface{} {
	return ev
}
