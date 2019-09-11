package jsonvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/Creespye/caf/internal/pkg/jsonwalker"
	"log"
	"strconv"
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
	Schema *schema `json:"$schema,omitempty"`

	// The value of $ref is a URI, and the part after # sign is in a format
	// called JSON Pointer.
	Ref *ref `json:"$ref,omitempty"`

	// The $id property is a URI that serves two purposes:
	// It declares a unique identifier for the schema
	// It declares a base URI against which $ref URIs are resolved.
	Id *id `json:"$id,omitempty"`

	// The $comment keyword is strictly intended for adding comments
	// to the JSON schema source. Its value must always be a string.
	Comment *comment `json:"$comment,omitempty"`

	// Title and Description used to describe the schema and not used for
	// validation.
	Title       *title       `json:"title,omitempty"`
	Description *description `json:"description,omitempty"`

	// The value of this keyword MUST be either a string or an array. If it is
	// an array, elements of the array MUST be strings and MUST be unique.
	// String values MUST be one of the six primitive types
	// ("null", "boolean", "object", "array", "number", or "string"),
	// or "integer" which matches any number with a zero fractional part.
	Type *_type `json:"type,omitempty"`

	// The default keyword specifies a default value for an item.
	Default _default `json:"default,omitempty"`

	// The examples keyword is a place to provide an array of examples
	// that validate against the schema.
	Examples examples `json:"examples,omitempty"`

	// The value of this keyword MUST be an array.
	// An instance validates successfully against this keyword if its value is
	// equal to one of the elements in this keyword's array value.
	Enum enum `json:"enum,omitempty"`

	// The value of this keyword MAY be of any type, including null.
	// An instance validates successfully against this keyword if its value is
	// equal to the value of the keyword.
	Const *_const `json:"const,omitempty"`

	// The "definitions" keywords provides a standardized location for schema
	// authors to inline re-usable JSON Schemas into a more general schema. The
	// keyword does not directly affect the validation result.
	// This keyword's value MUST be an object. Each member value of this
	// object MUST be a valid JSON Schema.
	Definitions definitions `json:"definitions,omitempty"`

	// The value of "properties" MUST be an object. Each value of this object
	// MUST be a valid JSON Schema.
	// This keyword determines how child instances validate for objects, and
	// does not directly validate the immediate instance itself.
	// Validation succeeds if, for each name that appears in both the instance
	// and as a name within this keyword's value, the child instance for that
	// name successfully validates against the corresponding schema.
	Properties properties `json:"properties,omitempty"`

	// The value of "additionalProperties" MUST be a valid JSON Schema.
	// This keyword determines how child instances validate for objects,
	// and does not directly validate the immediate instance itself.
	// Validation with "additionalProperties" applies only to the child values
	// of instance names that do not match any names in "properties", and do
	// not match any regular expression in "patternProperties".
	// For all such properties, validation succeeds if the child instance
	// validates against the "additionalProperties" schema.
	AdditionalProperties *additionalProperties `json:"additionalProperties,omitempty"`

	// The value of this keyword MUST be an array. Elements of this array,
	// if any, MUST be strings, and MUST be unique.
	// An object instance is valid against this keyword if every item in the
	// array is the name of a property in the instance.
	Required required `json:"required,omitempty"`

	// The value of "propertyNames" MUST be a valid JSON Schema.
	// If the instance is an object, this keyword validates if every property
	// name in the instance validates against the provided schema. Note the
	// property name that the schema is testing will always be a string
	PropertyNames *propertyNames `json:"propertyNames,omitempty"`

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
	Dependencies dependencies `json:"dependencies,omitempty"`

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
	PatternProperties patternProperties `json:"patternProperties,omitempty"`

	// The value of "items" MUST be either a valid JSON Schema or an array of
	// valid JSON Schemas.
	// This keyword determines how child instances validate for arrays, and
	// does not directly validate the immediate instance itself.
	// If "items" is a schema, validation succeeds if all elements in the array
	// successfully validate against that schema.
	// If "items" is an array of schemas, validation succeeds if each element
	// of the instance validates against the schema at the same position,
	// if any.
	Items items `json:"items,omitempty"`

	// The value of this keyword MUST be a boolean.
	// If this keyword has boolean value false, the instance validates
	// successfully. If it has boolean value true, the instance validates
	// successfully if all of its elements are unique.
	Contains *contains `json:"contains,omitempty"`

	// The value of "additionalItems" MUST be a valid JSON Schema.
	// This keyword determines how child instances validate for arrays, and
	// does not directly validate the immediate instance itself.
	// If "items" is an array of schemas, validation succeeds if every
	// instance element at a position greater than the size of "items"
	// validates against "additionalItems".
	// Otherwise, "additionalItems" MUST be ignored, as the "items" schema
	// (possibly the default value of an empty schema) is applied to all
	// elements.
	AdditionalItems *additionalItems `json:"additionalItems,omitempty"`

	// array limitations
	MinItems    *minItems    `json:"minItems,omitempty"`
	MaxItems    *maxItems    `json:"maxItems,omitempty"`
	UniqueItems *uniqueItems `json:"uniqueItems,omitempty"`

	// string limitations
	MinLength *minLength `json:"minLength,omitempty"`
	MaxLength *maxLength `json:"maxLength,omitempty"`
	Pattern   *pattern   `json:"pattern,omitempty"`
	Format    *format    `json:"format,omitempty"`

	// integer/number limitations
	MultipleOf       *multipleOf       `json:"multipleOf,omitempty"`
	Minimum          *minimum          `json:"minimum,omitempty"`
	Maximum          *maximum          `json:"maximum,omitempty"`
	ExclusiveMinimum *exclusiveMinimum `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *exclusiveMaximum `json:"exclusiveMaximum,omitempty"`

	// object size limitations
	MinProperties *minProperties `json:"minProperties,omitempty"`
	MaxProperties *maxProperties `json:"maxProperties,omitempty"`

	// The contentMediaType keyword specifies the MIME type of the contents
	// of a string.
	ContentMediaType *contentMediaType `json:"contentMediaType,omitempty"`

	// The contentEncoding keyword specifies the encoding used to store
	// the contents.
	ContentEncoding *contentEncoding `json:"contentEncoding,omitempty"`

	// Must be valid against any of the sub-schemas.
	AnyOf anyOf `json:"anyOf,omitempty"`

	// Must be valid against all of the sub-schemas.
	AllOf allOf `json:"allOf,omitempty"`

	// Must be valid against exactly one of the sub-schemas.
	OneOf oneOf `json:"oneOf,omitempty"`

	// Must not be valid against the given schema.
	Not *not `json:"not,omitempty"`

	// The if, then and else keywords allow the application of a sub-schema
	// based on the outcome of another schema.
	If   *_if   `json:"if,omitempty"`
	Then *_then `json:"then,omitempty"`
	Else *_else `json:"else,omitempty"`

	// If "readOnly" has a value of boolean true, it indicates that the value
	// of the instance is managed exclusively by the owning authority, and
	// attempts by an application to modify the value of this property are
	// expected to be ignored or rejected by that owning authority.
	ReadOnly *readOnly `json:"readOnly,omitempty"`

	// If "writeOnly" has a value of boolean true, it indicates that the value
	// is never present when the instance is retrieved from the owning
	// authority.
	// It can be present when sent to the owning authority to update or create
	// the document (or the resource it represents), but it will not be
	// included in any updated or newly created version of the instance.
	WriteOnly *writeOnly `json:"writeOnly,omitempty"`
}

func NewJsonSchema(bytes []byte) (*JsonSchema, error) {
	var schema *JsonSchema

	// Check if the string s is a valid json.
	err := json.Unmarshal(bytes, &schema)
	if err != nil {
		return nil, err
	}

	err = schema.connectRelatedKeywords("")
	if err != nil {
		fmt.Println("[JsonSchema DEBUG] connectRelatedKeywords() " +
			"failed: " + err.Error())
		return nil, err
	}

	return schema, nil
}

func (js *JsonSchema) connectRelatedKeywords(schemaPath string) error {
	// Connect sub-schemas in "properties" field.
	for key := range js.Properties {
		err := js.Properties[key].connectRelatedKeywords(schemaPath + "/properties/" + key)
		if err != nil {
			return err
		}
	}

	// Connect sub-schema in "additionalProperties" field.
	if js.AdditionalProperties != nil {
		err := js.AdditionalProperties.connectRelatedKeywords(schemaPath + "/additionalProperties")
		if err != nil {
			return err
		}

		// If "properties" field exists in the schema, save the keywordValidator's
		// address in "AdditionalProperties".
		if js.Properties != nil {
			js.AdditionalProperties.siblingProperties = &js.Properties
		}

		// If "patternProperties" field exists in the schema, save the keywordValidator's
		// address in "AdditionalProperties".
		if js.PatternProperties != nil {
			js.AdditionalProperties.siblingPatternProperties = &js.PatternProperties
		}
	}

	// Connect sub-schema in "propertyNames" field.
	if js.PropertyNames != nil {
		err := js.PropertyNames.connectRelatedKeywords(schemaPath + "/propertyNames")
		if err != nil {
			return err
		}
	}

	for key, value := range js.Dependencies {
		if v, ok := value.(map[string]interface{}); ok {
			var subSchema JsonSchema

			// Marshal the dependency in order to Unmarshal it into JsonSchema struct.
			rawDependency, err := json.Marshal(v)
			if err != nil {
				return SchemaCompilationError{
					schemaPath,
					err.Error(),
				}
			}

			// Unmarshal the raw data in order into a JsonSchema struct.
			err = json.Unmarshal(rawDependency, &subSchema)
			if err != nil {
				return SchemaCompilationError{
					schemaPath,
					err.Error(),
				}
			}

			err = subSchema.connectRelatedKeywords(schemaPath + "/dependencies/" + key)
			if err != nil {
				return err
			}

			js.Dependencies[key] = subSchema
		}
	}

	// Connect sub-schemas in "patternProperties" field.
	for key := range js.PatternProperties {
		err := js.PatternProperties[key].connectRelatedKeywords(schemaPath + "/patternProperties/" + key)
		if err != nil {
			return err
		}
	}

	// Connect sub-schemas in "definitions" field.
	for key := range js.Definitions {
		err := js.Definitions[key].connectRelatedKeywords(schemaPath + "/definitions/" + key)
		if err != nil {
			return err
		}
	}

	if js.Items != nil {
		var items interface{}
		err := json.Unmarshal(js.Items, &items)
		if err != nil {
			return SchemaCompilationError{
				schemaPath,
				err.Error(),
			}
		}

		switch v := items.(type) {
		case map[string]interface{}:
			{
				// Marshal the dependency in order to Unmarshal it into JsonSchema struct.
				rawSubSchema, err := json.Marshal(v)
				if err != nil {
					return SchemaCompilationError{
						schemaPath,
						err.Error(),
					}
				}

				// Create a new JsonSchema object.
				subSchema, err := NewJsonSchema(rawSubSchema)
				if err != nil {
					return err
				}

				err = subSchema.connectRelatedKeywords(schemaPath + "/items")
				if err != nil {
					return err
				}

				js.Items, err = json.Marshal(subSchema)
				if err != nil {
					return SchemaCompilationError{
						schemaPath,
						err.Error(),
					}
				}
			}
		case []interface{}:
			{
				for index, value := range v {
					// Marshal the dependency in order to Unmarshal it into JsonSchema struct.
					rawSubSchema, err := json.Marshal(value)
					if err != nil {
						return SchemaCompilationError{
							schemaPath,
							err.Error(),
						}
					}

					// Create a new JsonSchema object.
					subSchema, err := NewJsonSchema(rawSubSchema)
					if err != nil {
						return err
					}

					err = subSchema.connectRelatedKeywords(schemaPath + "/items/" + strconv.Itoa(index))
					if err != nil {
						return err
					}

					v[index] = subSchema
				}

				js.Items, err = json.Marshal(v)
				if err != nil {
					return SchemaCompilationError{
						schemaPath,
						err.Error(),
					}
				}
			}
		}
	}

	// Connect sub-schema in "additionalItems" field.
	if js.AdditionalItems != nil {
		err := js.AdditionalItems.connectRelatedKeywords(schemaPath + "/additionalItems")
		if err != nil {
			return err
		}

		// If "items" field exists in the schema, save the keywordValidator's
		// address in "AdditionalItems".
		if js.Items != nil {
			js.AdditionalItems.siblingItems = &js.Items
		}
	}

	// Connect sub-schema in "contains" field.
	if js.Contains != nil {
		err := js.Contains.connectRelatedKeywords(schemaPath + "/contains")
		if err != nil {
			return err
		}
	}

	// Connect sub-schemas in "anyOf" field.
	for index := range js.AnyOf {
		err := js.AnyOf[index].connectRelatedKeywords(schemaPath + "/anyOf/" + strconv.Itoa(index))
		if err != nil {
			return err
		}
	}

	// Connect sub-schemas in "allOf" field.
	for index := range js.AllOf {
		err := js.AllOf[index].connectRelatedKeywords(schemaPath + "/allOf/" + strconv.Itoa(index))
		if err != nil {
			return err
		}
	}

	// Connect sub-schemas in "oneOf" field.
	for index := range js.OneOf {
		err := js.OneOf[index].connectRelatedKeywords(schemaPath + "/oneOf/" + strconv.Itoa(index))
		if err != nil {
			return err
		}
	}

	// Connect sub-schema in "not" field.
	if js.Not != nil {
		err := js.Not.connectRelatedKeywords(schemaPath + "/not")
		if err != nil {
			return err
		}
	}

	// Connect sub-schema in "if" field.
	if js.If != nil {
		err := js.If.connectRelatedKeywords(schemaPath + "/if")
		if err != nil {
			return err
		}

		// Connect sub-schema in "then" field.
		if js.Then != nil {
			err := js.Then.connectRelatedKeywords(schemaPath + "/then")
			if err != nil {
				return err
			}

			// If "then" field exists in the schema, save the keywordValidator's
			// address in "If".
			js.If.siblingThen = js.Then
		}

		// Connect sub-schema in "else" field.
		if js.Else != nil {
			err := js.Else.connectRelatedKeywords(schemaPath + "/else")
			if err != nil {
				return err
			}

			// If "else" field exists in the schema, save the keywordValidator's
			// address in "If".
			js.If.siblingElse = js.Else
		}
	}

	return nil
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
		js.PropertyNames,
		js.Properties,
		js.AdditionalProperties,
		js.PatternProperties,
		js.Dependencies,
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
