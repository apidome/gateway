package jsonwalker

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
)

type JsonPointer []string

func NewJsonPointer(path string) (JsonPointer, error) {
	if len(path) == 0 {
		return JsonPointer{}, nil
	}

	if path[0] != '/' {
		// TODO: Create new error type.
		return nil, errors.New("first character of non-empty reference must be '/'")
	}

	tokens := strings.Split(path, "/")

	return JsonPointer(tokens), nil
}

func (jp JsonPointer) Evaluate(jsonData json.RawMessage) (interface{}, error) {
	// If the JsonPointer is an empty reference, return the whole data.
	if len(jp) == 0 {
		return jsonData, nil
	}

	var data interface{}

	//
	for i, token := range jp {
		if i == 0 {
			err := json.Unmarshal(jsonData, &data)
			if err != nil {
				return nil, err
			}
		} else {
			switch v := data.(type) {
			case bool, string, float64:
				{
					if len(jp) > i {
						// TODO: Create new error type.
						return nil, errors.New("json data does not contain the location the JsonPointer points to")
					}
				}
			case map[string]interface{}:
				{
					if len(jp) > i {
						err := json.Unmarshal(v[token].([]byte), &data)
						if err != nil {
							return nil, err
						}
					}
				}
			case []interface{}:
				{
					if len(jp) > i {
						index, err := strconv.Atoi(token)
						if err != nil {
							return nil, err
						}

						err = json.Unmarshal(v[index].([]byte), &data)
						if err != nil {
							return nil, err
						}
					}
				}
			default:
				{
					log.Print("[JsonPointer WARNING] Unexpected use case")
					// TODO: Create new error type.
					return nil, errors.New("unexpected possible data type in json data")
				}
			}
		}
	}

	return data, nil
}
