package middlewares

import (
	"log"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/httputils"

	"github.com/Creespye/caf/internal/pkg/middleman"
)

// RouteLogger is a middleware that prints the path of any route hit.
func RouteLogger() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) {
		log.Println("[RouteLogger]: " + req.Method + " " + req.RequestURI)
	}
}

// BodyReader reads the body of a request as a []byte and stores it in the store argument under 'body'
func BodyReader() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) {

		// If a request has the Content-Length header, it has a body
		contentLength := httputils.GetContentLength(req.Header)

		// If the Content-Length is higher than zero, allocate an array for the
		// body data and read it
		if contentLength > 0 {
			body := make([]byte, contentLength, contentLength)

			req.Body.Read(body)

			req.Body.Close()

			// Store the body data in the store for future middlewares to use freely
			store.RequestBody = body
		}
	}
}
