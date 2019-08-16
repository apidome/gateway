package proxy

import (
	"bytes"
	"log"
	"net/http"
	"strconv"

	"github.com/Creespye/caf/internal/pkg/middleman"
	"github.com/Creespye/caf/internal/pkg/middleman/middlewares"
	"github.com/Creespye/caf/internal/pkg/proxy/utils"
)

// Config is a struct that holds all configurations of the proxy server
type Config struct {
	Addr   string
	Target string
	Cert   string
	Key    string
}

// Start starts the proxy server and begins operating on requests
func Start(config Config) {
	// Creating a new middleman (middleware manager)
	mm := middleman.NewMiddleman(middleman.Config{
		Addr:     config.Addr,
		Target:   config.Target,
		CertFile: config.Cert,
		KeyFile:  config.Key,
	})

	mm.Use(middlewares.RouteLogger())

	// Create a new middleware for the root route
	mm.Use(func(res http.ResponseWriter, req *http.Request,
		store map[string]string, end middleman.End) {

		var targetReq *http.Request
		var reqError error

		switch req.Method {
		case http.MethodPost:
			// Store request body
			body := make([]byte, 255)

			_, err := req.Body.Read(body)

			if err != nil {
				if err.Error() != "EOF" {
					log.Println("[Reading request body error]:", err.Error())
				}
			}

			bodyReader := bytes.NewReader(body)

			targetReq, reqError = http.NewRequest(req.Method, config.Target+req.RequestURI, bodyReader)
		case http.MethodGet:
			targetReq, reqError = http.NewRequest(req.Method, config.Target+req.RequestURI, nil)
		}

		if reqError != nil {
			log.Println("[Target request creation error]: ", reqError.Error())
		}

		utils.CopyHeaders(req.Header, targetReq.Header)

		httpClient := http.Client{}

		targetRes, err := httpClient.Do(targetReq)

		if err != nil {
			log.Println("[Target request send error]:", err.Error())
		}

		utils.CopyHeaders(targetRes.Header, res.Header())

		contentLength, err := strconv.Atoi(targetRes.Header.Get("Content-Length"))

		if err != nil {
			log.Println("[Target response content length fetch error]:", err.Error())
		}

		targetResBody := make([]byte, contentLength, contentLength)

		_, err = targetRes.Body.Read(targetResBody)

		targetRes.Body.Close()

		if err != nil {
			if err.Error() != "EOF" {
				log.Println("[Target response read error]:", err.Error())
			}
		}

		res.Write(targetResBody)
	})

	// Begin listening
	err := mm.ListenAndServeTLS(func() {
		log.Println("[Middleman is listening on]:", config.Addr)
	})

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("[Failed creating a server]:", err)
	}
}
