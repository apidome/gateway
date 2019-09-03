package jsonvalidator

type keywordValidator interface {
	validate(string, interface{}) (bool, error)
}

/**********************/
/** Generic Keywords **/
/**********************/

type schema string

func (s schema) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type ref string

func (r ref) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type id string

func (i id) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type comment string

func (c comment) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type title string

func (t title) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type description string

func (d description) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type _default interface{}
type examples []interface{}
type enum []interface{}
type _const interface{}
type definitions map[string]*JsonSchema

/*********************/
/** String Keywords **/
/*********************/

type minLength int

func (ml minLength) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type maxLength int

func (ml maxLength) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type pattern string

func (p pattern) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type format string

func (f format) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

/*********************/
/** Number Keywords **/
/*********************/

type multipleOf int

func (mo multipleOf) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type minimum float64

func (m minimum) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type maximum float64

func (m maximum) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type exclusiveMinimum float64

func (em exclusiveMinimum) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type exclusiveMaximum float64

func (em exclusiveMaximum) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

/*********************/
/** Object Keywords **/
/*********************/

type properties map[string]*JsonSchema

func (p properties) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type additionalProperties interface{}
type required []string
type propertyNames map[string]interface{}
type dependencies map[string]interface{}
type patternProperties map[string]interface{}
type minProperties int
type maxProperties int

/********************/
/** Array Keywords **/
/********************/

type items interface{}
type contains interface{}
type additionalItems interface{}
type minItems int
type maxItems int
type uniqueItems bool

/********************/
/** Other Keywords **/
/********************/

type contentMediaType string
type contentEncoding string

/**************************/
/** Conditional Keywords **/
/**************************/

type anyOf []*JsonSchema
type allOf []*JsonSchema
type oneOf []*JsonSchema
type not *JsonSchema
type _if *JsonSchema
type _then *JsonSchema
type _else *JsonSchema
