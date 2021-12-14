package http

import (
	"errors"
	"fmt"
	goSdk "go-sdk/backtest"
	"io"
	_http "net/http"
	"net/url"
	"strings"
	"time"
)

var Token string
var IsBacktest bool
var Backtest *goSdk.Backtest

type Request _http.Request
type Response _http.Response
type Header _http.Header

type Config struct {
	baseUrl        string
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
		baseUrl:        "https://api.tradologics.com/v1",
		socketUrl:      "tcp://0.0.0.0:3003",
		defaultTimeout: 5,
	}

	cln := _http.DefaultClient
	//cln.Timeout = time.Duration(cfg.defaultTimeout)

	return &Client{
		cfg: cfg,
		cln: cln,
	}
}

var DefaultClient = NewDefaultClient()

func (c *Client) Do(req *_http.Request) (*_http.Response, error) {
	if includeProtocol(req.URL.String()) {
		req, err := _http.NewRequest(req.Method, req.URL.String(), nil)

		if err != nil {
			return nil, err
		}
		return DefaultClient.cln.Do(req)
	}
	return c.cln.Do(req)
}

func (c *Client) preparedDo(method, url, contentType string, body io.Reader) (*_http.Response, error) {
	req, err := _http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if _, ok := req.Header["Authorization"]; !ok {
		if Token == "" && !IsBacktest {
			return nil, errors.New("please use `SetToken(...)` first")
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	}

	if IsBacktest {
		res, err := Backtest.CallErocMethod(req)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	urlPath := req.URL.Path
	if strings.HasSuffix(urlPath, "/") {
		urlPath = urlPath[:len(urlPath)-1]
	}
	return DefaultClient.cln.Do(req)
}

func (c *Client) Head(url string) (resp *_http.Response, err error) {
	if includeProtocol(url) {
		return DefaultClient.cln.Head(url)
	}
	return c.preparedDo("HEAD", url, "", nil)
}

func (c *Client) Get(url string) (resp *_http.Response, err error) {
	if includeProtocol(url) {
		return DefaultClient.cln.Get(url)
	}
	return c.preparedDo("GET", url, "", nil)
}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *_http.Response, err error) {
	if includeProtocol(url) {
		return DefaultClient.cln.Post(url, contentType, body)
	}
	return c.preparedDo("POST", url, contentType, nil)
}

func (c *Client) PostForm(url string, data url.Values) (resp *_http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func SetToken(token string) {
	Token = token
}

func SetBacktestMode(start, end string) {
	Backtest = goSdk.NewBacktest(start, end, DefaultClient.cfg.socketUrl)
	IsBacktest = true
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
