package validators

import "regexp"

type keywordValidator interface {
	validate(string, interface{}) (bool, error)
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

type pattern regexp.Regexp

func (p pattern) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

type format string

func (f format) validate(path string, jsonData interface{}) (bool, error) {
	return true, nil
}

/*********************/
/** Object Keywords **/
/*********************/

/*********************/
/** Number Keywords **/
/*********************/

/********************/
/** Array Keywords **/
/********************/
