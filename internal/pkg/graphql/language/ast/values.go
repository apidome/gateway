package ast

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/kinds"
	"github.com/omeryahud/caf/internal/pkg/graphql/language/location"
)

type Name struct {
	value string
	loc   location.Location
}

func parseName(n string) Name {
	return Name{}
}

type Alias Name

type FragmentName Name

type Value interface {
	Kind() string
	Value() interface{}
}

type DefaultValue Value

type ObjectField struct {
	Name  Name
	Value Value
	loc   *location.Location
}

type ObjectValue struct {
	Values []ObjectField
	loc    *location.Location
}

func (ov ObjectValue) Kind() string {
	return kinds.ObjectValue
}

func (ov ObjectValue) Value() interface{} {
	return ov
}

type ListValue struct {
	Values []Value
	loc    *location.Location
}

func (lv ListValue) Kind() string {
	return kinds.ListValue
}

func (lv ListValue) Value() interface{} {
	return lv
}

type IntValue struct {
	value int
	loc   *location.Location
}

func (iv IntValue) Kind() string {
	return kinds.IntValue
}

func (iv IntValue) Value() interface{} {
	return iv
}

type FloatValue struct {
	value float64
	loc   *location.Location
}

func (fv FloatValue) Kind() string {
	return kinds.FloatValue
}

func (fv FloatValue) Value() interface{} {
	return fv
}

type StringValue struct {
	value string
	loc   *location.Location
}

func (sv StringValue) Kind() string {
	return kinds.StringValue
}

func (sv StringValue) Value() interface{} {
	return sv
}

type BooleanValue struct {
	value bool
	loc   *location.Location
}

func (bv BooleanValue) Kind() string {
	return kinds.BooleanValue
}

func (bv BooleanValue) Value() interface{} {
	return bv
}

type EnumValue struct {
	Name Name
	loc  *location.Location
}

func (ev EnumValue) Kind() string {
	return kinds.EnumValue
}

func (ev EnumValue) Value() interface{} {
	return ev
}

type Variable struct {
	Name Name
	loc  *location.Location
}
