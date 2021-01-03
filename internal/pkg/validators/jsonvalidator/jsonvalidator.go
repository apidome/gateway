package jsonvalidator

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
)

var (
	methods = []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	}
)

// JsonValidator is a struct that implements the Validator interface
// and validates json objects according to a json schema
type JsonValidator struct {
	draft      string
	schemaDict map[string]map[string]*RootJsonSchema
}

// NewJsonValidator returns a new instance of JsonValidator
func NewJsonValidator(draft string) (JsonValidator, error) {
	supportedDrafts := []string{"draft-07"}

	for _, supportedDraft := range supportedDrafts {
		if supportedDraft == draft {
			return JsonValidator{
				draft,
				make(map[string]map[string]*RootJsonSchema),
			}, nil
		}
	}

	return JsonValidator{}, InvalidDraftError(draft)
}

// LoadSchema is a function that handles addition of new schema to the
// JsonValidator's schemas list
func (jv JsonValidator) LoadSchema(path, method string, rawSchema []byte) error {
	// Check if the given method is correct
	for _, httpMethod := range methods {
		if method == httpMethod {
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
				return errors.Wrap(err, "failed to create a RootJsonSchema "+
					"instance")
			}

			// Add the schema to the appropriate map according to its path and
			// method.
			jv.schemaDict[path][method] = schema

			return nil
		}
	}

	return errors.New("could not load schema to path " +
		path +
		": unknown method \"" +
		method +
		"\"")
}

// Validate is the function that actually perform validation of json value
// according to a specific json schema
func (jv JsonValidator) Validate(path string, method string, body []byte) error {
	if _, isPathExist := jv.schemaDict[path]; isPathExist {
		if _, isMethodExist := jv.schemaDict[path][method]; isMethodExist {
			return jv.schemaDict[path][method].validateBytes(body)
		} else {
			return errors.New("could not validate request: unknown path \"" +
				path +
				"\"")
		}
	} else {
		return errors.New("could not validate to path " +
			path +
			": no schema exist for method \"" +
			method +
			"\"")
	}
}

// validateJsonSchema is a function that validates the schema's
// structure according to Json Schema.
func validateJsonSchema(draft string, rawSchema []byte) error {
	// Get the path of the current go file (including the path inside
	// the project).
	var absolutePath string
	if _, filename, _, ok := runtime.Caller(0); ok {
		absolutePath = path.Dir(filename)
	}

	// Open the meta-schema file.
	file, err := os.Open(absolutePath + "/meta-schemas/" + draft)
	if err != nil {
		return errors.Wrap(err, "json schema version \""+
			draft+
			"\" is not supported")
	}

	defer file.Close()

	// Read the data from the file.
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "could not read meta-schema from file")
	}

	// Create a new RootJsonSchema.
	metaSchema, err := NewRootJsonSchema(bytes)
	if err != nil {
		return errors.Wrap(err, "failed to create a RootJsonSchema instance "+
			"for meta-schema - "+
			draft)
	}

	return metaSchema.validateBytes(rawSchema)
}
