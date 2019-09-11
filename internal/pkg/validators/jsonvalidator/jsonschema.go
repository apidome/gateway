package jsonvalidator

import (
	"fmt"
	"github.com/Creespye/caf/internal/pkg/jsonwalker"
	"log"
	"strings"
)

// Valid Json Schema types
const (
	TYPE_OBJECT  = "object"
	TYPE_ARRAY   = "array"
	TYPE_STRING  = "string"
	TYPE_NUMBER  = "number"
	TYPE_INTEGER = "integer"
	TYPE_BOOLEAN = "boolean"
	TYPE_NULL    = "null"
)

// Valid values for "format" fields
const (
	FORMAT_DATE_TIME             = "date-time"
	FORMAT_TIME                  = "time"
	FORMAT_DATE                  = "date"
	FORMAT_EMAIL                 = "email"
	FORMAT_IDN_EMAIL             = "idn-email"
	FORMAT_HOSTNAME              = "hostname"
	FORMAT_IDN_HOSTNAME          = "idn-hostname"
	FORMAT_IPV4                  = "ipv4"
	FORMAT_IPV6                  = "ipv6"
	FORMAT_URI                   = "uri"
	FORMAT_URI_REFERENCE         = "uri-reference"
	FORMAT_IRI                   = "iri"
	FORMAT_IRI_REFERENCE         = "iri-reference"
	FORMAT_URI_TEMPLATE          = "uri-template"
	FORMAT_JSON_POINTER          = "json-pointer"
	FORMAT_RELATIVE_JSON_POINTER = "relative-json-pointer"
	FORMAT_REGEX                 = "regex"
)

// Valid values for "contentEncoding" field
const (
	ENCODING_7BIT             = "7bit"
	ENCODING_8bit             = "8bit"
	ENCODING_BINARY           = "binary"
	ENCODING_QUOTED_PRINTABLE = "quited-printable"
	ENCODING_BASE64           = "base64"
)

type JsonSchema struct {
	// The $schema keyword is used to declare that a JSON fragment is
	// actually a piece of JSON Schema.
	Schema *schema `json:"$schema"`

	// The value of $ref is a URI, and the part after # sign is in a format
	// called JSON Pointer.
	Ref *ref `json:"$ref"`

	// The $id property is a URI that serves two purposes:
	// It declares a unique identifier for the schema
	// It declares a base URI against which $ref URIs are resolved.
	Id *id `json:"$id"`

	// The $comment keyword is strictly intended for adding comments
	// to the JSON schema source. Its value must always be a string.
	Comment *comment `json:"$comment"`

	// Title and Description used to describe the schema and not used for
	// validation.
	Title       *title       `json:"title"`
	Description *description `json:"description"`

	// The value of this keyword MUST be either a string or an array. If it is
	// an array, elements of the array MUST be strings and MUST be unique.
	// String values MUST be one of the six primitive types
	// ("null", "boolean", "object", "array", "number", or "string"),
	// or "integer" which matches any number with a zero fractional part.
	Type *_type `json:"type"`

	// The default keyword specifies a default value for an item.
	Default _default `json:"default"`

	// The examples keyword is a place to provide an array of examples
	// that validate against the schema.
	Examples examples `json:"examples"`

	// The value of this keyword MUST be an array.
	// An instance validates successfully against this keyword if its value is
	// equal to one of the elements in this keyword's array value.
	Enum enum `json:"enum"`

	// The value of this keyword MAY be of any type, including null.
	// An instance validates successfully against this keyword if its value is
	// equal to the value of the keyword.
	Const *_const `json:"const"`

	// The "definitions" keywords provides a standardized location for schema
	// authors to inline re-usable JSON Schemas into a more general schema. The
	// keyword does not directly affect the validation result.
	// This keyword's value MUST be an object. Each member value of this
	// object MUST be a valid JSON Schema.
	Definitions definitions `json:"definitions"`

	// The value of "properties" MUST be an object. Each value of this object
	// MUST be a valid JSON Schema.
	// This keyword determines how child instances validate for objects, and
	// does not directly validate the immediate instance itself.
	// Validation succeeds if, for each name that appears in both the instance
	// and as a name within this keyword's value, the child instance for that
	// name successfully validates against the corresponding schema.
	Properties properties `json:"properties"`

	// The value of "additionalProperties" MUST be a valid JSON Schema.
	// This keyword determines how child instances validate for objects,
	// and does not directly validate the immediate instance itself.
	// Validation with "additionalProperties" applies only to the child values
	// of instance names that do not match any names in "properties", and do
	// not match any regular expression in "patternProperties".
	// For all such properties, validation succeeds if the child instance
	// validates against the "additionalProperties" schema.
	AdditionalProperties *additionalProperties `json:"additionalProperties"`

	// The value of this keyword MUST be an array. Elements of this array,
	// if any, MUST be strings, and MUST be unique.
	// An object instance is valid against this keyword if every item in the
	// array is the name of a property in the instance.
	Required required `json:"required"`

	// The value of "propertyNames" MUST be a valid JSON Schema.
	// If the instance is an object, this keyword validates if every property
	// name in the instance validates against the provided schema. Note the
	// property name that the schema is testing will always be a string
	PropertyNames *propertyNames `json:"propertyNames"`

	// This keyword specifies rules that are evaluated if the instance is an
	// object and contains a certain property.
	// This keyword's value MUST be an object. Each property specifies a
	// dependency. Each dependency value MUST be an array or a valid JSON
	// Schema.
	// If the dependency value is a subschema, and the dependency key is a
	// property in the instance, the entire instance must validate against the
	// dependency value.
	// If the dependency value is an array, each element in the array, if any,
	// MUST be a string, and MUST be unique. If the dependency key is a
	// property in the instance, each of the items in the dependency value
	// must be a property that exists in the instance.
	Dependencies dependencies `json:"dependencies"`

	// The value of "patternProperties" MUST be an object. Each property name
	// of this object SHOULD be a valid regular expression, according to the
	// ECMA 262 regular expression dialect. Each property value of this object
	// MUST be a valid JSON Schema.
	// This keyword determines how child instances validate for objects, and
	// does not directly validate the immediate instance itself. Validation of
	// the primitive instance type against this keyword always succeeds.
	// Validation succeeds if, for each instance name that matches any regular
	// expressions that appear as a property name in this keyword's value, the
	// child instance for that name successfully validates against each schema
	// that corresponds to a matching regular expression.
	PatternProperties patternProperties `json:"patternProperties"`

	// The value of "items" MUST be either a valid JSON Schema or an array of
	// valid JSON Schemas.
	// This keyword determines how child instances validate for arrays, and
	// does not directly validate the immediate instance itself.
	// If "items" is a schema, validation succeeds if all elements in the array
	// successfully validate against that schema.
	// If "items" is an array of schemas, validation succeeds if each element
	// of the instance validates against the schema at the same position,
	// if any.
	Items items `json:"items"`

	// The value of this keyword MUST be a boolean.
	// If this keyword has boolean value false, the instance validates
	// successfully. If it has boolean value true, the instance validates
	// successfully if all of its elements are unique.
	Contains *contains `json:"contains"`

	// The value of "additionalItems" MUST be a valid JSON Schema.
	// This keyword determines how child instances validate for arrays, and
	// does not directly validate the immediate instance itself.
	// If "items" is an array of schemas, validation succeeds if every
	// instance element at a position greater than the size of "items"
	// validates against "additionalItems".
	// Otherwise, "additionalItems" MUST be ignored, as the "items" schema
	// (possibly the default value of an empty schema) is applied to all
	// elements.
	AdditionalItems *additionalItems `json:"additionalItems"`

	// array limitations
	MinItems    *minItems    `json:"minItems"`
	MaxItems    *maxItems    `json:"maxItems"`
	UniqueItems *uniqueItems `json:"uniqueItems"`

	// string limitations
	MinLength *minLength `json:"minLength,omitempty"`
	MaxLength *maxLength `json:"maxLength"`
	Pattern   *pattern   `json:"pattern"`
	Format    *format    `json:"format"`

	// integer/number limitations
	MultipleOf       *multipleOf       `json:"multipleOf"`
	Minimum          *minimum          `json:"minimum"`
	Maximum          *maximum          `json:"maximum"`
	ExclusiveMinimum *exclusiveMinimum `json:"exclusiveMinimum"`
	ExclusiveMaximum *exclusiveMaximum `json:"exclusiveMaximum"`

	// object size limitations
	MinProperties *minProperties `json:"minProperties"`
	MaxProperties *maxProperties `json:"maxProperties"`

	// The contentMediaType keyword specifies the MIME type of the contents
	// of a string.
	ContentMediaType *contentMediaType `json:"contentMediaType"`

	// The contentEncoding keyword specifies the encoding used to store
	// the contents.
	ContentEncoding *contentEncoding `json:"contentEncoding"`

	// Must be valid against any of the sub-schemas.
	AnyOf anyOf `json:"anyOf"`

	// Must be valid against all of the sub-schemas.
	AllOf allOf `json:"allOf"`

	// Must be valid against exactly one of the sub-schemas.
	OneOf oneOf `json:"oneOf"`

	// Must not be valid against the given schema.
	Not *not `json:"not"`

	// The if, then and else keywords allow the application of a sub-schema
	// based on the outcome of another schema.
	If   *_if   `json:"if"`
	Then *_then `json:"then"`
	Else *_else `json:"else"`

	// If "readOnly" has a value of boolean true, it indicates that the value
	// of the instance is managed exclusively by the owning authority, and
	// attempts by an application to modify the value of this property are
	// expected to be ignored or rejected by that owning authority.
	ReadOnly *readOnly `json:"readOnly"`

	// If "writeOnly" has a value of boolean true, it indicates that the value
	// is never present when the instance is retrieved from the owning
	// authority.
	// It can be present when sent to the owning authority to update or create
	// the document (or the resource it represents), but it will not be
	// included in any updated or newly created version of the instance.
	WriteOnly *writeOnly `json:"writeOnly"`
}

func (js *JsonSchema) validateJsonData(jsonPath string, jsonData []byte) (bool, error) {
	fmt.Println("[JsonSchema DEBUG] Validating " + jsonPath)

	// Calculate the relative path in order to evaluate the data
	jsonTokens := strings.Split(jsonPath, "/")
	relativeJsonPath := "/" + jsonTokens[len(jsonTokens)-1]

	// Create a new JsonPointer.
	jsonPointer, err := jsonwalker.NewJsonPointer(relativeJsonPath)
	if err != nil {
		fmt.Println("[JsonSchema DEBUG] validateJsonData() " +
			"failed while trying to create JsonPointer " + jsonPath)
		return false, err
	}

	// Get the piece of json that the current schema describes.
	value, err := jsonPointer.Evaluate(jsonData)
	if err != nil {
		fmt.Println("[JsonSchema DEBUG] validateJsonData() " +
			"failed while trying to evaluate a JsonPointer " + jsonPath)
		return false, err
	}

	// Get a slice of all of JsonSchema's field in order to iterate them
	// and call each of their validate() functions.
	keywordValidators := getKeywordsSlice(js)

	// Iterate over the keywords.
	for _, keyword := range keywordValidators {
		// TODO: Check if keyword != nil

		// Validate the value that we extracted from the jsonData at each
		// keyword.
		valid, err := keyword.validate(jsonPath, value)
		if err != nil {
			log.Print("[JsonSchema DEBUG] validation failed in path: " +
				jsonPath + " - " + err.Error())
			return valid, err
		}
	}

	return true, nil
}

func getKeywordsSlice(js *JsonSchema) []keywordValidator {
	return []keywordValidator{
		js.Type,
		js.Const,
		js.Enum,
		js.MinLength,
		js.MaxLength,
		js.Pattern,
		js.Format,
		js.MultipleOf,
		js.Minimum,
		js.Maximum,
		js.ExclusiveMinimum,
		js.ExclusiveMaximum,
		js.Required,
		js.Properties,
		js.AdditionalProperties,
		js.PropertyNames,
		js.Dependencies,
		js.PatternProperties,
		js.MinProperties,
		js.MaxProperties,
		js.Items,
		js.Contains,
		js.AdditionalItems,
		js.MinItems,
		js.MaxItems,
		js.UniqueItems,
		js.AnyOf,
		js.AllOf,
		js.OneOf,
		js.Not,
		js.If,
		js.Then,
		js.Else,
	}
}
