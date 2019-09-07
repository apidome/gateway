package jsonwalker

import "fmt"

type JsonArrayIndexError struct {
	err   string
	index int
}

func (e JsonArrayIndexError) Error() string {
	return fmt.Sprintf("Index %d %s", e.index, e.err)
}

type JsonPointerSyntaxError struct {
	err  string
	path string
}

func (e JsonPointerSyntaxError) Error() string {
	return fmt.Sprintf("JsonPointer syntax error for \"%s\" - %s", e.path, e.err)
}
