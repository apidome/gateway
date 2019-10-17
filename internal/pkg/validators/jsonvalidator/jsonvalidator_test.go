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

}

func TestValidate(t *testing.T) {

}
