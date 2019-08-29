package validators

type Validator interface {
	// LoadSchema Gets a new schema and verifies that the schema is correct.
	LoadSchema(s string) error

	// Parser verifies that a piece of data fits to the validator's format.
	Parse(b string) (bool, error)

	// Validate enforces the schema's rules on a piece of data.
	Validate(b string) (bool, error)
}
