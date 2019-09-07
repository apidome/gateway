package jsonwalker

import "fmt"

type JsonArrayIndexError int

func (e JsonArrayIndexError) Error() string {
	return fmt.Sprintf("Index %d out of range", e)
}

type JsonPointerSyntaxError struct {
	err  string
	path string
}

func (e JsonPointerSyntaxError) Error() string {
	return fmt.Sprintf("JsonPointer syntax error for \"%s\" - %s", e.path, e.err)
}

type InvalidJsonPointerError struct {
	path         string
	missingToken string
}

func (e InvalidJsonPointerError) Error() string {
	return fmt.Sprintf("invalid json pointer \"%s\": missing json token - \"%s\"", e.path, e.missingToken)
}
