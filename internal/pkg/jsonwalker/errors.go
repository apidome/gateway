package jsonwalker

import "fmt"

type JsonArrayIndexError struct {
	err   string
	index int
}

func (e JsonArrayIndexError) Error() string {
	return fmt.Sprintf("Index %d %s", e.index, e.err)
}
