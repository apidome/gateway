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

/********************/
/** Array Keywords **/
/********************/
