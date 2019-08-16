package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
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

// ParseBody receives the interface type body and converts it to a []byte
func ParseBody(body interface{}) ([]byte, error) {
	if body != nil {
		var bodyBuf bytes.Buffer
		enc := gob.NewEncoder(&bodyBuf)
		err := enc.Encode(body)

		if err != nil {
			log.Println("[Body conversion error]:", err.Error())
		}

		return bodyBuf.Bytes(), nil
	} else {
		return nil, errors.New("Empty body")
	}
}
