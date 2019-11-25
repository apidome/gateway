package ast

import (
	"github.com/omeryahud/caf/internal/pkg/graphql/language/location"
)

type Name struct {
	Value string
	Loc   location.Location
}

func ParseName(n string) Name {
	// TODO: Implement Name validation according to
	// 	https://graphql.github.io/graphql-spec/draft/#Name
	return Name{}
}

type Alias Name
type FragmentName Name

type Value interface {
	GetKind() string
	GetValue() interface{}
}

var _ Value = (*DefaultValue)(nil)
var _ Value = (*ObjectField)(nil)
var _ Value = (*ObjectValue)(nil)
var _ Value = (*ListValue)(nil)
var _ Value = (*IntValue)(nil)
var _ Value = (*FloatValue)(nil)
var _ Value = (*StringValue)(nil)
var _ Value = (*BooleanValue)(nil)
var _ Value = (*EnumValue)(nil)

type DefaultValue struct {
	Kind  string
	Value Value
	Loc   location.Location
}

func (dv DefaultValue) GetKind() string {
	return dv.Kind
}

func (dv DefaultValue) GetValue() interface{} {
	return dv.Value.GetValue()
}

type ObjectField struct {
	Kind  string
	Name  Name
	Value Value
	Loc   *location.Location
}

func (of ObjectField) GetKind() string {
	return of.Kind
}

func (of ObjectField) GetValue() interface{} {
	return of.Value.GetValue()
}

type ObjectValue struct {
	Kind   string
	Values []ObjectField
	Loc    *location.Location
}

func (ov ObjectValue) GetKind() string {
	return ov.Kind
}

func (ov ObjectValue) GetValue() interface{} {
	return ov
}

type ListValue struct {
	Kind   string
	Values []Value
	Loc    *location.Location
}

func (lv ListValue) GetKind() string {
	return lv.Kind
}

func (lv ListValue) GetValue() interface{} {
	return lv
}

type IntValue struct {
	Kind  string
	Value int
	Loc   *location.Location
}

func (iv IntValue) GetKind() string {
	return iv.Kind
}

func (iv IntValue) GetValue() interface{} {
	return iv
}

type FloatValue struct {
	Kind  string
	Value float64
	Loc   *location.Location
}

func (fv FloatValue) GetKind() string {
	return fv.Kind
}

func (fv FloatValue) GetValue() interface{} {
	return fv
}

type StringValue struct {
	Kind  string
	Value string
	Loc   *location.Location
}

func (sv StringValue) GetKind() string {
	return sv.Kind
}

func (sv StringValue) GetValue() interface{} {
	return sv
}

type BooleanValue struct {
	Kind  string
	Value bool
	Loc   *location.Location
}

func (bv BooleanValue) GetKind() string {
	return bv.Kind
}

func (bv BooleanValue) GetValue() interface{} {
	return bv
}

type EnumValue struct {
	Kind string
	Name Name
	Loc  *location.Location
}

func (ev EnumValue) GetKind() string {
	return ev.Kind
}

func (ev EnumValue) GetValue() interface{} {
	return ev
}

type Variable struct {
	Name Name
	Loc  *location.Location
}
