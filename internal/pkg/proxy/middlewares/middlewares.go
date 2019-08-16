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
		store map[string]interface{}, end middleman.End) {

		parsedBody, _ := utils.ParseBody(store["body"])

		bodyReader := bytes.NewReader(parsedBody)

		tReq, err := http.NewRequest(req.Method, target, bodyReader)

		if err != nil {
			log.Println("[Request creation error]:", err.Error())
		}

		utils.CopyHeaders(req.Header, tReq.Header)

		c := http.Client{}

		tRes, err := c.Do(tReq)

		if err != nil {
			log.Println("[Request send error]:", err.Error())
		}

		utils.CopyHeaders(tRes.Header, res.Header())

		io.Copy(res, tRes.Body)
	}
}
