package jsonwalker

import (
	"errors"
	"strings"
)

type JsonPointer []string

func NewJsonPointer(path string) (JsonPointer, error) {
	if path[0] != '/' {
		return nil, errors.New("first character of JsonPointer must be \"/\"")
	}

	tokens := strings.Split(path, "/")

	return JsonPointer(tokens), nil
}
