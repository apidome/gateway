package httputils

import (
	"log"
	"net/http"
	"strconv"
)

// CopyHeaders copies headers from a Header object to another Header object
func CopyHeaders(src, dst http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

// GetContentLength returns the content length of a response by its header
func GetContentLength(header http.Header) int {
	contentLength := header.Get("Content-Length")

	if contentLength != "" {
		contentLengthNum, err := strconv.Atoi(contentLength)

		if err != nil {
			log.Println("[ContentLength conversion error]:", err.Error())
		} else {
			return contentLengthNum
		}
	}

	return 0
}
