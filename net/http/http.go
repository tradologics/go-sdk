package http

import (
	"errors"
	"fmt"
	goSdk "go-sdk/backtest"
	"io"
	"log"
	_http "net/http"
	"net/url"
	"strings"
	"time"
)

var Token string
var IsBacktest bool
var Backtest *goSdk.Backtest
var NewRequest = _http.NewRequest
var NewRequestWithContext = _http.NewRequestWithContext

type Request _http.Request
type Response _http.Response
type Header _http.Header

type Config struct {
	baseSchema     string
	basePath       string
	baseHost       string
	defaultTimeout int
	socketUrl      string
}

type Client struct {
	cln *_http.Client
	cfg *Config

	Transport     _http.RoundTripper
	Jar           _http.CookieJar
	Timeout       time.Duration
	CheckRedirect func(req *Request, via []*Request) error
}

func NewDefaultClient() *Client {
	cfg := &Config{
		baseHost:       "api.tradologics.com",
		baseSchema:     "https",
		basePath:       "/v1",
		socketUrl:      "tcp://0.0.0.0:3003",
		defaultTimeout: 5,
	}

	cln := _http.DefaultClient
	// TODO timeout
	//cln.Timeout = time.Duration(cfg.defaultTimeout)

	return &Client{
		cfg: cfg,
		cln: cln,
	}
}

var DefaultClient = NewDefaultClient()

func (c *Client) isConfigured() {
	if c.cln == nil {
		// TODO default client
	}
}

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
	c.isConfigured()

	// If url include protocol, then use default client
	if (req.URL != nil && req.URL.Scheme != "") || req.Host != "" && !includeProtocol(req.Host) {
		return c.cln.Do(req)
	}
	return c.processRequest(req.Method, req.URL.Path, "", req.Body, req.Header)

}

// TODO Is default exists

func (c *Client) processRequest(method, url, contentType string, body io.Reader, header _http.Header) (*_http.Response, error) {
	c.isConfigured()

	if includeProtocol(url) {
		req, err := newRequestWithContentType(method, url, contentType, body)
		if err != nil {
			return nil, err
		}
		req.Header = header

		return c.cln.Do(req)
	}

	if IsBacktest {
		// TODO handler error
		req, err := newRequestWithContentType(method, url, contentType, body)
		if err != nil {
			return nil, err
		}

		res, err := Backtest.CallErocMethod(req)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	fullUrl := fmt.Sprintf("%s://%s%s%s", c.cfg.baseSchema, c.cfg.baseHost, c.cfg.basePath, url)
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

	// TODO host and url FIX IT
	if strings.HasSuffix(req.URL.Path, "/") {
		req.URL.Path = req.URL.Path[:len(req.URL.Path)-1]
	}

	r, err := DefaultClient.cln.Do(req)
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
	Backtest, err = goSdk.NewBacktest(start, end, DefaultClient.cfg.socketUrl)
	if err != nil {
		return err
	}

	IsBacktest = true

	return nil
}

func SetCurrentBarInfo(info *goSdk.BarInfo) error {
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
