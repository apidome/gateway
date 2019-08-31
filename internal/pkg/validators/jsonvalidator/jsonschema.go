package validators

import "regexp"

type JsonSchema struct {
	// The $schema keyword is used to declare that a JSON fragment is
	// actually a piece of JSON Schema.
	_schema string

	// The value of $ref is a URI, and the part after # sign is in a format
	// called JSON Pointer.
	_ref string

	// The $id property is a URI that serves two purposes:
	// It declares a unique identifier for the schema
	// It declares a base URI against which $ref URIs are resolved.
	_id string

	// The $comment keyword is strictly intended for adding comments
	// to the JSON schema source. Its value must always be a string.
	_comment string

	// Title and Description used to describe the schema and not used for
	// validation.
	_title       string
	_description string

	// The default keyword specifies a default value for an item.
	_default interface{}

	// The examples keyword is a place to provide an array of examples
	// that validate against the schema.
	_examples []interface{}

	// The enum keyword is used to restrict a value to a fixed set of values.
	// It must be an array with at least one element, where each element
	// is unique.
	_enum []interface{}

	// The const keyword is used to restrict a value to a single value.
	_const interface{}

	// The definitions keyword is used to create entities that we recognize as
	// repetitive entities.
	// This ability maintains reuse in out Json Schema.
	_definitions map[string]*JsonSchema

	// The value of properties is an object, where each key is the name of a
	// property and each value is a JSON schema used to validate that property.
	_properties map[string]*JsonSchema

	// The additionalProperties keyword is used to control the handling of
	// extra stuff, that is, properties whose names are not listed in the
	// properties keyword.
	// By default any additional properties are allowed.
	// The additionalProperties keyword may be either a boolean or an object.
	// If additionalProperties is a boolean and set to false, no additional
	// properties will be allowed.
	// If additionalProperties is an object, that object is a schema that will be
	// used to validate any additional properties not listed in properties.
	_additionalProperties interface{}

	// The required keyword takes an array of zero or more strings.
	// Each of these strings must be unique.
	_required []string

	// The names of properties can be validated against a schema, irrespective
	// of their values.
	// This can be useful if you don’t want to enforce specific properties,
	// but you want to make sure that the names of those properties follow
	// a specific convention.
	_propertyNames map[string]interface{}

	// The dependencies keyword allows the schema of the object to change
	// based on the presence of certain special properties.
	_dependencies map[string]interface{}

	// TODO: Learn more about this keyword.
	_patternProperties map[string]interface{}

	// Items can be either an object or an array. If it is an object, it will
	// represent a schema that all the items in the array should match.
	// If it is an array, each item in that array is a different json schema
	// that should match the corresponding item in the inspected array
	// (In this case the index of each item is very important).
	_items interface{}

	// While the items schema must be valid for every item in the array,
	// the contains schema only needs to validate against one or more
	// items in the array.
	_contains interface{}

	// The additionalItems keyword controls whether it’s valid to have
	// additional items in the array beyond what is defined in items.
	_additionalItems interface{}

	// array limitations
	_minItems    int
	_maxItems    int
	_uniqueItems bool

	// string limitations
	_minLength int
	_maxLength int
	_pattern   regexp.Regexp
	_format    string

	// integer/number limitations
	_multipleOf       int
	_minimum          float64
	_maximum          float64
	_exclusiveMinimum float64
	_exclusiveMaximum float64

	// object size limitations
	_minProperties int
	_maxProperties int

	// The contentMediaType keyword specifies the MIME type of the contents
	// of a string.
	_contentMediaType string

	// The contentEncoding keyword specifies the encoding used to store
	// the contents.
	_contentEncoding string

	// Must be valid against any of the sub-schemas.
	_anyOf []*JsonSchema

	// Must be valid against all of the sub-schemas.
	_allOf []*JsonSchema

	// Must be valid against exactly one of the sub-schemas.
	_oneOf []*JsonSchema

	// Must not be valid against the given schema.
	_not *JsonSchema

	// The if, then and else keywords allow the application of a sub-schema
	// based on the outcome of another schema.
	_if   *JsonSchema
	_then *JsonSchema
	_else *JsonSchema
}

// Valid Json Schema types
const (
	OBJECT  = "object"
	ARRAY   = "array"
	STRING  = "string"
	NUMBER  = "number"
	INTEGER = "integer"
	BOOLEAN = "boolean"
	NULL    = "null"
)

// Valid values for "format" fields
const (
	DATE_TIME             = "date-time"
	TIME                  = "time"
	DATE                  = "date"
	EMAIL                 = "email"
	IDN_EMAIL             = "idn-email"
	HOSTNAME              = "hostname"
	IDN_HOSTNAME          = "idn-hostname"
	IPV4                  = "ipv4"
	IPV6                  = "ipv6"
	URI                   = "uri"
	URI_REFERENCE         = "uri-reference"
	IRI                   = "iri"
	IRI_REFERENCE         = "iri-reference"
	URI_TEMPLATE          = "uri-template"
	JSON_POINTER          = "json-pointer"
	RELATIVE_JSON_POINTER = "relative-json-pointer"
	REGEX                 = "regex"
)

// Valid values for "contentEncoding" field
var jsonSchemaContentEncodings = [...]string{
	"7bit",
	"8bit",
	"binary",
	"quoted-printable",
	"base64",
}
