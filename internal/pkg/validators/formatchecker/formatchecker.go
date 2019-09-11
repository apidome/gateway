package formatchecker

import "time"

// from RFC 3339, section 5.6 [RFC3339]
// https://tools.ietf.org/html/rfc3339#section-5.6
func IsValidDateTime(dateTime string) (bool, error) {
	if _, err := time.Parse(time.RFC3339, dateTime); err != nil {
		return false, err
	}
	return true, nil
}
