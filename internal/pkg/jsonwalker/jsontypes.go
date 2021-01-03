package jsonwalker

import (
	"encoding/json"
)

/*****************/
/** Json Object **/
/*****************/

// JsonObject is a type that represents a json object by keeping key-value
// pairs of string and json.RawMessage
type JsonObject map[string]json.RawMessage

// NewJsonObject initialize and returns a pointer to a new JsonObject
func NewJsonObject(data []byte) (*JsonObject, error) {
	object := new(JsonObject)

	err := json.Unmarshal(data, object)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// GetObject returns the value of an object property that holds an object.
func (jo JsonObject) GetObject(key string) (*JsonObject, error) {
	object, err := NewJsonObject([]byte{})
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jo[key], &object)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// GetString returns the value of an object property that holds a string.
func (jo JsonObject) GetString(key string) (*string, error) {
	var _string *string

	err := json.Unmarshal(jo[key], _string)
	if err != nil {
		return nil, err
	}

	return _string, nil
}

// GetInteger returns the value of an object property that holds a number.
func (jo JsonObject) GetNumber(key string) (*float64, error) {
	var number *float64

	err := json.Unmarshal(jo[key], number)
	if err != nil {
		return nil, err
	}

	return number, nil
}

// GetBoolean returns the value of an object property that holds a boolean.
func (jo JsonObject) GetBoolean(key string) (*bool, error) {
	var boolean *bool

	err := json.Unmarshal(jo[key], boolean)
	if err != nil {
		return nil, err
	}

	return boolean, nil
}

// GetArray returns the value of an object property that holds an array.
func (jo JsonObject) GetArray(key string) (*JsonArray, error) {
	var array *JsonArray

	err := json.Unmarshal(jo[key], array)
	if err != nil {
		return nil, err
	}

	return array, nil
}

/*****************/
/** Json Array **/
/*****************/

// JsonObject is a type that represents a json object by keeping a slice
// of json.RawMessages
type JsonArray []json.RawMessage

// NewJsonArray initializes and returns a pointer to JsonArray.
func NewJsonArray(data []byte) (*JsonArray, error) {
	var array *JsonArray

	err := json.Unmarshal(data, array)
	if err != nil {
		return nil, err
	}

	return array, nil
}

// GetObject returns the value of an object array item.
func (ja JsonArray) GetObject(index int) (*JsonObject, error) {
	if index > len(ja)-1 || index < 0 {
		return nil, JsonArrayIndexError(index)
	}

	object, err := NewJsonObject([]byte{})
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(ja[index], object)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// GetString returns the value of a string array item.
func (ja JsonArray) GetString(index int) (*string, error) {
	if index > len(ja)-1 || index < 0 {
		return nil, JsonArrayIndexError(index)
	}

	var _string *string

	err := json.Unmarshal(ja[index], _string)
	if err != nil {
		return nil, err
	}

	return _string, nil
}

// GetNumber returns the value of a numeric array item.
func (ja JsonArray) GetNumber(index int) (*float64, error) {
	if index > len(ja)-1 || index < 0 {
		return nil, JsonArrayIndexError(index)
	}

	var number *float64

	err := json.Unmarshal(ja[index], number)
	if err != nil {
		return nil, err
	}

	return number, nil
}

// GetBoolean returns the value of a boolean array item.
func (ja JsonArray) GetBoolean(index int) (*bool, error) {
	if index > len(ja)-1 || index < 0 {
		return nil, JsonArrayIndexError(index)
	}

	var boolean *bool

	err := json.Unmarshal(ja[index], boolean)
	if err != nil {
		return nil, err
	}

	return boolean, nil
}

// GetArray returns the value of an array array item.
func (ja JsonArray) GetArray(index int) (*JsonArray, error) {
	if index > len(ja)-1 || index < 0 {
		return nil, JsonArrayIndexError(index)
	}

	var array *JsonArray

	err := json.Unmarshal(ja[index], array)
	if err != nil {
		return nil, err
	}

	return array, nil
}

func (ja JsonArray) GetItemByIndex(index int) (interface{}, error) {
	if index > len(ja)-1 || index < 0 {
		return nil, JsonArrayIndexError(index)
	}

	return ja[index], nil
}
