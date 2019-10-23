package formatchecker_test

import (
	"github.com/Creespye/caf/internal/pkg/validators/formatchecker"
	"testing"
)

type test struct {
	data        string
	valid       bool
	description string
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
			description: "a valid date-time string",
			data:        "1985-04-12T23:20:50.52Z",
			valid:       true,
		},
		{
			description: "a valid date-time string",
			data:        "1996-12-19T16:39:57-08:00",
			valid:       true,
		},
		{
			description: "an invalid date-time string",
			data:        "06/19/1963 08:30:06 PST",
			valid:       false,
		},
	}

	isValidFormat(t, testCases, FORMAT_DATE_TIME, formatchecker.IsValidDateTime)
}

func TestIsValidDate(t *testing.T) {
	testCases := []test{
		{
			description: "a valid date string",
			data:        "1963-06-19",
			valid:       true,
		},
		{
			description: "an invalid date string (/ is invalid)",
			data:        "06/19/1963",
			valid:       false,
		},
		{
			description: "an invalid RFC3339 date",
			data:        "02-2002",
			valid:       false,
		},
		{
			description: "an invalid month 350",
			data:        "2010-350",
			valid:       false,
		},
	}
	isValidFormat(t, testCases, FORMAT_DATE, formatchecker.IsValidDate)
}

func TestIsValidTime(t *testing.T) {
	testCases := []test{
		{
			description: "a valid time",
			data:        "08:30:06.283185Z",
			valid:       true,
		},
		{
			description: "a valid time",
			data:        "10:05:08+01:00",
			valid:       true,
		},
		{
			description: "an invalid time",
			data:        "09:45:10 PST",
			valid:       false,
		},
		{
			description: "an invalid RFC3339 time",
			data:        "01:02:03,121212",
			valid:       false,
		},
		{
			description: "an invalid seconds",
			data:        "45:59:62",
			valid:       false,
		},
		{
			description: "an invalid time",
			data:        "1234",
			valid:       false,
		},
	}
	isValidFormat(t, testCases, FORMAT_TIME, formatchecker.IsValidTime)
}

func TestIsValidEmail(t *testing.T) {
	testCases := []test{
		{
			description: "a valid email",
			data:        "john@example.com",
			valid:       true,
		},
		{
			description: "an invalid email address",
			data:        "@",
			valid:       false,
		},
		{
			description: "@ is missing",
			data:        "john(at)example.com",
			valid:       false,
		},
		{
			description: "an invalid email address",
			data:        "1234",
			valid:       false,
		},
		{
			description: "an invalid email address",
			data:        "",
			valid:       false,
		},
	}
	isValidFormat(t, testCases, FORMAT_EMAIL, formatchecker.IsValidEmail)
}

func TestIsValidIdnEmail(t *testing.T) {
	testCases := []test{
		{
			description: "a valid idn email (example@example.test in Hangul)",
			data:        "실례@실례.테스트",
			valid:       true,
		},
		{
			description: "a valid idn email",
			data:        "john@example.com",
			valid:       true,
		},
		{
			description: "an invalid idn email",
			data:        "1234",
			valid:       false,
		},
		{
			description: "an invalid idn email",
			data:        "",
			valid:       false,
		},
	}
	isValidFormat(t, testCases, FORMAT_IDN_EMAIL, formatchecker.IsValidIdnEmail)
}

func TestIsValidHostname(t *testing.T) {
	testCases := []test{
		{
			description: "a valid host name",
			data:        "www.example.com",
			valid:       true,
		},
		{
			description: "a valid host name",
			data:        "xn--4gbwdl.xn--wgbh1c",
			valid:       true,
		},
		{
			description: "a host name containing illegal characters (_)",
			data:        "not_a_valid_host_name",
			valid:       false,
		},
		{
			description: "a host name starting with an illegal character",
			data:        "-a-host-name-that-starts-with--",
			valid:       false,
		},
		{
			description: "a host name with a component too long",
			data: "a-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-long-host-name-component",
			valid: false,
		},
	}
	isValidFormat(t, testCases, FORMAT_HOSTNAME, formatchecker.IsValidHostname)
}

func TestIsValidIdnHostname(t *testing.T) {
	testCases := []test{
		{
			description: "a valid host name (example.test in Hangul)",
			data:        "실례.테스트",
			valid:       true,
		},
		{
			description: "illegal first char",
			data:        "〮실례.테스트",
			valid:       false,
		},
		{
			description: "contains illegal",
			data:        "실〮례.테스트",
			valid:       false,
		},
	}
	isValidFormat(t, testCases, FORMAT_IDN_HOSTNAME, formatchecker.IsValidIdnHostname)
}

func TestIsValidIPv4(t *testing.T) {
	testCases := []test{
		{
			description: "a valid IPv4 address",
			data:        "192.168.0.1",
			valid:       true,
		},
		{
			description: "too many components",
			data:        "127.0.0.0.1",
			valid:       false,
		},
		{
			description: "IPv4 out of range",
			data:        "256.256.256.256",
			valid:       false,
		},
		{
			description: "not enough components (4 needed)",
			data:        "127",
			valid:       false,
		},
	}
	isValidFormat(t, testCases, FORMAT_IPV4, formatchecker.IsValidIPv4)
}

func TestIsValidIPv6(t *testing.T) {
	testCases := []test{
		{
			description: "a valid IPv6 address",
			data:        "::1",
			valid:       true,
		},
		{
			description: "IPv6 out of range",
			data:        "12345::",
			valid:       false,
		},
		{
			description: "too many components",
			data:        "1:1:1:1:1",
			valid:       false,
		},
		{
			description: "IPv6 containing illegal characters",
			data:        "::string",
			valid:       false,
		},
	}
	isValidFormat(t, testCases, FORMAT_IPV6, formatchecker.IsValidIPv6)
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
			t.Logf("\tTest %d: When trying to format %s => %s", index, testCase.data, testCase.description)
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
