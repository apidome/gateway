package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/middleman"
	"github.com/Creespye/caf/internal/pkg/proxy/utils"
)

// ForwardRequest forwards the request to the target
func ForwardRequest(target string) middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) {

		// Create a reader from the body data, this requires the BodyReader middleware from middleman
		bodyReader := bytes.NewReader(store.Body)

		// Create a target request
		tReq, err := http.NewRequest(req.Method, target+req.RequestURI, bodyReader)

		if err != nil {
			log.Println("[Request creation error]:", err.Error())
		}

		// Copy headers from the request to the target request
		utils.CopyHeaders(req.Header, tReq.Header)

		// Create an http client to send the target request
		c := http.Client{}

		// Send the target request
		tRes, err := c.Do(tReq)

		if err != nil {
			log.Println("[Request send error]:", err.Error())
		}

		// Copy headers from target response
		utils.CopyHeaders(tRes.Header, res.Header())

		// Copy target response to response
		io.Copy(res, tRes.Body)
	}
}
