package jsonvalidator

import "testing"

func TestNewJsonValidator(t *testing.T) {
	drafts := []string{"draft-07", "draft-06", ""}

	for _, draft := range drafts {
		if jsonValidator := NewJsonValidator(draft); jsonValidator == nil {
			t.Error("cannot create JsonValidator with draft " + draft)
		}
	}
}

func TestLoadSchema(t *testing.T) {

}

func TestValidate(t *testing.T) {

}
