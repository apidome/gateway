package utils

import (
	"net/http"
)

// CopyHeaders copies headers from a Header object to another Header object
func CopyHeaders(src, dst http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
