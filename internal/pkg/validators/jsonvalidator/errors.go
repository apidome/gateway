package jsonvalidator

import "fmt"

type KeywordValidationError struct {
	keyword string
	reason  string
}

func (e KeywordValidationError) Error() string {
	return fmt.Sprintf(e.keyword + " validation failed, reason: " + e.reason)
}
