package jsonwalker

import (
	"errors"
	"strings"
)

type JsonPointer []string

func NewJsonPointer(path string) (JsonPointer, error) {
	if len(path) == 0 {
		return JsonPointer{}, nil
	}

	if path[0] != '/' {
		return nil, errors.New("first character of non-empty reference must be \"/\"")
	}

	tokens := strings.Split(path, "/")

	return JsonPointer(tokens), nil
}
