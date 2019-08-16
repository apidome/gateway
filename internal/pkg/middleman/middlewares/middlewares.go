package middlewares

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Creespye/caf/internal/pkg/middleman"
)

// RouteLogger is a middleware that prints the path of any route hit.
func RouteLogger() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) {
		log.Println("[RouteLogger]: " + req.RequestURI)
	}
}

// BodyReader reads the body of a request as a []byte and stores it in the store argument under 'body'
func BodyReader() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) {

		// If a request has the Content-Length header, it has a body
		contentLength := req.Header.Get("Content-Length")

		// If the Content-Length header exist, convert it to an integer
		if contentLength != "" {
			contentLengthNum, err := strconv.Atoi(req.Header.Get("Content-Length"))

			if err != nil {
				log.Println("[Content-Length convert error]:", err.Error())
			} else {
				// If the Content-Length is higher than zero, allocate an array for the
				// body data and read it
				if contentLengthNum > 0 {
					body := make([]byte, contentLengthNum, contentLengthNum)

					req.Body.Read(body)

					req.Body.Close()

					// Store the body data in the store for future middlewares to use freely
					store.Body = body
				}
			}
		}
	}
}
