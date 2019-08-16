package middlewares

import (
	"log"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/middleman"
)

// RouteLogger is a middleware that prints the path of any route hit.
func RouteLogger() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store map[string]string, end middleman.End) {
		log.Println(req.RequestURI)
	}
}
