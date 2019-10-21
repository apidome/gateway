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
		description string
		method      string
		path        string
		schema      string
		data        string
		valid       bool
	}{
		{
			"",
			"GET",
			"/v1/a",
			`
				{
					
				}
			`,
			`
				{
					
				}	
			`,
			true,
		},
		{
			"",
			"GET",
			"/v1/b",
			`
				{
					
				}
			`,
			`
				{
					
				}	
			`,
			true,
		},
		{
			"",
			"GET",
			"/v1/c",
			`
				{
					
				}
			`,
			`
				{
					
				}	
			`,
			true,
		},
		{
			"",
			"GET",
			"/v1/d",
			`
				{
					
				}
			`,
			`
				{
					
				}	
			`,
			true,
		},
		{
			"",
			"GET",
			"/v1/e",
			`
				{
					
				}
			`,
			`
				{
					
				}	
			`,
			true,
		},
		{
			"",
			"GET",
			"/v1/f",
			`
				{
					
				}
			`,
			`
				{
					
				}	
			`,
			true,
		},
		{
			"",
			"GET",
			"/v1/g",
			`
				{
					
				}
			`,
			`
				{
					
				}	
			`,
			true,
		},
		{
			"",
			"GET",
			"/v1/h",
			`
				{
					
				}
			`,
			`
				{
					
				}	
			`,
			true,
		},
	}

	t.Log("Given the need to test json validation against json schema according to method and endpoint")
	{
		for index, testCase := range testCases {
			t.Logf("\tTest %d: When trying to validate %s against the schema belongs to %s %s",
				index, testCase.description, testCase.method, testCase.path)
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
						t.Logf("\t%s\tShould not be valid against the specified json schema: %v", succeed, err)
					} else {
						t.Errorf("\t%s\tShould not be valid against the specified json schema", failed)
					}
				}
			}
		}
	}
}
