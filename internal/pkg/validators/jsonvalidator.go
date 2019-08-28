package validators

// JsonValidator is a struct that implements the Validator interface
// and validates json objects according to a json schema
type JsonValidator struct {
	schema []byte
}

// LoadSchema is a function that handles addition of new schema to the
// JsonValidator's schemas list
func (jv JsonValidator) LoadSchema(s []byte) error {
	return nil
}

// Parse converts a string that represents a json value to a known
// data structure
func (jv JsonValidator) Parse(b []byte) (bool, error) {
	return false, nil
}

// Validate is the function that actually perform validation of json value
// according to a specific json schema
func (jv JsonValidator) Validate(b []byte) (bool, error) {
	return false, nil
}

// NewJsonValidator returns a new instance of JsonValidator
func NewJsonValidator(s []byte) JsonValidator {
	return JsonValidator{}
}
