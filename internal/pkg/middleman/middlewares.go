package middleman

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

// VariablesReader reads the variables from the request path
// and stores them in store["variables"]
func VariablesReader() Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store Store, end End) error {
		variables := strings.Split(req.URL.Path, "/")[1:]

		if variables[0] != "" {
			store["variables"] = variables
		} else {
			store["variables"] = []string{}
		}

		return nil
	}
}

// ParametersReader reads the query parameters from the request
// and stores them in store["parameters"]
func ParametersReader() Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store Store, end End) error {
		params := req.URL.Query()
		parameters := map[string]string{}

		for param := range params {
			parameters[param] = params.Get(param)
		}

		store["parameters"] = parameters

		return nil
	}
}
