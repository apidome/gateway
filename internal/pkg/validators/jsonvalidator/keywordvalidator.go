package jsonvalidator

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
)

/*
Implemented keywordValidators:
> schema: 					X
> ref: 						X
> id: 						X
> comment: 					X
> title: 					X
> description: 				X
> examples: 				X
> enum: 					X
> _default: 				X
> _const: 					X
> definitions: 				X
> _type: 					V
> minLength: 				V
> maxLength: 				X
> pattern: 					X
> format: 					X
> multipleOf: 				V
> minimum: 					V
> maximum: 					V
> exclusiveMinimum: 		V
> exclusiveMaximum: 		V
> properties: 				V
> additionalProperties: 	X
> required: 				Y
> propertyNames: 			X
> dependencies: 			X
> patternProperties: 		X
> minProperties: 			X
> maxProperties: 			X
> items: 					X
> contains: 				X
> additionalItems: 			X
> minItems: 				X
> maxItems: 				X
> uniqueItems: 				X
> contentMediaType: 		X
> contentEncoding: 			X
> anyOf: 					X
> allOf: 					X
> oneOf: 					X
> not: 						X
> _if: 						X
> _then: 					X
> _else: 					X
> readOnly: 				X
> writeOnly: 				X
*/

type keywordValidator interface {
	validate(interface{}) (bool, error)
}

/*****************/
/** Annotations **/
/*****************/

type schema string

func (s *schema) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type ref string

func (r *ref) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type id string

func (i *id) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type comment string

func (c *comment) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type title string

func (t *title) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type description string

func (d *description) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type examples []interface{}

func (e examples) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type _default json.RawMessage

func (d _default) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

func (d *_default) UnmarshalJSON(data []byte) error {
	*d = data
	return nil
}

/**********************/
/** Generic Keywords **/
/**********************/

type _type json.RawMessage

func (t *_type) validate(jsonData interface{}) (bool, error) {
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
					ok, _ := assertJsonType(v, jsonData)

					// If the assertion succeeded, return true
					if ok {
						return ok, nil
					}
				} else {
					return false, KeywordValidationError{
						"type",
						"type field in schema must be string or array of strings",
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
			return assertJsonType(typeFromSchema, jsonData)
		}
	default:
		{
			return false, KeywordValidationError{
				"type",
				"type field in schema must be string or array of strings",
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
	//		if v, ok := jsonData.(string); ok {
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
				"type field in schema must be string or array of strings",
			}
		}
	}
}

func (t *_type) UnmarshalJSON(data []byte) error {
	*t = data
	return nil
}

type enum []interface{}

func (e enum) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type _const json.RawMessage

func (c _const) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

/*********************/
/** String Keywords **/
/*********************/

type minLength int

func (ml *minLength) validate(jsonData interface{}) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if ml == nil {
		return true, nil
	}

	// If jsonData is a string, validate its length,
	// else, return a KeywordValidationError
	if v, ok := jsonData.(string); ok {
		if len(v) >= int(*ml) {
			return true, nil
		} else {
			return false, KeywordValidationError{
				"minLength",
				"inspected string shorter than " + string(*ml),
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

func (ml *maxLength) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type pattern string

func (p *pattern) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type format string

func (f *format) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

/*********************/
/** Number Keywords **/
/*********************/

type multipleOf float64

func (mo *multipleOf) validate(jsonData interface{}) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if mo == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.(float64); ok {
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

func (m *minimum) validate(jsonData interface{}) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if m == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.(float64); ok {
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

func (m *maximum) validate(jsonData interface{}) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if m == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.(float64); ok {
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

func (em *exclusiveMinimum) validate(jsonData interface{}) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if em == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.(float64); ok {
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
			"",
		}
	}
}

type exclusiveMaximum float64

func (em *exclusiveMaximum) validate(jsonData interface{}) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if em == nil {
		return true, nil
	}

	// If jsonData is float64, validate it. Else, return KeywordValidationError
	if v, ok := jsonData.(float64); ok {
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

func (p properties) validate(jsonData interface{}) (bool, error) {
	// If the receiver is nil, dont validate it (return true)
	if p == nil {
		return true, nil
	}

	var rawData json.RawMessage
	var err error

	// If the jsonData is already json.RawMessage, use it.
	// Else, Marshal it back to []byte (which is similar to json.RawMessage)
	// because JsonSchema.validateJsonData() requires a slice of bytes.
	if v, ok := jsonData.(json.RawMessage); ok {
		rawData = v
	} else {
		rawData, err = json.Marshal(jsonData)
		if err != nil {
			return false, err
		}
	}

	// For each "property" validate it according to its JsonSchema.
	for key, value := range p {
		valid, err := value.validateJsonData("/"+key, rawData)
		if err != nil {
			return valid, err
		}
	}

	return true, nil
}

type additionalProperties JsonSchema

func (ap *additionalProperties) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

//func (ap *additionalProperties) UnmarshalJSON(data []byte) error {
//	*ap = data
//	return nil
//}

type required []string

func (r required) validate(jsonData interface{}) (bool, error) {
	if r == nil {
		return true, nil
	}

	if v, ok := jsonData.(map[string]interface{}); ok {
		for _, property := range r {
			if v[property] == nil {
				return false, errors.New("Missing required property - " + property)
			}
		}
	}

	return true, nil
}

type propertyNames JsonSchema

func (pn *propertyNames) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type dependencies map[string]json.RawMessage

func (d dependencies) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type patternProperties map[string]*JsonSchema

func (pp patternProperties) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type minProperties int

func (mp *minProperties) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type maxProperties int

func (mp *maxProperties) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type definitions map[string]*JsonSchema

func (d definitions) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

/********************/
/** Array Keywords **/
/********************/

type items json.RawMessage

func (i items) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

func (i *items) UnmarshalJSON(data []byte) error {
	*i = data
	return nil
}

type contains json.RawMessage

func (c contains) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

func (c *contains) UnmarshalJSON(data []byte) error {
	*c = data
	return nil
}

type additionalItems json.RawMessage

func (ai additionalItems) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

func (ai *additionalItems) UnmarshalJSON(data []byte) error {
	*ai = data
	return nil
}

type minItems int

func (mi *minItems) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type maxItems int

func (mi *maxItems) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type uniqueItems bool

func (ui *uniqueItems) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

/********************/
/** Other Keywords **/
/********************/

type contentMediaType string

func (cm *contentMediaType) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type contentEncoding string

func (ce *contentEncoding) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

/**************************/
/** Conditional Keywords **/
/**************************/

type anyOf []*JsonSchema

func (af anyOf) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type allOf []*JsonSchema

func (af allOf) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type oneOf []*JsonSchema

func (of oneOf) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type not JsonSchema

func (n *not) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type _if JsonSchema

func (i *_if) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type _then JsonSchema

func (t *_then) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type _else JsonSchema

func (e *_else) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

/****************************/
/** Authorization Keywords **/
/****************************/

type readOnly bool

func (ro *readOnly) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type writeOnly bool

func (wo *writeOnly) validate(jsonData interface{}) (bool, error) {
	return true, nil
}
