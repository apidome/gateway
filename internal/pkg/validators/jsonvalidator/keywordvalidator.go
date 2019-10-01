package jsonvalidator

import (
	"encoding/json"
	"errors"
	"github.com/Creespye/caf/internal/pkg/validators/formatchecker"
	"math"
	"regexp"
	"strconv"
	"strings"
)

/*
Implemented keywordValidators:
> enum: 					V
> _const: 					V
> _type: 					V ***
> minLength: 				V
> maxLength: 				V
> pattern: 					V
> format: 					X
> multipleOf: 				V
> minimum: 					V
> maximum: 					V
> exclusiveMinimum: 		V
> exclusiveMaximum: 		V
> properties: 				V
> additionalProperties: 	V
> required: 				V
> propertyNames: 			V
> dependencies: 			V
> patternProperties: 		V
> minProperties: 			V
> maxProperties: 			V
> items: 					V ***
> contains: 				V
> additionalItems: 			V
> minItems: 				V
> maxItems: 				V
> uniqueItems: 				V
> anyOf: 					V
> allOf: 					V
> oneOf: 					V
> not: 						V
> _if: 						V
> _then: 					V
> _else: 					V

*** These keywords are being un-marshaled in their validate() function.
	We need to find a way to do that on startup and not on runtime.

*/

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

type keywordValidator interface {
	validate(string, jsonData, string) (bool, error)
}

/*****************/
/** Annotations **/
/*****************/

type ref string

func (r ref) validateByRef(jsonPath string, jsonData []byte, rootSchemaID string) (bool, error) {
	schemaURI := strings.Split(string(r), "#")[0]
	fragment := strings.Split(string(r), "#")[1]

	if schemaURI == "" {
		schemaURI = rootSchemaID
	}

	if rootSchema, ok := rootSchemaPool[schemaURI]; ok {
		if fragment != "" {
			if subSchema, ok := rootSchema.subSchemaMap[fragment]; ok {
				return subSchema.validateJsonData(jsonPath, jsonData, rootSchemaID)
			} else {
				return false, errors.New("TODO: implement new error type - sub-schema " + fragment +
					"in root schema " + schemaURI + " does not exist")
			}
		} else {
			return rootSchema.validateJsonData(jsonPath, jsonData, rootSchemaID)
		}
	} else {
		return false, errors.New("TODO: implement new error type - root-schema " + schemaURI + "does not exist")
	}
}

type schema string
type id string
type comment string
type title string
type description string
type examples []interface{}
type definitions map[string]*JsonSchema
type _default json.RawMessage

func (d *_default) UnmarshalJSON(data []byte) error {
	*d = data
	return nil
}

/**********************/
/** Generic Keywords **/
/**********************/

type _type json.RawMessage

func (t *_type) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if t == nil {
		return true, nil
	}

	var data interface{}

	// First we need to unmarshal the json data.
	err := json.Unmarshal(*t, &data)
	if err != nil {
		return false, err
	}

	// The "type" field in json schema can be represented by two different values:
	// - string - the inspected value can be only one json type.
	// - array - the inspected value can be a variety of json types.
	// - default - the schema is incorrect.
	switch typeFromSchema := data.(type) {
	case []interface{}:
		{
			// If we arrived this loop, it means "type" is an array of types.
			// We need to go over the existing types and perform
			// "json type assertion" of jsonData and the current json type.
			for _, typeFromList := range typeFromSchema {
				// A json type must be represented by a string.
				if v, ok := typeFromList.(string); ok {
					// Perform the "json type assertion"
					ok, _ := assertJsonType(v, jsonData.value)

					// If the assertion succeeded, return true
					if ok {
						return ok, nil
					}
				} else {
					return false, KeywordValidationError{
						"type",
						"\"type\" field in schema must be string or array of strings",
					}
				}
			}

			// JsonTypeMismatchError
			return false, KeywordValidationError{
				"type",
				"inspected value does not match any of the valid types in the schema",
			}
		}
	case string:
		{
			// In this case, there is only one valid type, so we
			// perform "json type assertion" of the json type and jsonData.
			return assertJsonType(typeFromSchema, jsonData.value)
		}
	default:
		{
			return false, KeywordValidationError{
				"type",
				"\"type\" field in schema must be string or array of strings",
			}
		}
	}
}

// assertJsonType is a function that gets a jsonType and some jsonData and
// returns true if the value belongs to the type.
// If it is not, the function will return an appropriate error.
func assertJsonType(jsonType string, jsonData interface{}) (bool, error) {
	switch jsonType {
	case TYPE_OBJECT:
		{
			if _, ok := jsonData.(map[string]interface{}); ok {
				return true, nil
			} else {
				return false, KeywordValidationError{
					"type",
					"inspected value expected to be a json object",
				}
			}
		}
	case TYPE_ARRAY:
		{
			if _, ok := jsonData.([]interface{}); ok {
				return true, nil
			} else {
				return false, KeywordValidationError{
					"type",
					"inspected value expected to be a json array",
				}
			}
		}
	case TYPE_STRING:
		{
			if _, ok := jsonData.(string); ok {
				return true, nil
			} else {
				return false, KeywordValidationError{
					"type",
					"inspected value expected to be a json string",
				}
			}
		}
	case TYPE_NUMBER, TYPE_INTEGER:
		{
			if _, ok := jsonData.(float64); ok {
				return true, nil
			} else {
				return false, KeywordValidationError{
					"type",
					"inspected value expected to be a json number",
				}
			}
		}
	case TYPE_BOOLEAN:
		{
			if _, ok := jsonData.(bool); ok {
				return true, nil
			} else {
				return false, KeywordValidationError{
					"type",
					"inspected value expected to be a json boolean",
				}
			}
		}
	//case TYPE_NULL:
	//	{
	//		if v, ok := jsonData.value.(string); ok {
	//			if v == "null" {
	//				return true, nil
	//			}
	//		} else {
	//			return false, KeywordValidationError{
	//				"type",
	//				"inspected value expected to be a json null",
	//			}
	//		}
	//	}
	default:
		{
			return false, KeywordValidationError{
				"type",
				"invalid json type " + jsonType,
			}
		}
	}
}

func (t *_type) UnmarshalJSON(data []byte) error {
	*t = data
	return nil
}

func (t *_type) MarshalJSON() ([]byte, error) {
	return []byte(*t), nil
}

type enum []interface{}

func (e enum) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if e == nil {
		return true, nil
	}

	// Iterate over the items in "enum" array.
	for _, item := range e {
		// Marshal the item from "enum" array back comparable value that does
		// not require type assertion.
		rawEnumItem, err := json.Marshal(item)
		if err != nil {
			return false, nil
		}

		// Convert both of the byte arrays to string for more convenient
		// comparison. If they are equal, the data is valid against "enum".
		if string(rawEnumItem) == string(jsonData.raw) {
			return true, nil
		}
	}

	// If we arrived here it means that the inspected value is not equal
	// to any of the values in "enum".
	return false, KeywordValidationError{
		"enum",
		"inspected value does not match any of the items in \"enum\" array",
	}
}

type _const json.RawMessage

func (c *_const) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if c == nil {
		return true, nil
	}

	// Convert both of the byte arrays to string for more convenient
	// comparison. If they are equal, the data is valid against "const".
	if string(*c) == string(jsonData.raw) {
		return true, nil
	} else {
		return false, KeywordValidationError{
			"const",
			"inspected value not equal to \"" + string(*c) + "\"",
		}
	}
}

func (c *_const) UnmarshalJSON(data []byte) error {
	// In this function we Unmarshal and then Marshal again
	// the argument data in order to remove special characters
	// like \n \t \r etc.

	var unmarshaledData interface{}

	err := json.Unmarshal(data, &unmarshaledData)
	if err != nil {
		return err
	}

	rawConst, err := json.Marshal(unmarshaledData)
	if err != nil {
		return err
	}

	*c = rawConst
	return nil
}

/*********************/
/** String Keywords **/
/*********************/

type minLength int

func (ml *minLength) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if ml == nil {
		return true, nil
	}

	// If jsonData is a string, validate its length,
	// else, return a KeywordValidationError
	if v, ok := jsonData.value.(string); ok {
		if len(v) >= int(*ml) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"minLength",
				"inspected string is less than " + strconv.Itoa(int(*ml)),
			}
		}
	} else {
		return false, KeywordValidationError{
			"minLength",
			"inspected value is not a string",
		}
	}
}

type maxLength int

func (ml *maxLength) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if ml == nil {
		return true, nil
	}

	// If jsonData is a string, validate its length,
	// else, return a KeywordValidationError
	if v, ok := jsonData.value.(string); ok {
		if len(v) <= int(*ml) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"maxLength",
				"inspected string is greater than " + strconv.Itoa(int(*ml)),
			}
		}
	} else {
		return false, KeywordValidationError{
			"maxLength",
			"inspected value is not a string",
		}
	}
}

type pattern string

func (p *pattern) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if p == nil {
		return true, nil
	}

	// If jsonData is a string, validate its length,
	// else, return a KeywordValidationError
	if v, ok := jsonData.value.(string); ok {
		match, err := regexp.MatchString(string(*p), v)

		// The pattern or the value is not in the right format (string)
		if err != nil {
			return false, KeywordValidationError{
				"pattern",
				err.Error(),
			}
		}

		if match {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"pattern",
				"value " + v + " does not match to pattern" + string(*p),
			}
		}
	} else {
		return false, KeywordValidationError{
			"pattern",
			"inspected value is not a string",
		}
	}
}

type format string

func (f *format) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if f == nil {
		return true, nil
	}

	if v, ok := jsonData.value.(string); ok {
		switch v {
		case FORMAT_DATE_TIME:
			if _, err := formatchecker.IsValidDateTime(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"date-time incorrectly formatted " + err.Error(),
				}
			}
		case FORMAT_DATE:
			if _, err := formatchecker.IsValidDate(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"date incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_TIME:
			if _, err := formatchecker.IsValidTime(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"time incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_EMAIL:
			if _, err := formatchecker.IsValidEmail(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"email incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_IDN_EMAIL:
			if _, err := formatchecker.IsValidIdnEmail(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"idn-email incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_HOSTNAME:
			if _, err := formatchecker.IsValidHostname(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"hostname incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_IDN_HOSTNAME:
			if _, err := formatchecker.IsValidIdnHostname(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"idn-hostname incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_IPV4:
			if _, err := formatchecker.IsValidIPv4(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"ipv4 incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_IPV6:
			if _, err := formatchecker.IsValidIPv6(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"ipv6 incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_URI:
			if _, err := formatchecker.IsValidURI(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"uri incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_URI_REFERENCE:
			if _, err := formatchecker.IsValidUriRef(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"uri-reference incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_IRI:
			if _, err := formatchecker.IsValidIri(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"iri incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_IRI_REFERENCE:
			if _, err := formatchecker.IsValidIriRef(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"iri-reference incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_URI_TEMPLATE:
			if _, err := formatchecker.IsValidURITemplate(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"uri-template incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_JSON_POINTER:
			if _, err := formatchecker.IsValidJSONPointer(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"json-pointer incorrectly formatted: " + err.Error(),
				}
			}
		case FORMAT_RELATIVE_JSON_POINTER:
			if _, err := formatchecker.IsValidRelJSONPointer(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"relative-json-pointer incorrectly formatted: " + err.Error(),
				}
			}
		case "regex":
			if _, err := formatchecker.IsValidRegex(v); err != nil {
				return false, KeywordValidationError{
					"format",
					"regex incorrectly formatted: " + err.Error(),
				}
			}
		default:
			return true, nil
		}
	} else {
		return false, KeywordValidationError{
			"pattern",
			"inspected value is not a string",
		}
	}
	return true, nil
}

/*********************/
/** Number Keywords **/
/*********************/

type multipleOf float64

func (mo *multipleOf) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if mo == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.value.(float64); ok {
		if math.Mod(v, float64(*mo)) == 0 {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"multipleOf",
				"inspected value is not a multiple of " + strconv.FormatFloat(float64(*mo),
					'f',
					6,
					64),
			}
		}
	} else {
		return false, KeywordValidationError{
			"multipleOf",
			"inspected value is not an integer",
		}
	}
}

type minimum float64

func (m *minimum) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if m == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.value.(float64); ok {
		if v >= float64(*m) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"minimum",
				"inspected value is less than " + strconv.FormatFloat(float64(*m),
					'f',
					6,
					64),
			}
		}
	} else {
		return false, KeywordValidationError{
			"minimum",
			"inspected value is not a number",
		}
	}
}

type maximum float64

func (m *maximum) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if m == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.value.(float64); ok {
		if v <= float64(*m) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"maximum",
				"inspected value is greater than " + strconv.FormatFloat(float64(*m),
					'f',
					6,
					64),
			}
		}
	} else {
		return false, KeywordValidationError{
			"maximum",
			"inspected value is not a number",
		}
	}
}

type exclusiveMinimum float64

func (em *exclusiveMinimum) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if em == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.value.(float64); ok {
		if v > float64(*em) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"exclusiveMinimum",
				"inspected value is not greater than " + strconv.FormatFloat(float64(*em),
					'f',
					6,
					64),
			}
		}
	} else {
		return false, KeywordValidationError{
			"exclusiveMinimum",
			"inspected value is not a number",
		}
	}
}

type exclusiveMaximum float64

func (em *exclusiveMaximum) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if em == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.value.(float64); ok {
		if v < float64(*em) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"exclusiveMaximum",
				"inspected value is not less than " + strconv.FormatFloat(float64(*em),
					'f',
					6,
					64),
			}
		}
	} else {
		return false, KeywordValidationError{
			"exclusiveMaximum",
			"inspected value is not a number",
		}
	}
}

/*********************/
/** Object Keywords **/
/*********************/

type properties map[string]*JsonSchema

func (p properties) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if p == nil {
		return true, nil
	}

	// First, we need to verify that jsonData is a json object
	if object, ok := jsonData.value.(map[string]interface{}); ok {
		// For each "property" validate it according to its JsonSchema.
		for key, value := range p {
			// Before we try to validate the data against the schema,
			// we make sure that the data actually contains the property.
			if _, ok := object[key]; ok {
				valid, err := value.validateJsonData(jsonPath+"/"+key, jsonData.raw, rootSchemaId)
				if err != nil {
					return valid, err
				}
			}
		}

		// If we arrived here, the validation of all the properties
		// succeeded.
		return true, nil
	} else {
		return false, KeywordValidationError{
			"properties",
			"inspected value expected to be a json object",
		}
	}
}

type additionalProperties struct {
	JsonSchema
	siblingProperties        *properties
	siblingPatternProperties *patternProperties
}

func (ap *additionalProperties) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if ap == nil {
		return true, nil
	}

	// First we need to verify that jsonData is a json object.
	if object, isObject := jsonData.value.(map[string]interface{}); isObject {
		// Iterate over the properties of the inspected object.
		for property := range object {
			if (*ap).siblingProperties != nil {
				// Check if the property does not have corresponding schema in
				// "properties" field
				if _, ok := (*ap.siblingProperties)[property]; !ok {
					// Iterate over the patterns in "patternProperties" field.
					for pattern := range *ap.siblingPatternProperties {
						// Check if the inspected property matches to the pattern.
						match, err := regexp.MatchString(pattern, property)

						// The pattern or the value is not in the right format (string)
						if err != nil {
							return false, KeywordValidationError{
								"additionalProperties",
								err.Error(),
							}
						}

						// If there is no match, validate the value of the property against
						// the given schema in "additionalProperties" field.
						if !match {
							valid, err := (*ap).validateJsonData(jsonPath+"/"+property, jsonData.raw, rootSchemaId)

							// If the validation fails, return an error.
							if !valid {
								return false, KeywordValidationError{
									"additionalProperties",
									"property \"" +
										property +
										"\" failed in validation: \n" + err.Error(),
								}
							}
						}
					}
				}
			}
		}

		// If we arrived here, none of the properties failed in validation,
		// and we return true.
		return true, nil
	} else {
		return false, KeywordValidationError{
			"properties",
			"inspected value expected to be a json object",
		}
	}
}

type required []string

func (r required) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if r == nil {
		return true, nil
	}

	// First, we must verify that jsonData is a json object.
	if v, ok := jsonData.value.(map[string]interface{}); ok {
		// For each property in the required list, check if it exists.
		for _, property := range r {
			if v[property] == nil {
				return false, KeywordValidationError{
					"required",
					"Missing required property - " + property,
				}
			}
		}
	} else {
		return false, KeywordValidationError{
			"required",
			"all items \"required\" field must be strings",
		}
	}

	// Is we arrived here, all the properties exist.
	return true, nil
}

type propertyNames struct {
	JsonSchema
}

func (pn *propertyNames) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if pn == nil {
		return true, nil
	}

	// First, we need to verify that jsonData is a json object
	if object, ok := jsonData.value.(map[string]interface{}); ok {
		// Iterate over the object's properties.
		for property := range object {
			// Validate the property name against the schema stored in "propertyNames" field
			valid, err := pn.validateJsonData("", []byte("\""+property+"\""), rootSchemaId)

			// If the property name could be validated against the scheme return an error
			if !valid {
				return false, KeywordValidationError{
					"propertyNames",
					"property name \"" + property + "\" failed in validation: " + err.Error(),
				}
			}
		}

		// If we arrived here it means that all the property names validated successfully against
		// the schema stored in "propertyNames".
		return true, nil
	} else {
		return false, KeywordValidationError{
			"propertyNames",
			"inspected value expected to be a json object",
		}
	}
}

type dependencies map[string]interface{}

func (d dependencies) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if d == nil {
		return true, nil
	}

	// First we need to verify that jsonData is a json object.
	if object, ok := jsonData.value.(map[string]interface{}); ok {

		// Iterate over the dependencies object from the schema.
		for propertyName, dependency := range d {
			// A dependency may be a json array (consist of strings) of a json
			// object which is a json schema that the inspected value need to
			// validated against.
			switch v := dependency.(type) {

			// In this case the dependency is a sub-schema.
			case *JsonSchema:
				{
					// Check if the propertyName (which is the key in the "dependencies" object)
					// is present in the data. If it is, validate the whole instance against the
					// sub-schema.
					if _, ok := object[propertyName]; ok {
						// Validate the whole data against the given sub-schema.
						valid, err := v.validateJsonData("", jsonData.raw, rootSchemaId)
						if !valid {
							return false, KeywordValidationError{
								"dependencies",
								"inspected value failed in validation against sub-schema given in \"" +
									propertyName +
									"\" dependency: " +
									err.Error(),
							}
						}
					}
				}
			// In this case the dependency is a list of required property names.
			case []interface{}:
				{
					// Iterate over the items in the dependency array.
					for index, value := range v {
						// Verify that the value is actually a string.
						// If not, return an error
						if requiredProperty, ok := value.(string); ok {
							// Check if the required property name is missing. If it is,
							// return an error.
							if _, ok := object[requiredProperty]; !ok {
								return false, KeywordValidationError{
									"dependencies",
									"missing property \"" +
										requiredProperty +
										"\" although it is required according to \"" +
										propertyName +
										"\" dependency",
								}
							}
						} else {
							return false, KeywordValidationError{
								"dependencies",
								"all items in dependency array must be strings, item at position " +
									strconv.Itoa(index) +
									" is not a string",
							}
						}
					}
				}
			default:
				{
					return false, KeywordValidationError{
						"dependencies",
						"dependency value must be a json object or a json array",
					}
				}
			}
		}

		// If we arrived here it means that all the validations succeeded.
		return true, nil
	} else {
		return false, KeywordValidationError{
			"dependencies",
			"inspected value expected to be a json object",
		}
	}
}

type patternProperties map[string]*JsonSchema

func (pp patternProperties) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if pp == nil {
		return true, nil
	}

	// First we need to verify that jsonData is a json object.
	if object, ok := jsonData.value.(map[string]interface{}); ok {
		// Iterate over the given patterns.
		for pattern, subSchema := range pp {
			// Iterate over the properties in the inspected value.
			for property := range object {
				// Check if the property matches to the pattern.
				match, err := regexp.MatchString(pattern, property)

				// The pattern or the value is not in the right format (string)
				if err != nil {
					return false, KeywordValidationError{
						"patternProperties",
						err.Error(),
					}
				}

				// If there is a match, validate the value of the property against
				// the given schema.
				if match {
					valid, err := subSchema.validateJsonData(jsonPath+"/"+property, jsonData.raw, rootSchemaId)

					// If the validation fails, return an error.
					if !valid {
						return false, KeywordValidationError{
							"patternProperties",
							"property \"" +
								property +
								"\" that matches the pattern \"" +
								pattern +
								"\" failed in validation: \n" + err.Error(),
						}
					}
				}
			}
		}

		// If we arrived here it means that none of the properties failed in
		// validation against any of the given schemas.
		return true, nil
	} else {
		return false, KeywordValidationError{
			"patternProperties",
			"inspected value expected to be a json object",
		}
	}
}

type minProperties int

func (mp *minProperties) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if mp == nil {
		return true, nil
	}

	// First, we must verify that jsonData is a json object.
	// If it is not a json object, we return an error.
	if v, ok := jsonData.value.(map[string]interface{}); ok {
		// If the amount of keys in jsonData equals to or greater
		// than minProperties.
		// Else, return an error.
		if len(v) >= int(*mp) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"minProperties",
				"inspected value must contains at least " + strconv.Itoa(int(*mp)) + " properties",
			}
		}
	} else {
		return false, KeywordValidationError{
			"minProperties",
			"inspected value must be a json object",
		}
	}
}

type maxProperties int

func (mp *maxProperties) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if mp == nil {
		return true, nil
	}

	// First, we must verify that jsonData is a json object.
	// If it is not a json object, we return an error.
	if v, ok := jsonData.value.(map[string]interface{}); ok {
		// If the amount of keys in jsonData equals to or less
		// than maxProperties.
		// Else, return an error.
		if len(v) <= int(*mp) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"minProperties",
				"inspected value may contains at most " +
					strconv.Itoa(int(*mp)) +
					" properties",
			}
		}
	} else {
		return false, KeywordValidationError{
			"minProperties",
			"inspected value must be a json object",
		}
	}
}

/********************/
/** Array Keywords **/
/********************/

type items json.RawMessage

func (i items) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if i == nil {
		return true, nil
	}

	// First, we need to verify that json Data is an array
	if array, ok := jsonData.value.([]interface{}); ok {
		var data interface{}

		// Unmarshal the value in items in order to figure out if it is a
		// json object or json array
		err := json.Unmarshal(i, &data)
		if err != nil {
			return false, err
		}

		// Handle the value in items according to its json type.
		switch itemsField := data.(type) {
		// If jsonData is a json object, which means that is holds a single schema,
		// we validate the all the items in the inspected array against the given
		// schema.
		case map[string]interface{}:
			{
				// This is the JsonSchema instance that should hold the schema in
				// "items" field.
				var schema JsonSchema

				// Unmarshal the rawSchema into the JsonSchema struct.
				err = json.Unmarshal(i, &schema)
				if err != nil {
					return false, err
				}

				// Iterate over the items in the inspected array and validate each
				// item against the schema in "items" field.
				for index := 0; index < len(array); index++ {
					valid, err := schema.validateJsonData(jsonPath+"/"+strconv.Itoa(index), jsonData.raw, rootSchemaId)
					if !valid {
						return valid, err
					}
				}

				// If we arrived here it means that all the items in the inspected array
				// validated successfully against the given schema.
				return true, nil
			}
		// If jsonData is a json array, which means that is holds multiple json schema objects,
		// we validate each item in the inspected array against the schema at the same position.
		case []interface{}:
			{
				if len(itemsField) > len(array) {
					return false, KeywordValidationError{
						"items",
						"when \"items\" field contains a list of Json Schema objects, the " +
							"inspected array must contain at least the same amount of items",
					}
				}

				// Iterate over the schemas in "items" field.
				for index, schemaFromItems := range itemsField {
					// Marshal the current schema in "items" field in order to Unmarshal it
					// into JsonSchema instance.
					rawSchema, err := json.Marshal(schemaFromItems)
					if err != nil {
						return false, err
					}

					// This is the JsonSchema instance that should hold the current
					// working schema.
					var schema JsonSchema

					// Unmarshal the rawSchema into the JsonSchema struct.
					err = json.Unmarshal(rawSchema, &schema)
					if err != nil {
						return false, err
					}

					// Validate the item against the schema at the same position.
					valid, err := schema.validateJsonData(jsonPath+"/"+strconv.Itoa(index), jsonData.raw, rootSchemaId)
					if !valid {
						return valid, err
					}
				}

				// If we arrived here it means that all the items in the inspected array
				// validated successfully against corresponding schema.
				return true, nil
			}
		// The default case indicates that the value in items field is not a json schema or
		// a list of json schema.
		default:
			{
				return false, KeywordValidationError{
					"items",
					"\"items\" field value in schema must be a valid Json Schema or an array of Json Schema",
				}
			}
		}
	} else {
		return false, KeywordValidationError{
			"items",
			"inspected value expected to be a json array",
		}
	}
}

func (i *items) UnmarshalJSON(data []byte) error {
	*i = data
	return nil
}

type additionalItems struct {
	JsonSchema
	siblingItems *items
}

func (ai *additionalItems) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if ai == nil {
		return true, nil
	}

	// Unmarshal the sibling field "items" in order to check it's json type.
	var siblingItems interface{}
	err := json.Unmarshal([]byte(*ai.siblingItems), &siblingItems)
	if err != nil {
		return false, err
	}

	// If "items" is a json array, "additionalItems" needs to verify the items
	// that the schema in "items" field did not validate.
	if itemsArray, ok := siblingItems.([]interface{}); ok {
		// Check if jsonData is a json array.
		if array, ok := jsonData.value.([]interface{}); ok {
			// Iterate over the inspected array from the position that items stopped
			// validating.
			for index := range array[len(itemsArray):] {
				// Validate the inspected item against the schema given in "additionalItems".
				valid, err := ai.validateJsonData(jsonPath+"/"+strconv.Itoa(index), jsonData.raw, rootSchemaId)
				if !valid {
					return false, KeywordValidationError{
						"additionalItems",
						"item at position " +
							strconv.Itoa(index) +
							" failed in validation: " +
							err.Error(),
					}
				}
			}

			// If we arrived here it means that no item failed in validation.
			return true, nil
		} else {
			return false, KeywordValidationError{
				"additionalItems",
				"inspected value expected to be a json array",
			}
		}
	} else {
		// If "items" field is not an array of json schema, additionalItems
		// is meaningless so we return true.
		return true, nil
	}
}

type contains struct {
	JsonSchema
}

func (c *contains) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if c == nil {
		return true, nil
	}

	// First, we need to verify that jsonData is a json array.
	if array, ok := jsonData.value.([]interface{}); ok {
		// Go over all the items in the array in order to inspect them.
		for index := range array {
			// If the item is valid against the given schema, which means that
			// the array contains the required value.
			valid, _ := (*c).validateJsonData(jsonPath+"/"+strconv.Itoa(index), jsonData.raw, rootSchemaId)
			if valid {
				return true, nil
			}
		}
	}

	// If we arrived here it means that we could not validate any of the array's
	// items against the given schema.
	return false, KeywordValidationError{
		"contains",
		"could validate any of the inspected array's items against the given schema",
	}
}

type minItems int

func (mi *minItems) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if mi == nil {
		return true, nil
	}

	// First, we need to verify that jsonData is an array.
	if v, ok := jsonData.value.([]interface{}); ok {
		// Check that the number of items in the array is equal to
		// or greater than minItems.
		if len(v) >= int(*mi) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"minItems",
				"inspected array must contain at least " + strconv.Itoa(int(*mi)) + " items",
			}
		}
	} else {
		return false, KeywordValidationError{
			"minItems",
			"inspected value expected to be json array",
		}
	}
}

type maxItems int

func (mi *maxItems) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if mi == nil {
		return true, nil
	}

	// First, we need to verify that jsonData is an array.
	if v, ok := jsonData.value.([]interface{}); ok {
		// Check that the number of items in the array is equal to
		// or less than maxItems.
		if len(v) <= int(*mi) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"maxItems",
				"inspected array must contain at most " + strconv.Itoa(int(*mi)) + " items",
			}
		}
	} else {
		return false, KeywordValidationError{
			"maxItems",
			"inspected value expected to be json array",
		}
	}
}

type uniqueItems bool

func (ui *uniqueItems) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if ui == nil {
		return true, nil
	}

	// First, we need to verify that jsonData is an array.
	if array, ok := jsonData.value.([]interface{}); ok {
		// Create a map that will help us to check if we already met the
		// item by using the map's hashing mechanism.
		uniqueSet := make(map[string]int)

		// Iterate over the items in the inspected array.
		for index, item := range array {
			// Marshal the item back to hash-able value, because maps (json object)
			// and slices (json arrays) are not a hash-able values.
			rawItem, err := json.Marshal(item)
			if err != nil {
				return false, err
			}

			// If ok is true it means that the value exists in the map, which means
			// we already met it in one of the previous iterations.
			// Else, insert the item into the map as key, and the index as value.
			if v, ok := uniqueSet[string(rawItem)]; ok {
				return false, KeywordValidationError{
					"uniqueItems",
					"the inspected array contains two equal items at indices: " +
						strconv.Itoa(v) +
						", " +
						strconv.Itoa(index),
				}
			} else {
				uniqueSet[string(rawItem)] = index
			}
		}

		// If we arrived here it means that we did not meat any item which is
		// similar to another item in the array.
		return true, nil
	} else {
		return false, KeywordValidationError{
			"uniqueItems",
			"inspected value expected to be json array",
		}
	}
}

/********************/
/** Other Keywords **/
/********************/

type contentMediaType string
type contentEncoding string

/**************************/
/** Conditional Keywords **/
/**************************/

type anyOf []*JsonSchema

func (af anyOf) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if af == nil {
		return true, nil
	}

	// Validate jsonData.raw against each of the schemas until on of them succeeds.
	for _, schema := range af {
		valid, err := schema.validateJsonData("", jsonData.raw, rootSchemaId)
		if valid {
			return valid, err
		}
	}

	// If we arrived here, the validation of jsonData failed against all schemas.
	return false, KeywordValidationError{
		"anyOf",
		"inspected value could not be validated against any of the given schemas",
	}
}

type allOf []*JsonSchema

func (af allOf) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if af == nil {
		return true, nil
	}

	// Validate jsonData.raw against each of the schemas.
	// If one of them fails, return error.
	for _, schema := range af {
		valid, err := schema.validateJsonData("", jsonData.raw, rootSchemaId)
		if !valid {
			return valid, err
		}
	}

	// If we arrived here, the validation of jsonData succeeded against all
	// given schemas.
	return false, KeywordValidationError{
		"allOf",
		"inspected value could not be validated against all of the given schemas",
	}
}

type oneOf []*JsonSchema

func (of oneOf) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if of == nil {
		return true, nil
	}

	var oneValidationAlreadySucceeded bool

	// Validate jsonData.raw against each of the schemas until on of them succeeds.
	for _, schema := range of {
		valid, _ := schema.validateJsonData("", jsonData.raw, rootSchemaId)
		if valid {
			if oneValidationAlreadySucceeded {
				return false, KeywordValidationError{
					"oneOf",
					"inspected data is valid against more than one given schema",
				}
			} else {
				oneValidationAlreadySucceeded = true
			}
		}
	}

	if oneValidationAlreadySucceeded {
		return true, nil
	} else {
		// If we arrived here, the validation of jsonData failed against all schemas.
		return false, KeywordValidationError{
			"oneOf",
			"inspected value could not be validated against any of the given schemas",
		}
	}
}

type not struct {
	JsonSchema
}

func (n *not) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if n == nil {
		return true, nil
	}

	valid, _ := (*n).validateJsonData(jsonPath, jsonData.raw, rootSchemaId)
	if !valid {
		return true, nil
	} else {
		return false, KeywordValidationError{
			"not",
			"inspected value did not fail on validation against the schema defined by this keyword",
		}
	}
}

type _if struct {
	JsonSchema
	siblingThen *_then
	siblingElse *_else
}

func (i *_if) validate(jsonPath string, jsonData jsonData, rootSchemaId string) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if i == nil {
		return true, nil
	}

	// Validate the data against the given schema in "if".
	valid, _ := i.validateJsonData("", jsonData.raw, rootSchemaId)

	// If the validation succeeded, validate the data against the given schema
	// in "then".
	// Else, validate the data against the given schema in "else".
	if valid {
		return i.siblingThen.validateJsonData(jsonPath, jsonData.raw, rootSchemaId)
	} else {
		return i.siblingElse.validateJsonData(jsonPath, jsonData.raw, rootSchemaId)
	}
}

type _then struct {
	JsonSchema
}

type _else struct {
	JsonSchema
}

/****************************/
/** Authorization Keywords **/
/****************************/

type readOnly bool
type writeOnly bool
