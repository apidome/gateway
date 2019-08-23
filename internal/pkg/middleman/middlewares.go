package middleman

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// RouteLogger is a middleware that prints the path of any route hit
func RouteLogger() Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *Store, end End) error {
		log.Println("[RouteLogger]: " + req.Method + " " + req.RequestURI)

		return nil
	}
}

// BodyReader reads the body of a request as a []byte and stores it
// in the store argument under 'body'
func BodyReader() Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store *Store, end End) error {
		// Read the body of the request
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			return errors.New("Request body read error: " + err.Error())
		}

		// Close the request body
		req.Body.Close()

		// Store the body data in the store for future
		// middlewares to use freely
		store.RequestBody = body

		return nil
	}
}
