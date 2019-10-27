package validators

type Validator interface {
	// LoadSchema Gets a new schema and verifies that the schema is correct.
	LoadSchema(path string, method string, schema []byte) error

	// Validate enforces the schema's rules on a piece of data.
	Validate(path string, method string, body []byte) error
}
