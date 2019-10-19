package formatchecker_test

import (
	"github.com/Creespye/caf/internal/pkg/validators/formatchecker"
	"testing"
)

type test struct {
	data  string
	valid bool
}

type format func(string) error

const succeed = "V"
const failed = "X"

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

	isValidFormat(t, testCases, FORMAT_DATE_TIME, formatchecker.IsValidDateTime)
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
	isValidFormat(t, testCases, FORMAT_DATE, formatchecker.IsValidDate)
}

func TestIsValidTime(t *testing.T) {
	testCases := []test{
		{
			data:  "08:30:06.283185Z",
			valid: true,
		},
		{
			data:  "10:05:08+01:00",
			valid: true,
		},
		{
			data:  "09:45:10 PST",
			valid: false,
		},
		{
			data:  "01:02:03,121212",
			valid: false,
		},
		{
			data:  "45:60:62",
			valid: false,
		},
		{
			data:  "1234",
			valid: false,
		},
	}
	isValidFormat(t, testCases, FORMAT_TIME, formatchecker.IsValidTime)
}

func TestIsValidEmail(t *testing.T) {
	testCases := []test{
		{
			data:  "john@example.com",
			valid: true,
		},
		{
			data:  "@",
			valid: false,
		},
		{
			data:  "john(at)example.com",
			valid: false,
		},
	}
	isValidFormat(t, testCases, FORMAT_EMAIL, formatchecker.IsValidEmail)
}

func TestIsValidIdnEmail(t *testing.T) {

}

func TestIsValidHostname(t *testing.T) {

}

func TestIsValidIdnHostname(t *testing.T) {

}

func TestIsValidIPv4(t *testing.T) {
}

func TestIsValidIPv6(t *testing.T) {

}

func TestIsValidURI(t *testing.T) {

}

func TestIsValidUriRef(t *testing.T) {

}

func TestIsValidIri(t *testing.T) {

}

func TestIsValidIriRef(t *testing.T) {

}

func TestIsValidURITemplate(t *testing.T) {

}

func TestIsValidJSONPointer(t *testing.T) {

}

func TestIsValidRelJSONPointer(t *testing.T) {

}

func TestIsValidRegex(t *testing.T) {

}

func isValidFormat(t *testing.T, tests []test, formatType string, fn format) {
	t.Logf("Given the need to test %s format", formatType)
	{
		for index, testCase := range tests {
			t.Logf("\tTest %d: When trying to format %s", index, testCase.data)
			{
				var valid bool
				if err := fn(testCase.data); err != nil {
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
