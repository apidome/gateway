package jsonvalidator

import (
	"encoding/json"
)

// JsonValidator is a struct that implements the Validator interface
// and validates json objects according to a json schema
type JsonValidator struct {
	schemaDict map[string]map[string]JsonSchema
}

// LoadSchema is a function that handles addition of new schema to the
// JsonValidator's schemas list
func (jv JsonValidator) LoadSchema(path, method, s string) error {
	var schema JsonSchema

	// Check if the string s is a valid json.
	err := json.Unmarshal([]byte(s), &schema)
	if err != nil {
		return err
	}

	isSchemaValid, err := validateJsonSchema(schema)
	if err != nil {
		return err
	}

	if isSchemaValid {
		// Create a new empty method-JsonSchema map for the current path.
		jv.schemaDict[path] = make(map[string]JsonSchema)

		// Add the schema to the appropriate map according to its path and
		// method.
		jv.schemaDict[path][method] = schema
	}

	return nil
}

// Parse converts a string that represents a json value to a known
// data structure
func (jv JsonValidator) Parse(b string) (bool, error) {
	return false, nil
}

// Validate is the function that actually perform validation of json value
// according to a specific json schema
func (jv JsonValidator) Validate(path, method, b string) (bool, error) {
	return false, nil
}

// NewJsonValidator returns a new instance of JsonValidator
func NewJsonValidator() JsonValidator {
	return JsonValidator{
		make(map[string]map[string]JsonSchema),
	}
}

// validateJsonSchema is a recursive function that validates the schema's
// structure according to Json Schema draft 7
func validateJsonSchema(schema JsonSchema) (bool, error) {
	return true, nil
}
