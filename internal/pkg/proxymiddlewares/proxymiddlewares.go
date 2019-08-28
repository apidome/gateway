package proxymiddlewares

import (
	"log"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/httputils"

	"github.com/Creespye/caf/internal/pkg/proxy"

	"github.com/Creespye/caf/internal/pkg/middleman"
)

// CreateRequest creates a new request as a copy
// of the request from the client
func CreateRequest(pr *proxy.Proxy) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {

		tReq, err := pr.CreateRequest(req.Method,
			req.URL.Path,
			req.URL.RawQuery,
			req.Header,
			store["requestBody"].([]byte))

		store["targetRequest"] = tReq

		return err
	}
}

// SendRequest forwards the target request to the target
// and stores the target response in store.TargetResponse
func SendRequest(pr *proxy.Proxy) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {

		tRes, err :=
			pr.SendRequest(store["targetRequest"].(*http.Request))

		store["targetResponse"] = tRes

		return err
	}
}

// ReadResponseBody will read the target response body and store it in
// store.TargetResponseBody
func ReadResponseBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		body, err :=
			httputils.ReadResponseBody(store["targetResponse"].(*http.Response))

		store["targetResponseBody"] = body

		return err
	}
}

// SendResponse sends the target response to the client
func SendResponse() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		err := proxy.CopyResponseToClient(res,
			store["targetResponse"].(*http.Response),
			store["targetResponseBody"].([]byte))

		return err
	}
}

// PrintRequestBody prints the request body
func PrintRequestBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		log.Println(store["requestBody"].([]byte))

		return nil
	}
}

// PrintTargetResponseBody prints the request body
func PrintTargetResponseBody() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		log.Println(store["targetResponseBody"].([]byte))

		return nil
	}
}

// ValidateRequest is a middleware that handles validation of an HTTP request.
func ValidateRequest(validator *validators.Validator) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {
		//validator.Validate()
		return nil
	}
}
