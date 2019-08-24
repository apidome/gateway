package proxy

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Creespye/caf/internal/pkg/httputils"
)

// Proxy is a struct that holds the proxy information
type Proxy struct {
	Client http.Client
}

// NewProxy creates a new default proxy
func NewProxy() Proxy {
	return Proxy{
		Client: http.Client{},
	}
}

// CreateTargetRequest creates a new request as a copy
// of the request from the client
func (pr *Proxy) CreateTargetRequest(method, target, path, query string,
	body io.Reader, header http.Header) (*http.Request, error) {
	// Create a target request
	tReq, err := http.NewRequest(method,
		target,
		body)

	if err != nil {
		return nil, errors.New("Target request creation error:" + err.Error())
	}

	tReq.URL.Path = path
	tReq.URL.RawQuery = query

	// Copy headers from the request to the target request
	httputils.CopyHeaders(header, tReq.Header)

	return tReq, nil
}

// SendTargetRequest forwards the target request to the target
// and returnes the response
func (pr *Proxy) SendTargetRequest(req *http.Request) (*http.Response, error) {
	// Send the target request
	tRes, err := pr.Client.Do(req)

	if err != nil {
		return nil, errors.New("Target request send error:" + err.Error())
	}

	return tRes, nil
}

// ReadTargetResponseBody will read the target response body and return it
func (pr *Proxy) ReadTargetResponseBody(tRes *http.Response) ([]byte, error) {
	// Read the target response body
	targetResBody, err :=
		ioutil.ReadAll(tRes.Body)

	if err != nil {
		errMsg := "Target response body read error: " + err.Error()
		return nil, errors.New(errMsg)
	}

	// Close the target response body
	err = tRes.Body.Close()

	if err != nil {
		errMsg := "Target response body close error: " + err.Error()
		return nil, errors.New(errMsg)
	}

	return targetResBody, nil
}

// SendTargetResponse sends the target response to the client
func (pr *Proxy) SendTargetResponse(res http.ResponseWriter,
	targetRes *http.Response, body []byte) error {
	// Copy headers from target response
	httputils.CopyHeaders(targetRes.Header, res.Header())

	// Write the header of the target response
	res.WriteHeader(targetRes.StatusCode)

	// Write target response body to response
	_, err := res.Write(body)

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

// PrintRequestBody prints the request body
func (pr *Proxy) PrintRequestBody(body []byte) {
	log.Println("[RequestBody]:", body)
}

// PrintTargetResponseBody prints the request body
func (pr *Proxy) PrintTargetResponseBody(body []byte) {
	log.Println("[TargetResponseBody]:", body)
}
