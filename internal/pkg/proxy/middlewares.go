package proxy

import (
	"bytes"
	"errors"
	"io/ioutil"
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
			return errors.New("Request creation error:" + err.Error())
		}

		tReq.URL.Path = req.URL.Path
		tReq.URL.RawQuery = req.URL.RawQuery

		// Copy headers from the request to the target request
		httputils.CopyHeaders(req.Header, tReq.Header)

		// Create an http client to send the target request
		c := http.Client{}

		// Send the target request
		tRes, err := c.Do(tReq)

		if err != nil {
			return errors.New("Request send error:" + err.Error())
		}

		// Store the target response in the middleware store
		store.TargetResponse = tRes

		return nil
	}
}

// ReadResponseBody will read the response body and store it in
// store.TargetResponseBody
func ReadResponseBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {
		// Read the target response body
		targetResBody, err :=
			ioutil.ReadAll(store.TargetResponse.Body)

		if err != nil {
			return errors.New("Target response body read error: " + err.Error())
		}

		// Store it in the middleware store
		store.TargetResponseBody = targetResBody

		// Close the target response body
		err = store.TargetResponse.Body.Close()

		if err != nil {
			return errors.New("Body close error: " + err.Error())
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

		// Write target response body to response
		_, err := res.Write(store.TargetResponseBody)

		if err != nil {
			return errors.New("Send response error: " + err.Error())
		}

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
