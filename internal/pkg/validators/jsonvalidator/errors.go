package jsonvalidator

import "fmt"

type KeywordValidationError struct {
	keyword string
	reason  string
}

func (e KeywordValidationError) Error() string {
	return fmt.Sprintf("\"" + e.keyword + "\" validation failed, reason: " + e.reason)
}

type SchemaCompilationError struct {
	path string
	err  string
}

func (e SchemaCompilationError) Error() string {
	return fmt.Sprintf("schema compilation failed in path " + e.path + ": " + e.err)
}
