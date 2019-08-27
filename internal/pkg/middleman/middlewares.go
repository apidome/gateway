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
		store Store, end End) error {
		log.Println("[RouteLogger]: " + req.Method + " " + req.RequestURI)

		return nil
	}
}

// BodyReader reads the body of a request as a []byte and stores it
// in store["body"]
func BodyReader() Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store Store, end End) error {

		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			return errors.New("Request body read error: " + err.Error())
		}

		req.Body.Close()

		store["requestBody"] = body

		return nil
	}
}
