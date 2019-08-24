package proxymiddlewares

import (
	"bytes"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/proxy"

	"github.com/Creespye/caf/internal/pkg/middleman"
)

// CreateTargetRequest creates a new request as a copy
// of the request from the client
func CreateTargetRequest(pr *proxy.Proxy, target string) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {

		bodyReader := bytes.NewReader(store["requestBody"].([]byte))

		tReq, err := pr.CreateTargetRequest(req.Method,
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
func SendTargetRequest(pr *proxy.Proxy) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {

		tRes, err :=
			pr.SendTargetRequest(store["targetRequest"].(*http.Request))

		store["targetResponse"] = tRes

		return err
	}
}

// ReadTargetResponseBody will read the target response body and store it in
// store.TargetResponseBody
func ReadTargetResponseBody(pr *proxy.Proxy) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		body, err :=
			pr.ReadTargetResponseBody(store["targetResponse"].(*http.Response))

		store["targetResponseBody"] = body

		return err
	}
}

// SendTargetResponse sends the target response to the client
func SendTargetResponse(pr *proxy.Proxy) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		err := pr.SendTargetResponse(res,
			store["targetResponse"].(*http.Response),
			store["targetResponseBody"].([]byte))

		return err
	}
}

// PrintRequestBody prints the request body
func PrintRequestBody(pr *proxy.Proxy) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		pr.PrintRequestBody(store["requestBody"].([]byte))

		return nil
	}
}

// PrintTargetResponseBody prints the request body
func PrintTargetResponseBody(pr *proxy.Proxy) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		pr.PrintTargetResponseBody(store["targetResponseBody"].([]byte))

		return nil
	}
}
