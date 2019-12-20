package language

import "github.com/pkg/errors"

type syntaxError error

func newSyntaxError(msg string) syntaxError {
	return errors.New(msg)
}

type nonExistingObject error

func newNonExistingObject(msg string) nonExistingObject {
	return errors.New(msg)
}
