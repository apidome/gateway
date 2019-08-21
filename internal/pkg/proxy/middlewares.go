package proxy

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/httputils"
	"github.com/Creespye/caf/internal/pkg/middleman"
)

// SendRequest forwards the request to the target
func SendRequest(target string) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {

		// Create a reader from the body data, this requires the
		// BodyReader middleware from middleman
		bodyReader := bytes.NewReader(store.RequestBody)

		// Create a target request
		tReq, err := http.NewRequest(req.Method,
			target+req.RequestURI,
			bodyReader)

		if err != nil {
			return errors.New("[Request creation error]:" + err.Error())
		}

		// Copy headers from the request to the target request
		httputils.CopyHeaders(req.Header, tReq.Header)

		// Create an http client to send the target request
		c := http.Client{}

		// Send the target request
		tRes, err := c.Do(tReq)

		if err != nil {
			return errors.New("[Request send error]:" + err.Error())
		}

		// Store the target response in the middleware store
		store.TargetResponse = tRes

		// Read the content length header from the target response
		contentLength := httputils.GetContentLength(tRes.Header)

		// If the content length header exists,
		// read the body of the target response
		if contentLength > 0 {
			// Create a buffer to read the body
			body := make([]byte, contentLength, contentLength)

			tRes.Body.Read(body)

			tRes.Body.Close()

			// Store the target response body in the middleware store
			store.TargetResponseBody = body

		}

		return nil
	}
}

// SendResponse sends the target response to the client
func SendResponse() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {

		// Copy headers from target response
		httputils.CopyHeaders(store.TargetResponse.Header, res.Header())

		// Create a reader to read the target response body from
		targetResBody := bytes.NewReader(store.TargetResponseBody)

		// Copy target response to response
		io.Copy(res, targetResBody)

		return nil
	}
}

// PrintRequestBody prints the request body
func PrintRequestBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {
		log.Println("[RequestBody]:", store.RequestBody)

		return nil
	}
}

// PrintTargetResponseBody prints the request body
func PrintTargetResponseBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {
		log.Println("[TargetResponseBody]:", store.TargetResponseBody)

		return nil
	}
}
