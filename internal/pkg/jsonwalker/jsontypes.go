package jsonwalker

import "encoding/json"

/*****************/
/** Json Object **/
/*****************/

type JsonObject struct {
	data map[string]json.RawMessage
}

func NewJsonObject(data []byte) (*JsonObject, error) {
	object := JsonObject{
		make(map[string]json.RawMessage),
	}

	err := json.Unmarshal(data, &object.data)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

func (jo *JsonObject) GetObject(key string) (*JsonObject, error) {
	var object JsonObject

	err := json.Unmarshal(jo.data[key], &object.data)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

func (jo *JsonObject) GetString(key string) (*string, error) {
	var _string string

	err := json.Unmarshal(jo.data[key], &_string)
	if err != nil {
		return nil, err
	}

	return &_string, nil
}

func (jo *JsonObject) GetInteger(key string) (*int, error) {
	var integer int

	err := json.Unmarshal(jo.data[key], &integer)
	if err != nil {
		return nil, err
	}

	return &integer, nil
}

func (jo *JsonObject) GetBoolean(key string) (*bool, error) {
	var boolean bool

	err := json.Unmarshal(jo.data[key], &boolean)
	if err != nil {
		return nil, err
	}

	return &boolean, nil
}

func (jo *JsonObject) GetArray(key string) (*JsonArray, error) {
	var array JsonArray

	err := json.Unmarshal(jo.data[key], &array)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

/*****************/
/** Json Array **/
/*****************/

type JsonArray struct {
	data []json.RawMessage
}

func NewJsonArray(data []byte) (*JsonArray, error) {
	var array JsonArray

	err := json.Unmarshal(data, &array.data)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

func (ja *JsonArray) GetObject(index int) (*JsonObject, error) {
	var object JsonObject

	err := json.Unmarshal(ja.data[index], &object.data)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

func (ja *JsonArray) GetString(index int) (*string, error) {
	var _string string

	err := json.Unmarshal(ja.data[index], &_string)
	if err != nil {
		return nil, err
	}

	return &_string, nil
}

func (ja *JsonArray) GetInteger(index int) (*int, error) {
	var integer int

	err := json.Unmarshal(ja.data[index], &integer)
	if err != nil {
		return nil, err
	}

	return &integer, nil
}

func (ja *JsonArray) GetBoolean(index int) (*bool, error) {
	var boolean bool

	err := json.Unmarshal(ja.data[index], &boolean)
	if err != nil {
		return nil, err
	}

	return &boolean, nil
}

func (ja *JsonArray) GetArray(index int) (*JsonArray, error) {
	var array JsonArray

	err := json.Unmarshal(ja.data[index], &array)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

//
//
///*****************/
///** Json Number **/
///*****************/
//
//type JsonNumber struct {
//	data json.RawMessage
//}
//
//func NewJsonNumber(data []byte) (interface{}, error) {
//	// TODO: Unmarshal here
//	return nil, nil
//}
//
//
///******************/
///** Json Boolean **/
///******************/
//
//type JsonBoolean struct {
//	data json.RawMessage
//}
//
//func NewJsonBoolean(data []byte) (interface{}, error) {
//	// TODO: Unmarshal here
//	return nil, nil
//}
//
//
///*****************/
///** Json String **/
///*****************/
//
//type JsonString struct {
//	data json.RawMessage
//}
//
//func NewJsonString(data []byte) (interface{}, error) {
//	// TODO: Unmarshal here
//	return nil, nil
//}
//
//
///******************/
///** Json Integer **/
///******************/
//
//type JsonInteger struct {
//	data json.RawMessage
//}
//
//func NewJsonInteger(data []byte) (interface{}, error) {
//	// TODO: Unmarshal here
//	return nil, nil
//}
//
//
///***************/
///** Json Null **/
///***************/
//
//type JsonNull struct {
//	data json.RawMessage
//}
//
//func NewJsonNull(data []byte) (interface{}, error) {
//	// TODO: Unmarshal here
//	return nil, nil
//}
