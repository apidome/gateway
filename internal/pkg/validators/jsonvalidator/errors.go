package jsonvalidator

import "fmt"

type KeywordValidationError struct {
	keyword string
	path    string
	reason  string
}

func (e KeywordValidationError) Error() string {
	return fmt.Sprintf(e.keyword + " validation failed in path " + e.path + ": " + e.reason)
}
