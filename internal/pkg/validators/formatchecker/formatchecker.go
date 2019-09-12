package formatchecker

import (
	"fmt"
	"net/mail"
	"time"
)

// from RFC 3339, section 5.6 [RFC3339]
// https://tools.ietf.org/html/rfc3339#section-5.6
func IsValidDateTime(dateTime string) (bool, error) {
	if _, err := time.Parse(time.RFC3339, dateTime); err != nil {
		return false, err
	}
	return true, nil
}

// RFC 3339, section 5.6 [RFC3339]
// https://tools.ietf.org/html/rfc3339#section-5.6
func IsValidDate(date string) (bool, error) {
	timeToAppend := "T00:00:00.0Z"
	dateTime := fmt.Sprintf("%s%s", date, timeToAppend)
	return IsValidDateTime(dateTime)
}

// RFC 3339, section 5.6 [RFC3339]
// https://tools.ietf.org/html/rfc3339#section-5.6
func IsValidTime(time string) (bool, error) {
	dateToAppend := "1991-02-21"
	dateTime := fmt.Sprintf("%sT%s", dateToAppend, time)
	return IsValidDateTime(dateTime)
}

// RFC 5322, section 3.4.1 [RFC5322].
// https://tools.ietf.org/html/rfc5322#section-3.4.1
func IsValidEmail(email string) (bool, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return false, err
	}
	return true, nil
}
