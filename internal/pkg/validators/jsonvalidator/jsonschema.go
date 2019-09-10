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

	// Type specifies the acceptable value type of the schema.
	Type *_type `json:"type"`

	// The default keyword specifies a default value for an item.
	Default _default `json:"default"`

	// The examples keyword is a place to provide an array of examples
	// that validate against the schema.
	Examples examples `json:"examples"`

	// The enum keyword is used to restrict a value to a fixed set of values.
	// It must be an array with at least one element, where each element
	// is unique.
	Enum enum `json:"enum"`

	// The const keyword is used to restrict a value to a single value.
	Const _const `json:"const"`

	// The definitions keyword is used to create entities that we recognize as
	// repetitive entities.
	// This ability maintains reuse in out Json Schema.
	Definitions definitions `json:"definitions"`

	// The value of properties is an object, where each key is the name of a
	// property and each value is a JSON schema used to validate that property.
	Properties properties `json:"properties"`

	// The additionalProperties keyword is used to control the handling of
	// extra stuff, that is, properties whose names are not listed in the
	// properties keyword.
	// By default any additional properties are allowed.
	// The additionalProperties keyword may be either a boolean or an object.
	// If additionalProperties is a boolean and set to false, no additional
	// properties will be allowed.
	// If additionalProperties is an object, that object is a schema that will be
	// used to validate any additional properties not listed in properties.
	AdditionalProperties *additionalProperties `json:"additionalProperties"`

	// The required keyword takes an array of zero or more strings.
	// Each of these strings must be unique.
	Required required `json:"required"`

	// The names of properties can be validated against a schema, irrespective
	// of their values.
	// This can be useful if you don’t want to enforce specific properties,
	// but you want to make sure that the names of those properties follow
	// a specific convention.
	PropertyNames *propertyNames `json:"propertyNames"`

	// The dependencies keyword allows the schema of the object to change
	// based on the presence of certain special properties.
	Dependencies dependencies `json:"dependencies"`

	// TODO: Learn more about this keyword.
	PatternProperties patternProperties `json:"patternProperties"`

	// Items can be either an object or an array. If it is an object, it will
	// represent a schema that all the items in the array should match.
	// If it is an array, each item in that array is a different json schema
	// that should match the corresponding item in the inspected array
	// (In this case the index of each item is very important).
	Items items `json:"items"`

	// While the items schema must be valid for every item in the array,
	// the contains schema only needs to validate against one or more
	// items in the array.
	Contains *contains `json:"contains"`

	// The additionalItems keyword controls whether it’s valid to have
	// additional items in the array beyond what is defined in items.
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
		fmt.Println("[JsonSchema DEBUG] validateJsonData() failed while trying to create JsonPointer " + jsonPath)
		return false, err
	}

	// Get the piece of json that the current schema describes.
	value, err := jsonPointer.Evaluate(jsonData)
	if err != nil {
		fmt.Println("[JsonSchema DEBUG] validateJsonData() failed while trying to evaluate a JsonPointer " + jsonPath)
		return false, nil
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
			log.Print("[JsonSchema DEBUG] validation failed in path: " + jsonPath + " - " + err.Error())
			return valid, err
		}
	}

	return true, nil
}

func getKeywordsSlice(js *JsonSchema) []keywordValidator {
	return []keywordValidator{
		js.Type,
		js.Enum,
		js.Const,
		js.MinLength,
		js.MaxLength,
		js.Pattern,
		js.Format,
		js.MultipleOf,
		js.Minimum,
		js.Maximum,
		js.ExclusiveMinimum,
		js.ExclusiveMaximum,
		js.Properties,
		js.AdditionalProperties,
		js.Required,
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
