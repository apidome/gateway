package jsonwalker

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type JsonPointer []string

func NewJsonPointer(path string) (JsonPointer, error) {
	if len(path) == 0 {
		return JsonPointer{}, nil
	}

	if path[0] != '/' {
		return nil, JsonPointerSyntaxError{
			"first character of non-empty reference must be '/'",
			path,
		}
		//return nil, errors.New("first character of non-empty reference must be '/'")
	}

	tokens := strings.Split(path, "/")

	return JsonPointer(tokens[1:]), nil
}

func (jp JsonPointer) Evaluate(jsonData json.RawMessage) (interface{}, error) {
	// If the JsonPointer is an empty reference, return the whole data.
	if len(jp) == 0 {
		return jsonData, nil
	}

	var data interface{}

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	//
	for _, token := range jp {
		data, err = jp.evaluateToken(token, data)
		if err != nil {
			return nil, errors.New("invalid json pointer - " + strings.Join(jp, "/") + ": " + err.Error())
		}
	}

	return data, nil
}

func (jp JsonPointer) evaluateToken(token string, jsonData interface{}) (interface{}, error) {
	switch v := jsonData.(type) {
	case map[string]interface{}:
		{
			return v[token], nil
		}
	case []interface{}:
		{
			index, err := strconv.Atoi(token)
			if err != nil {
				return nil, err
			}

			return v[index], nil
		}
	default:
		{
			// TODO: Create new error type.
			return nil, errors.New("json token - " + token + " does not exist in data")
		}
	}
}
