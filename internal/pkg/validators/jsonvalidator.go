package validators

import "encoding/json"

type JsonSchema map[string]interface{}

// Valid keywords in Json Schema (draft 7)
var jsonSchemaKeywords = [...]string{

	/**********************/
	/** GENERIC KEYWORDS **/
	/**********************/

	// The $schema keyword is used to declare that a JSON fragment is
	// actually a piece of JSON Schema.
	"$schema",

	// The value of $ref is a URI, and the part after # sign is in a format
	// called JSON Pointer.
	"$ref",

	// The $id property is a URI that serves two purposes:
	// It declares a unique identifier for the schema
	// It declares a base URI against which $ref URIs are resolved.
	"$id",

	// The $comment keyword is strictly intended for adding comments
	// to the JSON schema source. Its value must always be a string.
	"$comment",

	// Title and Description used to describe the schema and not used for
	// validation.
	"title",
	"description",

	// The default keyword specifies a default value for an item.
	"default",

	// The examples keyword is a place to provide an array of examples
	// that validate against the schema.
	"examples",

	// The enum keyword is used to restrict a value to a fixed set of values.
	// It must be an array with at least one element, where each element
	// is unique.
	"enum",

	// The const keyword is used to restrict a value to a single value.
	"const",

	// The definitions keyword is used to create entities that we recognize as
	// repetitive entities.
	// This ability maintains reuse in out Json Schema.
	"definitions",

	/************************/
	/** OBJECT DESCRIPTORS **/
	/************************/

	// The value of properties is an object, where each key is the name of a
	// property and each value is a JSON schema used to validate that property.
	"properties",

	// The additionalProperties keyword is used to control the handling of
	// extra stuff, that is, properties whose names are not listed in the
	// properties keyword.
	// By default any additional properties are allowed.
	// The additionalProperties keyword may be either a boolean or an object.
	// If additionalProperties is a boolean and set to false, no additional
	// properties will be allowed.
	// If additionalProperties is an object, that object is a schema that will be
	// used to validate any additional properties not listed in properties.
	"additionalProperties",

	// The required keyword takes an array of zero or more strings.
	// Each of these strings must be unique.
	"required",

	// The names of properties can be validated against a schema, irrespective
	// of their values.
	// This can be useful if you don’t want to enforce specific properties,
	// but you want to make sure that the names of those properties follow
	// a specific convention.
	"propertyNames",

	// The dependencies keyword allows the schema of the object to change
	// based on the presence of certain special properties.
	"dependencies",

	// TODO: Learn more about this keyword.
	"patternProperties",

	/***********************/
	/** ARRAY DESCRIPTORS **/
	/***********************/

	// Items can be either an object or an array. If it is an object, it will
	// represent a schema that all the items in the array should match.
	// If it is an array, each item in that array is a different json schema
	// that should match the corresponding item in the inspected array
	// (In this case the index of each item is very important).
	"items",

	// While the items schema must be valid for every item in the array,
	// the contains schema only needs to validate against one or more
	// items in the array.
	"contains",

	// The additionalItems keyword controls whether it’s valid to have
	// additional items in the array beyond what is defined in items.
	"additionalItems",

	/*****************/
	/** LIMITATIONS **/
	/*****************/

	// array limitations
	"minItems",
	"maxItems",
	"uniqueItems",

	// string limitations
	"minLength",
	"maxLength",
	"pattern",
	"format",

	// integer/number limitations
	"multipleOf",
	"minimum",
	"maximum",
	"exclusiveMinimum",
	"exclusiveMaximum",

	// object size limitations
	"minProperties",
	"maxProperties",

	/*****************/
	/** MEDIA TYPES **/
	/*****************/

	// The contentMediaType keyword specifies the MIME type of the contents
	// of a string.
	"contentMediaType",

	// The contentEncoding keyword specifies the encoding used to store
	// the contents.
	"contentEncoding",

	/*************************/
	/** SCHEMAS COMBINATION **/
	/*************************/

	// Must be valid against any of the sub-schemas.
	"anyOf",

	// Must be valid against all of the sub-schemas.
	"allOf",

	// Must be valid against exactly one of the sub-schemas.
	"oneOf",

	// Must not be valid against the given schema
	"not",

	// The if, then and else keywords allow the application of a sub-schema
	// based on the outcome of another schema.
	"if",
	"then",
	"else",
}

// Valid values for "type" field
var jsonSchemaTypes = [...]string{
	"object",
	"array",
	"string",
	"number",
	"integer",
	"boolean",
	"null",
}

// Valid values for "format" field
var jsonSchemaBuiltInFormats = [...]string{
	"date-time",
	"time",
	"date",
	"email",
	"idn-email",
	"hostname",
	"idn-hostname",
	"ipv4",
	"ipv6",
	"uri",
	"uri-reference",
	"iri",
	"iri-reference",
	"uri-template",
	"json-pointer",
	"relative-json-pointer",
	"regex",
}

// Valid values for "contentEncoding" field
var jsonSchemaContentEncodings = [...]string{
	"7bit",
	"8bit",
	"binary",
	"quoted-printable",
	"base64",
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
