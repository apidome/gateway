package proxy

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Creespye/caf/internal/pkg/httputils"
	"github.com/Creespye/caf/internal/pkg/middleman"
)

// CreateTargetRequests creates a new request as a copy
// of the request from the client
func CreateTargetRequest(target string) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {
		bodyReader := bytes.NewReader(store.RequestBody)

		// Create a target request
		tReq, err := http.NewRequest(req.Method,
			target+req.RequestURI,
			bodyReader)

		if err != nil {
			return errors.New("Target request creation error:" + err.Error())
		}

		tReq.URL.Path = req.URL.Path
		tReq.URL.RawQuery = req.URL.RawQuery

		// Copy headers from the request to the target request
		httputils.CopyHeaders(req.Header, tReq.Header)

		// Store the target request for all middlewares to use
		store.TargetRequest = tReq

		return nil
	}
}

// SendTargetRequest forwards the target request to the target
// and stores the target response in store.TargetResponse
func SendTargetRequest() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {

		// Create an http client to send the target request
		c := http.Client{}

		// Send the target request
		tRes, err := c.Do(store.TargetRequest)

		if err != nil {
			return errors.New("Target request send error:" + err.Error())
		}

		// Store the target response in the middleware store
		store.TargetResponse = tRes

		return nil
	}
}

// ReadTargetResponseBody will read the target response body and store it in
// store.TargetResponseBody
func ReadTargetResponseBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {
		// Read the target response body
		targetResBody, err :=
			ioutil.ReadAll(store.TargetResponse.Body)

		if err != nil {
			errMsg := "Target response body read error: " + err.Error()
			return errors.New(errMsg)
		}

		// Store it in the middleware store
		store.TargetResponseBody = targetResBody

		// Close the target response body
		err = store.TargetResponse.Body.Close()

		if err != nil {
			errMsg := "Target response body close error: " + err.Error()
			return errors.New(errMsg)
		}

		return nil
	}
}

// SendTargetResponse sends the target response to the client
func SendTargetResponse() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {

		// Copy headers from target response
		httputils.CopyHeaders(store.TargetResponse.Header, res.Header())

		// Write the header of the target response
		res.WriteHeader(store.TargetResponse.StatusCode)

		// Write target response body to response
		_, err := res.Write(store.TargetResponseBody)

		// If the response status code does not support body it will not
		// be written and can be ignored
		if err != nil {
			if !strings.HasSuffix(err.Error(),
				"request method or response status code does not allow body") {
				return errors.New("Send response error: " + err.Error())
			}
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
