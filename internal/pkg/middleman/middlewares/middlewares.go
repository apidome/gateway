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
		store map[string]interface{}, end middleman.End) {
		log.Println("[RouteLogger]: " + req.RequestURI)
	}
}

// BodyReader reads the body of a request as a []byte and stores it in the store argument under 'body'
func BodyReader() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store map[string]interface{}, end middleman.End) {

		contentLength := req.Header.Get("Content-Length")

		if contentLength != "" {
			contentLengthNum, err := strconv.Atoi(req.Header.Get("Content-Length"))

			if err != nil {
				log.Println("[Content-Length convert error]:", err.Error())
			} else {
				if contentLengthNum > 0 {
					body := make([]byte, contentLengthNum, contentLengthNum)

					req.Body.Read(body)

					req.Body.Close()

					store["body"] = body
				}
			}
		}
	}
}
