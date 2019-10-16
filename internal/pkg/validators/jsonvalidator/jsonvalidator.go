package jsonvalidator

import (
	"github.com/Creespye/caf/internal/pkg/configs"
	"github.com/pkg/errors"
)

// JsonValidator is a struct that implements the Validator interface
// and validates json objects according to a json schema
type JsonValidator struct {
	draft      string
	schemaDict map[string]map[string]*RootJsonSchema
}

// NewJsonValidator returns a new instance of JsonValidator
func NewJsonValidator(draft string) *JsonValidator {
	return &JsonValidator{
		draft,
		make(map[string]map[string]*RootJsonSchema),
	}
}

// LoadSchema is a function that handles addition of new schema to the
// JsonValidator's schemas list
func (jv JsonValidator) LoadSchema(path, method string, rawSchema []byte) error {
	// Validate the given schema against draft-07 meta-schema.
	err := validateJsonSchema(jv.draft, rawSchema)
	if err != nil {
		return errors.Wrap(err, "validation against meta-schema failed")
	}

	// If the schema is valid make a new map and insert the new schema to it.
	if jv.schemaDict[path] == nil {
		// Create a new empty method-JsonSchema map for the current path.
		jv.schemaDict[path] = make(map[string]*RootJsonSchema)
	}

	// Create a new JsonSchema object.
	schema, err := NewRootJsonSchema(rawSchema)
	if err != nil {
		return errors.Wrap(err, "failed to create a RootJsonSchema instance")
	}

	// Add the schema to the appropriate map according to its path and
	// method.
	jv.schemaDict[path][method] = schema

	return nil
}

// Validate is the function that actually perform validation of json value
// according to a specific json schema
func (jv JsonValidator) Validate(path string, method string, body []byte) error {
	return jv.schemaDict[path][method].validateBytes(body)
}

// validateJsonSchema is a function that validates the schema's
// structure according to Json Schema.
func validateJsonSchema(draft string, rawSchema []byte) error {
	config, err := configs.GetConfiguration()
	if err != nil {
		return errors.Wrap(err, "could not access configuration module")
	}

	if rawMetaSchema, ok := config.General.JsonMetaSchema[draft]; ok {
		metaSchema, err := NewRootJsonSchema([]byte(rawMetaSchema))
		if err != nil {
			return errors.Wrap(err, "failed to create a RootJsonSchema instance for meta-schema - "+draft)
		}

		return metaSchema.validateBytes(rawSchema)
	} else {
		return InvalidDraftError(draft)
	}
}
