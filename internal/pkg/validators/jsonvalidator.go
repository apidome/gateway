package validators

import "encoding/json"

type JsonSchema map[string]interface{}

var jsonSchemaKeywords = [...]string{
	"$id",
	"$schema",
	"$ref",
	"title",
	"description",
	"default",
	"examples",
	"enum",
	"properties",
	"definitions",
	"required",
}

var jsonSchemaTypes = [...]string{
	"object",
	"array",
	"string",
	"number",
	"boolean",
	"null",
}

// JsonValidator is a struct that implements the Validator interface
// and validates json objects according to a json schema
type JsonValidator struct {
	schemaList map[string][]JsonSchema
}

// LoadSchema is a function that handles addition of new schema to the
// JsonValidator's schemas list
func (jv JsonValidator) LoadSchema(path string, s string) error {
	var schema JsonSchema

	// Check if the string s is a valid json.
	err := json.Unmarshal([]byte(s), &schema)
	if err != nil {
		return err
	}

	// TODO: Continue writing validation code here.

	// Add the schema to the
	jv.schemaList[path] = append(jv.schemaList[path], schema)

	return nil
}

// Parse converts a string that represents a json value to a known
// data structure
func (jv JsonValidator) Parse(b string) (bool, error) {
	return false, nil
}

// Validate is the function that actually perform validation of json value
// according to a specific json schema
func (jv JsonValidator) Validate(b string) (bool, error) {
	return false, nil
}

// NewJsonValidator returns a new instance of JsonValidator
func NewJsonValidator() JsonValidator {
	return JsonValidator{}
}
