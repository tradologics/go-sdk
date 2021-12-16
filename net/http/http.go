package http

import (
	"errors"
	"fmt"
	"github.com/tradologics/go-sdk/backtest"
	"io"
	"log"
	_http "net/http"
	"net/url"
	"strings"
)

const (
	baseHost       = "api.tradologics.com"
	baseSchema     = "https"
	basePath       = "/v1"
	socketUrl      = "tcp://0.0.0.0:3003"
	defaultTimeout = 5
)

var Token string
var IsBacktest bool
var Backtest *backtest.Backtest
var NewRequest = _http.NewRequest
var NewRequestWithContext = _http.NewRequestWithContext

type Request _http.Request
type Response _http.Response
type Header _http.Header
type Client _http.Client

func NewDefaultClient() *Client {
	return &Client{Timeout: defaultTimeout}
}

var DefaultClient = NewDefaultClient()
var httpDefaultClient = _http.DefaultClient

func newRequestWithContentType(method, url string, contentType string, body io.Reader) (*_http.Request, error) {
	req, err := NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, err
}

func (c *Client) Do(req *_http.Request) (*_http.Response, error) {

	// If url include protocol, then use default http client
	if (req.URL != nil && req.URL.Scheme != "") || req.Host != "" && !includeProtocol(req.Host) {
		return httpDefaultClient.Do(req)
	}
	return c.processRequest(req.Method, req.URL.Path, "", req.Body, req.Header)

}

func (c *Client) processRequest(method, url, contentType string, body io.Reader, header _http.Header) (*_http.Response, error) {

	if includeProtocol(url) {
		req, err := newRequestWithContentType(method, url, contentType, body)
		if err != nil {
			return nil, err
		}
		req.Header = header

		return httpDefaultClient.Do(req)
	}

	if IsBacktest {
		// TODO error handler
		req, err := newRequestWithContentType(method, url, contentType, body)
		if err != nil {
			return nil, err
		}

		res := Backtest.CallErocMethod(req)

		return res, nil
	}

	fullUrl := fmt.Sprintf("%s://%s%s%s", baseSchema, baseHost, basePath, url)
	req, err := newRequestWithContentType(method, fullUrl, contentType, body)
	if err != nil {
		return nil, err
	}

	if _, ok := req.Header["Authorization"]; !ok {
		if Token == "" && !IsBacktest {
			return nil, errors.New("please use `SetToken(...)` first")
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	}

	if strings.HasSuffix(req.URL.Path, "/") {
		req.URL.Path = req.URL.Path[:len(req.URL.Path)-1]
	}

	r, err := httpDefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return r, nil
}

func Head(url string) (resp *_http.Response, err error) {
	return DefaultClient.Head(url)
}

func (c *Client) Head(url string) (resp *_http.Response, err error) {
	return c.processRequest("HEAD", url, "", nil, nil)
}

func Get(url string) (resp *_http.Response, err error) {
	return DefaultClient.Get(url)
}

func (c *Client) Get(url string) (resp *_http.Response, err error) {
	return c.processRequest("GET", url, "", nil, nil)
}

func Post(url, contentType string, body io.Reader) (resp *_http.Response, err error) {
	return DefaultClient.Post(url, contentType, body)
}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *_http.Response, err error) {
	return c.processRequest("POST", url, contentType, body, nil)
}

func PostForm(url string, data url.Values) (resp *_http.Response, err error) {
	return DefaultClient.PostForm(url, data)
}

func (c *Client) PostForm(url string, data url.Values) (resp *_http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func SetToken(token string) {
	Token = token
}

func SetBacktestMode(start, end string) (err error) {
	Backtest, err = backtest.NewBacktest(start, end, socketUrl)
	if err != nil {
		return err
	}

	IsBacktest = true

	return nil
}

func SetCurrentBarInfo(info *backtest.BarInfo) error {
	if Backtest != nil {
		Backtest.SetCurrentBarInfo(info)

		return nil
	}
	return errors.New("please set backtest mode first")
}

func GetRuntimeEvents() (interface{}, error) {
	if Backtest != nil {
		return Backtest.GetRuntimeEvents(), nil
	}
	return nil, errors.New("please set backtest mode first")
}
