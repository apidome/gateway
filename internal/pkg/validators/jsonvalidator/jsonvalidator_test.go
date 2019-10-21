package jsonvalidator_test

import (
	"github.com/Creespye/caf/internal/pkg/validators/jsonvalidator"
	"testing"
)

//const succeed = "\u2713"
//const failed = "\u2717"
const succeed = "V"
const failed = "X"

func TestNewJsonValidator(t *testing.T) {
	testCases := []struct {
		draft string
		valid bool
	}{
		{
			"draft-07",
			true,
		},
		{
			"draft-06",
			false,
		},
		{
			"",
			false,
		},
	}

	t.Log("Given the need to test creation of new JsonValidator")
	{
		for index, testCase := range testCases {
			t.Logf("\tTest %d: When trying to create a JsonValidator with %s", index, testCase.draft)
			{
				if testCase.valid {
					if _, err := jsonvalidator.NewJsonValidator(testCase.draft); err != nil {
						t.Errorf("\t%s\tShould be able to get a reference to a JsonValidator: %v", failed, err)
					} else {
						t.Logf("\t%s\tShould be able to get a reference to a JsonValidator", succeed)
					}
				} else {
					if _, err := jsonvalidator.NewJsonValidator(testCase.draft); err == nil {
						t.Errorf("\t%s\tShould not be able to get a reference to a JsonValidator", failed)
					} else {
						t.Logf("\t%s\tShould not be able to get a reference to a JsonValidator: %v", succeed, err)
					}
				}
			}
		}
	}
}

func TestLoadSchema(t *testing.T) {
	testCases := []struct {
		description string
		method      string
		path        string
		schema      string
		valid       bool
	}{
		{
			"the json boolean \"true\" as a schema",
			"GET",
			"/v1/a",
			"true",
			true,
		},
		{
			"the json boolean \"false\" as a schema",
			"GET",
			"/v1/a",
			"false",
			true,
		},
		{
			"empty json object as a schema",
			"GET",
			"/v1/a",
			"{}",
			true,
		},
		{
			"a valid json schema as a schema",
			"GET",
			"/v1/a",
			"{\"type\": \"string\"}",
			true,
		},
		{
			"a json object that contains only a non-standard keywords as a schema",
			"GET",
			"/v1/a",
			"{\"someNonStandardKeyword\": 4}",
			true,
		},
		{
			"any json string as a schema",
			"GET",
			"/v1/a",
			"'someJsonString'",
			false,
		},
		{
			"any json number as a schema",
			"GET",
			"/v1/a",
			"45.7",
			false,
		},
		{
			"\"GET\" as method",
			"GET",
			"/v1/a",
			"{}",
			true,
		},
		{
			"\"POST\" as method",
			"POST",
			"/v1/a",
			"{}",
			true,
		},
		{
			"\"PUT\" as method",
			"PUT",
			"/v1/a",
			"{}",
			true,
		},
		{
			"\"PATCH\" as method",
			"PATCH",
			"/v1/a",
			"{}",
			true,
		},
		{
			"\"DELETE\" as method",
			"DELETE",
			"/v1/a",
			"{}",
			true,
		},
		{
			"a non-standard http method - \"GET1\" as method",
			"GET1",
			"/v1/a",
			"{}",
			false,
		},
	}

	t.Log("Given the need to test loading of new json schema to JsonValidator")
	{
		for index, testCase := range testCases {
			t.Logf("\tTest %d: When trying to load %s", index, testCase.description)
			{
				jv, err := jsonvalidator.NewJsonValidator("draft-07")
				if err != nil {
					t.Fatalf("\t%s\tShould be able to create a new JsonValidator: %v", failed, err)
				}
				t.Logf("\t%s\tShould be able to create a new JsonValidator", succeed)

				err = jv.LoadSchema(testCase.path, testCase.method, []byte(testCase.schema))
				if testCase.valid {
					if err != nil {
						t.Errorf("\t%s\tShould be able to Load schema: %v", failed, err)
					} else {
						t.Logf("\t%s\tShould be able to Load schema", succeed)
					}
				} else {
					if err != nil {
						t.Logf("\t%s\tShould not be able to Load schema: %v", succeed, err)
					} else {
						t.Errorf("\t%s\tShould not be able to Load schema", failed)
					}
				}
			}
		}
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		testedJsonSchemaKeyword string
		schemaDescription       string
		dataDescription         string
		method                  string
		path                    string
		schema                  string
		data                    string
		valid                   bool
	}{
		{
			"type",
			"a json schema that accepts any json string",
			"a json string",
			"GET",
			"/v1/a",
			`
				{
					"type": "string"
				}
			`,
			`"some json string"`,
			true,
		},
		{
			"type",
			"a json schema that accepts any json string",
			"a json object",
			"GET",
			"/v1/a",
			`
				{
					"type": "string"
				}
			`,
			`{}`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json string",
			"a json array",
			"GET",
			"/v1/a",
			`
				{
					"type": "string"
				}
			`,
			`[3, 4]`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json string",
			"a json boolean",
			"GET",
			"/v1/a",
			`
				{
					"type": "string"
				}
			`,
			`true`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json string",
			"a json null",
			"GET",
			"/v1/a",
			`
				{
					"type": "string"
				}
			`,
			`null`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json string",
			"a json number",
			"GET",
			"/v1/a",
			`
				{
					"type": "string"
				}
			`,
			`78.2`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json string",
			"a json integer",
			"GET",
			"/v1/a",
			`
				{
					"type": "string"
				}
			`,
			`8`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json number (float and integer)",
			"a json number",
			"GET",
			"/v1/b",
			`
				{
					"type": "number"
				}
			`,
			`78.9`,
			true,
		},
		{
			"type",
			"a json schema that accepts any json number (float and integer)",
			"a json integer",
			"GET",
			"/v1/b",
			`
				{
					"type": "number"
				}
			`,
			`99`,
			true,
		},
		{
			"type",
			"a json schema that accepts any json number (float and integer)",
			"a json object",
			"GET",
			"/v1/b",
			`
				{
					"type": "number"
				}
			`,
			`{}`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json number (float and integer)",
			"a json array",
			"GET",
			"/v1/b",
			`
				{
					"type": "number"
				}
			`,
			`["a", "b", "c"]`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json number (float and integer)",
			"a json string",
			"GET",
			"/v1/b",
			`
				{
					"type": "number"
				}
			`,
			`"some json string"`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json number (float and integer)",
			"a json boolean",
			"GET",
			"/v1/b",
			`
				{
					"type": "number"
				}
			`,
			`false`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json number (float and integer)",
			"a json null",
			"GET",
			"/v1/b",
			`
				{
					"type": "number"
				}
			`,
			`null`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json integer",
			"a json object",
			"GET",
			"/v1/b",
			`
				{
					"type": "integer"
				}
			`,
			`{}`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json integer",
			"a json array",
			"GET",
			"/v1/b",
			`
				{
					"type": "integer"
				}
			`,
			`["", "4"]`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json integer",
			"a json string",
			"GET",
			"/v1/b",
			`
				{
					"type": "integer"
				}
			`,
			`"some json string"`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json integer",
			"a json boolean",
			"GET",
			"/v1/b",
			`
				{
					"type": "integer"
				}
			`,
			`true`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json integer",
			"a json null",
			"GET",
			"/v1/b",
			`
				{
					"type": "integer"
				}
			`,
			`null`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json integer",
			"a json integer",
			"GET",
			"/v1/b",
			`
				{
					"type": "integer"
				}
			`,
			`45`,
			true,
		},
		{
			"type",
			"a json schema that accepts any json integer",
			"a json number (which is not an integer)",
			"GET",
			"/v1/b",
			`
				{
					"type": "integer"
				}
			`,
			`33.3`,
			false,
		},
		{
			"type",
			"a json schema that accepts json boolean",
			"a json object",
			"GET",
			"/v1/b",
			`
				{
					"type": "boolean"
				}
			`,
			`{}`,
			false,
		},
		{
			"type",
			"a json schema that accepts json boolean",
			"a json array",
			"GET",
			"/v1/b",
			`
				{
					"type": "boolean"
				}
			`,
			`[]`,
			false,
		},
		{
			"type",
			"a json schema that accepts json boolean",
			"a json string",
			"GET",
			"/v1/b",
			`
				{
					"type": "boolean"
				}
			`,
			`"some json string"`,
			false,
		},
		{
			"type",
			"a json schema that accepts json boolean",
			"a json number",
			"GET",
			"/v1/b",
			`
				{
					"type": "boolean"
				}
			`,
			`78.85`,
			false,
		},
		{
			"type",
			"a json schema that accepts json boolean",
			"a json integer",
			"GET",
			"/v1/b",
			`
				{
					"type": "boolean"
				}
			`,
			`7`,
			false,
		},
		{
			"type",
			"a json schema that accepts json boolean",
			"a json null",
			"GET",
			"/v1/b",
			`
				{
					"type": "boolean"
				}
			`,
			`null`,
			false,
		},
		{
			"type",
			"a json schema that accepts json boolean",
			"a json boolean",
			"GET",
			"/v1/b",
			`
				{
					"type": "boolean"
				}
			`,
			`false`,
			true,
		},
		{
			"type",
			"a json schema that accepts any json object",
			"a json boolean",
			"GET",
			"/v1/b",
			`
				{
					"type": "object"
				}
			`,
			`true`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json object",
			"a json string",
			"GET",
			"/v1/b",
			`
				{
					"type": "object"
				}
			`,
			`"some json string"`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json object",
			"a json array",
			"GET",
			"/v1/b",
			`
				{
					"type": "object"
				}
			`,
			`[]`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json object",
			"a json integer",
			"GET",
			"/v1/b",
			`
				{
					"type": "object"
				}
			`,
			`9`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json object",
			"a json number",
			"GET",
			"/v1/b",
			`
				{
					"type": "object"
				}
			`,
			`78.5`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json object",
			"a json null",
			"GET",
			"/v1/b",
			`
				{
					"type": "object"
				}
			`,
			`null`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json object",
			"a json ",
			"GET",
			"/v1/b",
			`
				{
					"type": "object"
				}
			`,
			`{}`,
			true,
		},
		{
			"type",
			"a json schema that accepts any json object",
			"a json array",
			"GET",
			"/v1/b",
			`
				{
					"type": "array"
				}
			`,
			`[4, "a", null]`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json array",
			"a json boolean",
			"GET",
			"/v1/b",
			`
				{
					"type": "array"
				}
			`,
			`true`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json array",
			"a json object",
			"GET",
			"/v1/b",
			`
				{
					"type": "array"
				}
			`,
			`{}`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json array",
			"a json number",
			"GET",
			"/v1/b",
			`
				{
					"type": "array"
				}
			`,
			`34.5`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json array",
			"a json integer",
			"GET",
			"/v1/b",
			`
				{
					"type": "array"
				}
			`,
			`45`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json array",
			"a json string",
			"GET",
			"/v1/b",
			`
				{
					"type": "array"
				}
			`,
			`"some json string"`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json array",
			"a json null",
			"GET",
			"/v1/b",
			`
				{
					"type": "array"
				}
			`,
			`null`,
			false,
		},
		{
			"type",
			"a json schema that accepts only null",
			"a json null",
			"GET",
			"/v1/b",
			`
				{
					"type": "null"
				}
			`,
			`null`,
			true,
		},
		{
			"type",
			"a json schema that accepts only null",
			"a json integer",
			"GET",
			"/v1/b",
			`
				{
					"type": "null"
				}
			`,
			`13`,
			false,
		},
		{
			"type",
			"a json schema that accepts only null",
			"a json number",
			"GET",
			"/v1/b",
			`
				{
					"type": "null"
				}
			`,
			`6.6`,
			false,
		},
		{
			"type",
			"a json schema that accepts only null",
			"a json array",
			"GET",
			"/v1/b",
			`
				{
					"type": "null"
				}
			`,
			`[true]`,
			false,
		},
		{
			"type",
			"a json schema that accepts only null",
			"a json object",
			"GET",
			"/v1/b",
			`
				{
					"type": "null"
				}
			`,
			`{}`,
			false,
		},
		{
			"type",
			"a json schema that accepts only null",
			"a json string",
			"GET",
			"/v1/b",
			`
				{
					"type": "null"
				}
			`,
			`"some json string"`,
			false,
		},
		{
			"type",
			"a json schema that accepts only null",
			"a json boolean",
			"GET",
			"/v1/b",
			`
				{
					"type": "null"
				}
			`,
			`true`,
			false,
		},
		{
			"type",
			"a json schema that accepts any json object or array",
			"a json object",
			"GET",
			"/v1/b",
			`
				{
					"type": ["object", "array"]
				}
			`,
			`{}`,
			true,
		},
		{
			"type",
			"a json schema that accepts any json object or array",
			"a json array",
			"GET",
			"/v1/b",
			`
				{
					"type": ["object", "array"]
				}
			`,
			`[]`,
			true,
		},
		{
			"type",
			"a json schema that accepts any json object or array",
			"a json boolean",
			"GET",
			"/v1/b",
			`
				{
					"type": ["object", "array"]
				}
			`,
			`{}`,
			false,
		},
	}

	t.Log("Given the need to test json validation against json schema according to method and endpoint")
	{
		for index, testCase := range testCases {
			t.Logf("\tTest %d: When trying to validate %s against %s",
				index, testCase.dataDescription, testCase.schemaDescription)
			{
				jv, err := jsonvalidator.NewJsonValidator("draft-07")
				if err != nil {
					t.Fatalf("\t%s\tShould be able to create a new JsonValidator: %v", failed, err)
				}
				t.Logf("\t%s\tShould be able to create a new JsonValidator", succeed)

				err = jv.LoadSchema(testCase.path, testCase.method, []byte(testCase.schema))
				if err != nil {
					t.Errorf("\t%s\tShould be able to Load schema: %v", failed, err)
				} else {
					t.Logf("\t%s\tShould be able to Load schema", succeed)
				}

				err = jv.Validate(testCase.path, testCase.method, []byte(testCase.data))
				if testCase.valid {
					if err != nil {
						t.Errorf("\t%s\tData should be valid against the specified json schema: %v", failed, err)
					} else {
						t.Logf("\t%s\tData should be valid against the specified json schema", succeed)
					}
				} else {
					if err != nil {
						t.Logf("\t%s\tData should not be valid against the specified json schema: %v", succeed, err)
					} else {
						t.Errorf("\t%s\tData should not be valid against the specified json schema", failed)
					}
				}
			}
		}
	}
}
