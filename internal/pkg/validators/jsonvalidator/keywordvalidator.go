package jsonvalidator

import (
	"encoding/json"
)

type keywordValidator interface {
	validate(string, interface{}) (bool, error)
}

/**********************/
/** Generic Keywords **/
/**********************/

type schema string

func (s *schema) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type ref string

func (r *ref) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type id string

func (i *id) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type comment string

func (c *comment) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type title string

func (t *title) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type description string

func (d *description) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type examples []interface{}

func (e examples) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type enum []interface{}

func (e enum) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type _default json.RawMessage

func (d _default) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

func (d *_default) UnmarshalJSON(data []byte) error {
	*d = data
	return nil
}

type _const json.RawMessage

func (c _const) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type definitions map[string]*JsonSchema

func (d definitions) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type _type json.RawMessage

func (t *_type) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

func (t *_type) UnmarshalJSON(data []byte) error {
	*t = data
	return nil
}

/*********************/
/** String Keywords **/
/*********************/

type minLength int

func (ml *minLength) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type maxLength int

func (ml *maxLength) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type pattern string

func (p *pattern) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type format string

func (f *format) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

/*********************/
/** Number Keywords **/
/*********************/

type multipleOf int

func (mo *multipleOf) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type minimum float64

func (m *minimum) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type maximum float64

func (m *maximum) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type exclusiveMinimum float64

func (em *exclusiveMinimum) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type exclusiveMaximum float64

func (em *exclusiveMaximum) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

/*********************/
/** Object Keywords **/
/*********************/

type properties map[string]*JsonSchema

func (p properties) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type additionalProperties json.RawMessage

func (ap additionalProperties) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

func (ap *additionalProperties) UnmarshalJSON(data []byte) error {
	*ap = data
	return nil
}

type required []string

func (r required) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type propertyNames map[string]interface{}

func (pn propertyNames) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type dependencies map[string]interface{}

func (d dependencies) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type patternProperties map[string]interface{}

func (pp patternProperties) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type minProperties int

func (mp *minProperties) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type maxProperties int

func (mp *maxProperties) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

/********************/
/** Array Keywords **/
/********************/

type items json.RawMessage

func (i items) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

func (i *items) UnmarshalJSON(data []byte) error {
	*i = data
	return nil
}

type contains json.RawMessage

func (c contains) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

func (c *contains) UnmarshalJSON(data []byte) error {
	*c = data
	return nil
}

type additionalItems json.RawMessage

func (ai additionalItems) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

func (ai *additionalItems) UnmarshalJSON(data []byte) error {
	*ai = data
	return nil
}

type minItems int

func (mi *minItems) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type maxItems int

func (mi *maxItems) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type uniqueItems bool

func (ui *uniqueItems) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

/********************/
/** Other Keywords **/
/********************/

type contentMediaType string

func (cm *contentMediaType) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type contentEncoding string

func (ce *contentEncoding) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

/**************************/
/** Conditional Keywords **/
/**************************/

type anyOf []*JsonSchema

func (af anyOf) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type allOf []*JsonSchema

func (af allOf) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type oneOf []*JsonSchema

func (of oneOf) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type not JsonSchema

func (n *not) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type _if JsonSchema

func (i *_if) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type _then JsonSchema

func (t *_then) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type _else JsonSchema

func (e *_else) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}
