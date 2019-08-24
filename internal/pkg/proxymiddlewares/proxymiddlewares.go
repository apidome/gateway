package proxymiddlewares

import (
	"bytes"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/middleman"
	"github.com/Creespye/caf/internal/pkg/proxy"
)

// CreateTargetRequest creates a new request as a copy
// of the request from the client
func CreateTargetRequest(target string) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {

		bodyReader := bytes.NewReader(store["requestBody"].([]byte))

		tReq, err := proxy.CreateTargetRequest(req.Method,
			target+req.RequestURI,
			req.URL.Path,
			req.URL.RawQuery,
			bodyReader,
			req.Header)

		store["targetRequest"] = tReq

		return err
	}
}

// SendTargetRequest forwards the target request to the target
// and stores the target response in store.TargetResponse
func SendTargetRequest() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {

		tRes, err :=
			proxy.SendTargetRequest(store["targetRequest"].(*http.Request))

		store["targetResponse"] = tRes

		return err
	}
}

// ReadTargetResponseBody will read the target response body and store it in
// store.TargetResponseBody
func ReadTargetResponseBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		body, err :=
			proxy.ReadTargetResponseBody(store["targetResponse"].(*http.Response))

		store["targetResponseBody"] = body

		return err
	}
}

// SendTargetResponse sends the target response to the client
func SendTargetResponse() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		err := proxy.SendTargetResponse(res,
			store["targetResponse"].(*http.Response),
			store["targetResponseBody"].([]byte))

		return err
	}
}

// PrintRequestBody prints the request body
func PrintRequestBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		proxy.PrintRequestBody(store["requestBody"].([]byte))

		return nil
	}
}

// PrintTargetResponseBody prints the request body
func PrintTargetResponseBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		proxy.PrintTargetResponseBody(store["targetResponseBody"].([]byte))

		return nil
	}
}
