package ast

type Name struct {
	Value string
	Locator
}

func ParseName(n string) Name {
	// TODO: Implement Name validation according to
	// 	https://graphql.github.io/graphql-spec/draft/#Name
	return Name{}
}

type Alias Name
type FragmentName Name

type Value interface {
	GetValue() interface{}
	LocatorInterface
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
	Value Value
	Locator
}

func (dv DefaultValue) GetValue() interface{} {
	return dv.Value.GetValue()
}

type ObjectField struct {
	Name  Name
	Value Value
	Locator
}

func (of ObjectField) GetValue() interface{} {
	return of.Value.GetValue()
}

type ObjectValue struct {
	Values []ObjectField
	Locator
}

func (ov ObjectValue) GetValue() interface{} {
	return ov
}

type ListValue struct {
	Values []Value
	Locator
}

func (lv ListValue) GetValue() interface{} {
	return lv
}

type IntValue struct {
	Value int64
	Locator
}

func (iv IntValue) GetValue() interface{} {
	return iv
}

type FloatValue struct {
	Value float64
	Locator
}

func (fv FloatValue) GetValue() interface{} {
	return fv
}

type StringValue struct {
	Value string
	Locator
}

func (sv StringValue) GetValue() interface{} {
	return sv
}

type BooleanValue struct {
	Value bool
	Locator
}

func (bv BooleanValue) GetValue() interface{} {
	return bv
}

type EnumValue struct {
	Name Name
	Locator
}

func (ev EnumValue) GetValue() interface{} {
	return ev
}

type NullValue struct {
	Locator
}

func (nv NullValue) GetValue() interface{} {
	return nv
}

type Variable struct {
	Name Name
	Locator
}
