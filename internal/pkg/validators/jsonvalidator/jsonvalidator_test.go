package jsonvalidator_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/omeryahud/caf/internal/pkg/validators/jsonvalidator"
)

//const succeed = "\u2713"
//const failed = "\u2717"
const succeed = "V"
const failed = "X"

type testCase struct {
	Keyword      string
	Descriptions string          `json:"description"`
	Schema       json.RawMessage `json:"schema"`
	Path         string          `json:"path"`
	Method       string          `json:"Method"`
	Tests        []struct {
		Description string          `json:"description"`
		Data        json.RawMessage `json:"data"`
		Valid       bool            `json:"valid"`
	} `json:"tests"`
}

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
			"a schema that accepts only empty object or \"true\" using the enum keyword",
			"GET",
			"/v1/a",
			`{
			"enum": [{}, true]
		}`,
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
		jv, err := jsonvalidator.NewJsonValidator("draft-07")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create a new JsonValidator: %v", failed, err)
		}
		t.Logf("\t%s\tShould be able to create a new JsonValidator", succeed)

		for index, testCase := range testCases {
			t.Logf("\tTest %d: When trying to load %s", index, testCase.description)
			{
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
	keywords := []string{"type", "enum", "const", "minLength", "maxLength", "pattern", "format", "multipleOf",
		"minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum", "properties", "patternProperties",
		"additionalProperties", "required", "propertyNames", "minProperties", "maxProperties", "items", "contains",
		"additionalItems", "minItems", "maxItems", "uniqueItems", "anyOf", "allOf", "oneOf", "not",
		"if_then_else", "ref"}
	testCases := make([]testCase, 0)

	// Read all the test data from the files and append them to the main slice.
	for _, keyword := range keywords {
		testData := make([]testCase, 0)
		rawTestData, err := readTestDataFromFile(keyword + ".json")
		if err != nil {
			t.Fatalf("Could not read test data from file: %v", err)
		}

		err = json.Unmarshal(rawTestData, &testData)
		if err != nil {
			t.Fatalf("Could not unmarshal test data to test cases slice, "+
				"probably one or more cases is not in the correct format in %s.json: %v", keyword, err)
		}

		for index := range testData {
			testData[index].Keyword = keyword
		}

		testCases = append(testCases, testData...)
	}

	t.Log("Given the need to test json validation against json schema according to method and endpoint")
	{
		jv, err := jsonvalidator.NewJsonValidator("draft-07")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create a new JsonValidator: %v", failed, err)
		}
		t.Logf("\t%s\tShould be able to create a new JsonValidator", succeed)

		for i, testCase := range testCases {
			subTest := func(t *testing.T) {
				t.Logf("\t[%s] Test Schema %d: %s", testCase.Keyword, i, testCase.Descriptions)
				{
					for j, test := range testCase.Tests {
						t.Logf("\t\tTest %d.%d: When trying to validate %s against the given schema", i, j, test.Description)
						{
							err = jv.LoadSchema(testCase.Path, testCase.Method, testCase.Schema)
							if err != nil {
								t.Errorf("\t\t%s\tShould be able to Load schema: %v", failed, err)
							}

							err = jv.Validate(testCase.Path, testCase.Method, test.Data)
							if test.Valid {
								if err != nil {
									t.Errorf("\t\t%s\tData should be valid against the specified json schema: %v", failed, err)
								} else {
									t.Logf("\t\t%s\tData should be valid against the specified json schema", succeed)
								}
							} else {
								if err != nil {
									t.Logf("\t\t%s\tData should not be valid against the specified json schema: %v", succeed, err)
								} else {
									t.Errorf("\t\t%s\tData should not be valid against the specified json schema", failed)
								}
							}
						}
					}
				}
				t.Log()
			}

			t.Run(testCase.Keyword, subTest)
		}
	}
}

func readTestDataFromFile(fileName string) ([]byte, error) {
	// Get the path of the current go file (including the path inside
	// the project).
	var absolutePath string
	if _, filename, _, ok := runtime.Caller(0); ok {
		absolutePath = path.Dir(filename)
	}

	// Open the meta-schema file.
	file, err := os.Open(absolutePath + "/testdata/" + fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// Read the data from the file.
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
