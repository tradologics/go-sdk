package backtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

func NewBacktest(start, end, socketUrl string) (*Backtest, error) {
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

func (b *Backtest) CallErocMethod(req *http.Request) *http.Response {
	erocRequestData := ErocRequestData{}

	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return b.errorHandler(req, err)
		}

		err = json.Unmarshal(body, &erocRequestData)
		if err != nil {
			return b.errorHandler(req, err)
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
		return b.errorHandler(req, err)
	}

	var erocResponse ErocResponse
	err = b.zmqConn.ReceiveJSON(&erocResponse)
	if err != nil {
		return b.errorHandler(req, err)
	}

	b.runtimeEvents = erocResponse.Events

	erocJSONResponse, err := json.Marshal(BacktestResponse{
		Errors: erocResponse.Errors,
		Data:   erocResponse.Data,
	})
	if err != nil {
		return b.errorHandler(req, err)
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

func (b *Backtest) errorHandler(req *http.Request, err error) *http.Response {
	// TODO find correct error message
	erocJSONResponse, err := json.Marshal(ErocResponse{
		Status: 502,
		Errors: []ErocError{{ID: "internal_server_error", Message: "internal server error"}},
		Data:   make(map[string]interface{}),
		Events: RuntimeEvents{},
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

func (b *Backtest) SetCurrentBarInfo(info *BarInfo) {
	b.currentBarInfo = info
}

func (b *Backtest) GetRuntimeEvents() interface{} {
	return b.runtimeEvents
}

func (b *Backtest) Close() {
	b.zmqConn.Close()
}
