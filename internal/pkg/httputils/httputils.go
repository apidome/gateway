package httputils

import (
	"io/ioutil"
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

// ReadRequestBody reads the requests body, closes the reader and returns
// the request body raw data
func ReadRequestBody(req *http.Request) ([]byte, error) {

	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	err = req.Body.Close()

	if err != nil {
		return nil, err
	}

	return reqBody, nil
}

// ReadResponseBody reads the response body, closes the reader and returns
// the response body raw data
func ReadResponseBody(res *http.Response) ([]byte, error) {

	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	err = res.Body.Close()

	if err != nil {
		return nil, err
	}

	return resBody, nil
}
