package jsonvalidator

import (
	"encoding/json"
	"math"
	"strconv"
)

type keywordValidator interface {
	validate(interface{}) (bool, error)
}

/**********************/
/** Generic Keywords **/
/**********************/

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

type enum []interface{}

func (e enum) validate(jsonData interface{}) (bool, error) {
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

type _const json.RawMessage

func (c _const) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type definitions map[string]*JsonSchema

func (d definitions) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

type _type json.RawMessage

func (t *_type) validate(jsonData interface{}) (bool, error) {
	return true, nil
}

func (t *_type) UnmarshalJSON(data []byte) error {
	*t = data
	return nil
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
