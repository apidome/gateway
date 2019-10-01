package jsonvalidator

import "github.com/Creespye/caf/internal/pkg/configs"

// JsonValidator is a struct that implements the Validator interface
// and validates json objects according to a json schema
type JsonValidator struct {
	schemaDict map[string]map[string]*RootJsonSchema
}

// LoadSchema is a function that handles addition of new schema to the
// JsonValidator's schemas list
func (jv JsonValidator) LoadSchema(path, method string, rawSchema []byte) error {
	// Validate the given schema against draft-07 meta-schema.
	isSchemaValid, err := validateJsonSchema(rawSchema)
	if err != nil {
		return err
	}

	// If the schema is valid make a new map and insert the new schema to it.
	if isSchemaValid {
		if jv.schemaDict[path] == nil {
			// Create a new empty method-JsonSchema map for the current path.
			jv.schemaDict[path] = make(map[string]*RootJsonSchema)
		}

		// Create a new JsonSchema object.
		schema, err := NewRootJsonSchema(rawSchema)
		if err != nil {
			return err
		}

		// Add the schema to the appropriate map according to its path and
		// method.
		jv.schemaDict[path][method] = schema
	}

	return nil
}

// Validate is the function that actually perform validation of json value
// according to a specific json schema
func (jv JsonValidator) Validate(path string, method string, body []byte) (bool, error) {
	return jv.schemaDict[path][method].validate(body)
}

// NewJsonValidator returns a new instance of JsonValidator
func NewJsonValidator() JsonValidator {
	return JsonValidator{
		make(map[string]map[string]*RootJsonSchema),
	}
}

// validateJsonSchema is a function that validates the schema's
// structure according to Json Schema.
func validateJsonSchema(rawSchema []byte) (bool, error) {
	config, err := configs.GetConfiguration()
	if err != nil {
		return false, err
	}

	metaSchema, err := NewRootJsonSchema([]byte(config.General.JsonMetaSchema["draft-07"]))
	if err != nil {
		return false, err
	}

	return metaSchema.validate(rawSchema)
}
