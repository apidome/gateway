package proxy

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/apidome/gateway/internal/pkg/httputils"
)

var (
	// ErrHijackingNotOk is returned when hijacking a request is not supported
	ErrHijackingNotOk = errors.New("Hijacking not supported")
)

// Proxy is a struct that holds the proxy information
type Proxy struct {
	target string
	Client http.Client
}

// InitProxy initializes a Proxy instance
func InitProxy(pr *Proxy, target string) {
	pr.Client = http.Client{}
	pr.target = target
}

// NewProxy creates a new proxy
func NewProxy(target string) *Proxy {
	pr := &Proxy{}

	InitProxy(pr, target)

	return pr
}

// CreateRequest creates a new request
func (pr *Proxy) CreateRequest(method,
	path,
	rawQuery string,
	headers http.Header,
	body []byte) (*http.Request, error) {

	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest(method,
		pr.target,
		bodyReader)

	if err != nil {
		return nil, err
	}

	req.URL.Path = path
	req.URL.RawQuery = rawQuery

	httputils.CopyHeaders(headers, req.Header)

	return req, nil
}

// SendRequest forwards the target request to the target
// and returnes the response
func (pr *Proxy) SendRequest(req *http.Request) (*http.Response, error) {
	// Send the target request
	res, err := pr.Client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CopyResponseToClient sends the target response to the client
func CopyResponseToClient(res http.ResponseWriter,
	targetRes *http.Response,
	body []byte) error {

	httputils.CopyHeaders(targetRes.Header, res.Header())

	res.WriteHeader(targetRes.StatusCode)

	_, err := res.Write(body)

	// If the response status code does not support body it will not
	// be written and can be ignored
	if err != nil {
		if err != http.ErrBodyNotAllowed {
			return err
		}
	}

	return nil
}

// TunnelConnection tunnels a connection to the target server
func TunnelConnection(res http.ResponseWriter,
	req *http.Request,
	target string) error {
	destConn, err := net.DialTimeout("tcp", target, 10*time.Second)

	if err != nil {
		return err
	}

	res.WriteHeader(http.StatusOK)

	hijacker, ok := res.(http.Hijacker)

	if !ok {
		return ErrHijackingNotOk
	}

	clientCon, _, err := hijacker.Hijack()

	if err != nil {
		return err
	}

	go connectPipes(destConn, clientCon)
	go connectPipes(clientCon, destConn)

	return nil
}

// connectPipes copies data from a socket to another
func connectPipes(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()

	io.Copy(dst, src)
}
