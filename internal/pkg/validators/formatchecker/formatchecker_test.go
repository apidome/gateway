package formatchecker_test

import (
	"github.com/Creespye/caf/internal/pkg/validators/formatchecker"
	"testing"
)

type test struct {
	data  string
	valid bool
}

const succeed = "V"
const failed = "X"

func TestIsValidDateTime(t *testing.T) {
	testCases := []test{
		{
			data:  "1985-04-12T23:20:50.52Z",
			valid: true,
		},
		{
			data:  "1996-12-19T16:39:57-08:00",
			valid: true,
		},
		{
			data:  "06/19/1963 08:30:06 PST",
			valid: false,
		},
	}

	t.Log("Given the need to test date-time format")
	{
		for index, testCase := range testCases {
			t.Logf("\tTest %d: When trying to format %s", index, testCase.data)
			{
				var valid bool
				if err := formatchecker.IsValidDateTime(testCase.data); err != nil {
					valid = false
				} else {
					valid = true
				}

				if valid != testCase.valid {
					t.Errorf("\t%s\tShould get valid = %t but got valid = %t, %s", failed, testCase.valid, valid)
				} else {
					t.Logf("\t%s\tvalid = %t", succeed, testCase.valid)
				}
			}
		}
	}
}

func TestIsValidDate(t *testing.T) {
	testCases := []test{
		{
			data:  "1963-06-19",
			valid: true,
		},
		{
			data:  "06/19/1963",
			valid: false,
		},
		{
			data:  "02-2002",
			valid: false,
		},
		{
			data:  "2010-350",
			valid: false,
		},
	}

	t.Log("Given the need to test date format")
	{
		for index, testCase := range testCases {
			t.Logf("\tTest %d: When trying to format %s", index, testCase.data)
			{
				var valid bool
				if err := formatchecker.IsValidDate(testCase.data); err != nil {
					valid = false
				} else {
					valid = true
				}

				if valid != testCase.valid {
					t.Errorf("\t%s\tShould get valid = %t but got valid = %t", failed, testCase.valid, valid)
				} else {
					t.Logf("\t%s\tvalid = %t", succeed, testCase.valid)
				}
			}
		}
	}
}
