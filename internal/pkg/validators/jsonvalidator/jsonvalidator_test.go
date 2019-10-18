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
			"the json boolean \"true\"",
			"GET",
			"/v1/a",
			"true",
			true,
		},
		{
			"the json boolean \"false\"",
			"GET",
			"/v1/a",
			"false",
			true,
		},
		{
			"empty schema",
			"GET",
			"/v1/a",
			"{}",
			true,
		},
		{
			"a valid json schema",
			"GET",
			"/v1/a",
			"{\"type\": \"string\"}",
			true,
		},
		{
			"a json object that is not a json schema",
			"GET",
			"/v1/a",
			"{\"someNonStandardKeyword\": 4}",
			false,
		},
		{
			"any json string",
			"GET",
			"/v1/a",
			"'someJsonString'",
			false,
		},
		{
			"any json number",
			"GET",
			"/v1/a",
			"45.7",
			false,
		},
	}

	t.Log("Given the need to test loading of new json schema to JsonValidator")
	{
		for index, testCase := range testCases {
			t.Logf("\tTest %d: When trying to load %s as a schema", index, testCase.description)
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
			}
		}
	}
}

func TestValidate(t *testing.T) {

}
