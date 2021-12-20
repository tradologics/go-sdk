package backtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const DefaultErrorMessage = "Something bad happen"

type ErocRequestHeader struct {
	Start      string `json:"start"`
	End        string `json:"end"`
	Datetime   string `json:"datetime"`
	Resolution string `json:"resolution"`
}

type ErocRequestData map[string]interface{}

type ErocRequest struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Data    ErocRequestData   `json:"data"`
	Headers ErocRequestHeader `json:"headers"`
}

type ErocError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type RuntimeEvents map[string]interface{}

type ErocResponse struct {
	Status int           `json:"status"`
	Errors []ErocError   `json:"errors"`
	Data   interface{}   `json:"data"`
	Events RuntimeEvents `json:"events"`
}

type BacktestResponse struct {
	Errors []ErocError `json:"errors"`
	Data   interface{} `json:"data"`
}

type BarInfo struct {
	Datetime   string
	Resolution string
}

type Backtest struct {
	start          string
	end            string
	currentBarInfo *BarInfo
	runtimeEvents  RuntimeEvents
	zmqConn        *ZmqConn
}

// NewBacktest create new Backtest object with selected start,
// end dates and create new ZMQ connection using chosen socket URL
func NewBacktest(start, end, socketUrl string) (*Backtest, error) {

	// Create new ZMQ connection
	zmqConn, err := NewZmq(socketUrl)
	if err != nil {
		return nil, err
	}

	return &Backtest{
		start:          start,
		end:            end,
		currentBarInfo: &BarInfo{},
		zmqConn:        zmqConn,
	}, nil
}

// CallErocMethod parse client request data and use it to create new EROC request and send data using ZMQ;
// Returns EROC response as HTTP response.
func (b *Backtest) CallErocMethod(req *http.Request) *http.Response {

	// Parse request data as JSON to erocRequestData structure
	erocRequestData := ErocRequestData{}

	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return b.errorHandler(req, err, DefaultErrorMessage)
		}

		err = json.Unmarshal(body, &erocRequestData)
		if err != nil {
			return b.errorHandler(req, err, "Invalid JSON")
		}
	}

	erocRequest := &ErocRequest{
		Method: req.Method,
		Url:    req.URL.String(),
		Data:   erocRequestData,
		Headers: ErocRequestHeader{
			Start:      b.start,
			End:        b.end,
			Datetime:   b.currentBarInfo.Datetime,
			Resolution: b.currentBarInfo.Resolution,
		},
	}

	err := b.zmqConn.SendJSON(&erocRequest)
	if err != nil {
		return b.errorHandler(req, err, DefaultErrorMessage)
	}

	var erocResponse ErocResponse
	err = b.zmqConn.ReceiveJSON(&erocResponse)
	if err != nil {
		return b.errorHandler(req, err, DefaultErrorMessage)
	}

	// Set runtime events
	b.runtimeEvents = erocResponse.Events

	erocJSONResponse, err := json.Marshal(BacktestResponse{
		Errors: erocResponse.Errors,
		Data:   erocResponse.Data,
	})
	if err != nil {
		return b.errorHandler(req, err, DefaultErrorMessage)
	}

	res := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(erocJSONResponse)),

		StatusCode: erocResponse.Status,
		Status:     fmt.Sprintf("%d %s", erocResponse.Status, http.StatusText(erocResponse.Status)),

		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,

		Request: req,
	}
	req.Response = res

	return res

}

// errorHandler returns HTTP Bad Gateway error if something unexpected happened inside CallErocMethod function
func (b *Backtest) errorHandler(req *http.Request, err error, message string) *http.Response {

	// Log source error
	if err != nil {
		log.Println(err)
	}

	erocJSONResponse, err := json.Marshal(ErocResponse{
		Status: 502,
		Errors: []ErocError{{ID: "internal_server_error", Message: message}},
		Data:   make(map[string]interface{}),
	})
	if err != nil {
		log.Fatalln(err)
	}

	res := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(erocJSONResponse)),

		StatusCode: http.StatusBadGateway,
		Status:     fmt.Sprintf("%d %s", http.StatusBadGateway, http.StatusText(http.StatusBadGateway)),

		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,

		Request: req,
	}
	req.Response = res

	return res
}

// SetCurrentBarInfo set currentBarInfo datetime and resolution
func (b *Backtest) SetCurrentBarInfo(info *BarInfo) {
	b.currentBarInfo = info
}

// GetRuntimeEvents returns current Backtest events data
func (b *Backtest) GetRuntimeEvents() map[string]interface{} {
	return b.runtimeEvents
}

// Close ZMQ connection
func (b *Backtest) Close() {
	b.zmqConn.Close()
}
