package backtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ErocRequestHeader struct {
	Start      string `json:"start"`
	End        string `json:"end"`
	Datetime   string `json:"datetime"`
	Resolution string `json:"resolution"`
}

type ErocRequest struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Data    []byte            `json:"data"`
	Headers ErocRequestHeader `json:"headers"`
}

type ErocResponseData struct {
	Name        string `json:"name"`
	StrategyID  string `json:"strategy_id"`
	Description string `json:"description"`
	AsTradelet  string `json:"as_tradelet"`
	Mode        string `json:"mode"`
	Url         string `json:"url"`
	Public      bool   `json:"public"`
	Datetime    string `json:"datetime"`
}

type ErocError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type RuntimeEvents map[string]string

type ErocResponse struct {
	Status int              `json:"status"`
	Errors []ErocError      `json:"errors"`
	Data   ErocResponseData `json:"data"`
	Events RuntimeEvents    `json:"events"`
}

type BarInfo struct {
	datetime   string
	resolution string
}

type Backtest struct {
	start          string
	end            string
	currentBarInfo *BarInfo
	runtimeEvents  RuntimeEvents
	zmqConn        *ZmqConn
}

func NewBacktest(start, end, socketUrl string) *Backtest {
	zmqConn := NewZmq(socketUrl)

	return &Backtest{
		start:          start,
		end:            end,
		currentBarInfo: &BarInfo{},
		zmqConn:        zmqConn,
	}
}

func (b *Backtest) CallErocMethod(req *http.Request) (*http.Response, error) {
	var erocResponse ErocResponse

	erocRequest := &ErocRequest{
		Method: req.Method,
		Url:    req.URL.String(),
		Headers: ErocRequestHeader{
			Start:      b.start,
			End:        b.end,
			Datetime:   b.currentBarInfo.datetime,
			Resolution: b.currentBarInfo.resolution,
		},
	}

	err := b.zmqConn.SendJSON(&erocRequest)
	if err != nil {
		return nil, err
	}

	err = b.zmqConn.ReceiveJSON(&erocResponse)
	if err != nil {
		return nil, err
	}

	b.runtimeEvents = erocResponse.Events

	erocJSONResponse, err := json.Marshal(erocResponse)
	if err != nil {
		return nil, err
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

	return res, nil

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
