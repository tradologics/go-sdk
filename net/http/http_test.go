package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-sdk/backtest"
	"io"
	"log"
	"strings"
	"testing"
)

const invalidErrorMsg = "invalid response"

// TODO Remove it from git
const token = "***REMOVED***"

func turnOnBacktestModel() {
	SetBacktestMode("2020-07-01 21:00:00.000000", "2020-07-01 21:00:00.000000")
}

func removeBacktestMode() {
	Backtest = nil
}

func removeToken() {
	Token = ""
}

func cls(b io.ReadCloser) {
	b.Close()
}

func TestHttpGetHost(t *testing.T) {
	res, err := Get("https://google.com")
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "google") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestHttpGetHostInvalidSchema(t *testing.T) {
	_, err := Get("htt://google.com")
	if err != nil {
		assert.Equal(t, "Get \"htt://google.com\": unsupported protocol scheme \"htt\"", err.Error())
	} else {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestHttpPostHost(t *testing.T) {
	payload, err := json.Marshal(map[string]interface{}{
		"title":     "my simple todo",
		"completed": false,
	})

	if err != nil {
		assert.Error(t, err)
	}

	res, err := Post("https://google.com", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 405, res.StatusCode)
}

func TestTradologicsGetWithEmptyToken(t *testing.T) {
	_, err := Get("/me")
	if err != nil {
		assert.Equal(t, "please use `SetToken(...)` first", err.Error())
	} else {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestTradologicsGetWithToken(t *testing.T) {
	SetToken(token)
	defer removeToken()

	res, err := Get("/me")
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "{\"errors\":[],\"data\":{\"name\"") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}

}

func TestTradologicsGetWithInvalidToken(t *testing.T) {
	SetToken("Bearer")
	defer removeToken()

	res, err := Get("/me")
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 401, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "{\"errors\":[{\"id\":\"authentication_error\",\"message\":\"Token cannot be validated. Please make sure you are using a valid and active token.\"}],\"data\":{}}") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestTradologicsPostWithEmptyToken(t *testing.T) {
	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.Error(t, err)
	}

	_, err = Post("/accounts", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		assert.Equal(t, "please use `SetToken(...)` first", err.Error())
	} else {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestTradologicsPostWithToken(t *testing.T) {
	SetToken(token)
	defer removeToken()

	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.Error(t, err)
	}

	res, err := Post("/accounts", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 400, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "{\"errors\":[{\"id\":\"invalid_request\",\"message\":\"data should have required property") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestTradologicsPostWithInvalidToken(t *testing.T) {
	SetToken("Bearer")
	defer removeToken()

	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.Error(t, err)
	}

	res, err := Post("/accounts", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 401, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "{\"errors\":[{\"id\":\"authentication_error\",\"message\":\"Token cannot be validated. Please make sure you are using a valid and active token.\"}],\"data\":{}}") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestTradologicsNewRequestWithEmptyToken(t *testing.T) {
	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.Error(t, err)
	}

	req, err := NewRequest(MethodGet, "/me", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	c := DefaultClient
	_, err = c.Do(req)
	if err != nil {
		assert.Equal(t, "please use `SetToken(...)` first", err.Error())
	} else {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestTradologicsNewRequestWithToken(t *testing.T) {
	SetToken(token)
	defer removeToken()

	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.Error(t, err)
	}

	req, err := NewRequest(MethodGet, "/me", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	c := DefaultClient
	res, err := c.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		assert.Error(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)
}

func TestTradologicsNewRequestWithInvalidToken(t *testing.T) {
	SetToken("Bearer")
	defer removeToken()

	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.Error(t, err)
	}

	req, err := NewRequest(MethodGet, "/me", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	c := DefaultClient
	res, err := c.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		assert.Error(t, err)
	}

	assert.Equal(t, 401, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "{\"errors\":[{\"id\":\"authentication_error\",\"message\":\"Token cannot be validated. Please make sure you are using a valid and active token.\"}],\"data\":{}}") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestHttpNewRequestWithEmptyToken(t *testing.T) {
	req, err := NewRequest(MethodGet, "https://google.com/", nil)
	req.Header.Set("Content-Type", "application/json")

	c := DefaultClient
	res, err := c.Do(req)
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "google") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestSetCurrentBarInfoWithoutBacktest(t *testing.T) {
	err := SetCurrentBarInfo(&backtest.BarInfo{
		Datetime:   "2020-07-01 21:00:00.000000",
		Resolution: "1m",
	})
	if err != nil {
		assert.Equal(t, "please set backtest mode first", err.Error())
	} else {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestSetCurrentBarInfoWithBacktest(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	err := SetCurrentBarInfo(&backtest.BarInfo{
		Datetime:   "2020-07-01 21:00:00.000000",
		Resolution: "1m",
	})
	if err != nil {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestBacktestGetWithInvalidUrl(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	res, err := Get("/com")
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 502, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "{\"status\":502,\"errors\":[{\"id\":\"internal_server_error\",\"message\":\"Endpoint not found\"}],\"data\":{},\"events\":{}}") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestBacktestGetWithoutInfo(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	res, err := Get("/accounts")
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "{\"status\":200,\"errors\":[],\"data\":[{\"account\":null,\"account_id\":\"backtest\",\"blocked\":false,\"broker\":\"tradologics\",\"buying_power\":null,\"cash\":100000,\"currency\":\"USD\",\"daytrade_count\":null,\"daytrading_buying_power\":null,\"equity\":null,\"initial_margin\":1,\"maintenance_margin\":null,\"multiplier\":null,\"name\":\"paper\",\"pattern_day_trader\":false,\"regt_buying_power\":null,\"shorting_enabled\":false,\"sma\":null,\"status\":null}],\"events\":{}}") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestBacktestGetWithInfo(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	err := SetCurrentBarInfo(&backtest.BarInfo{
		Datetime:   "2020-07-01 21:00:00.000000",
		Resolution: "1m",
	})
	if err != nil {
		assert.Error(t, err)
	}

	res, err := Get("/accounts")
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}

	if !strings.Contains(string(body), "{\"status\":200,\"errors\":[],\"data\":[{\"account\":null,\"account_id\":\"backtest\",\"blocked\":false,\"broker\":\"tradologics\",\"buying_power\":null,\"cash\":100000,\"currency\":\"USD\",\"daytrade_count\":null,\"daytrading_buying_power\":null,\"equity\":null,\"initial_margin\":1,\"maintenance_margin\":null,\"multiplier\":null,\"name\":\"paper\",\"pattern_day_trader\":false,\"regt_buying_power\":null,\"shorting_enabled\":false,\"sma\":null,\"status\":null}],\"events\":{}}") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestGetRuntimeEventsWithoutBacktest(t *testing.T) {
	_, err := GetRuntimeEvents()
	if err != nil {
		assert.Equal(t, "please set backtest mode first", err.Error())
	} else {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestGetRuntimeEventsWithBacktest(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	data, err := GetRuntimeEvents()
	if err != nil {
		assert.Equal(t, true, false, invalidErrorMsg)
	} else {
		assert.NotEqual(t, nil, data)
	}
}

// TODO
func TestBacktestPostWithoutInfoAndValidData(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	type Rule struct {
		Type   string `json:"type"`
		Target int    `json:"target"`
	}

	type Data struct {
		Type       string   `json:"type"`
		Rule       Rule     `json:"rule"`
		Asset      string   `json:"asset"`
		Price      float64  `json:"price"`
		Strategies []string `json:"strategies"`
	}

	data, err := json.Marshal(Data{
		Type: "price",
		Rule: Rule{
			Type:   "above",
			Target: 10,
		},
		Asset:      "AAPL",
		Price:      123.0,
		Strategies: []string{"my-strategy", "my-second-strategy"},
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := Post("/monitors", "application/json", bytes.NewBuffer(data))
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 400, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}
	fmt.Println(string(body))
	if !strings.Contains(string(body), "{\"status\":200,\"errors\":[],\"data\":[{\"account\":null,\"account_id\":\"backtest\",\"blocked\":false,\"broker\":\"tradologics\",\"buying_power\":null,\"cash\":100000,\"currency\":\"USD\",\"daytrade_count\":null,\"daytrading_buying_power\":null,\"equity\":null,\"initial_margin\":1,\"maintenance_margin\":null,\"multiplier\":null,\"name\":\"paper\",\"pattern_day_trader\":false,\"regt_buying_power\":null,\"shorting_enabled\":false,\"sma\":null,\"status\":null}],\"events\":{}}") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestBacktestPostWithoutInfoAndInvalidData(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	data, err := json.Marshal(map[string]interface{}{
		"test": "my simple test",
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := Post("/monitors", "application/json", bytes.NewBuffer(data))
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, 400, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err)
	}
	if !strings.Contains(string(body), "{\"status\":400,\"errors\":[{\"id\":\"invalid_request\",\"message\":\"data should have required property 'type'\"},{\"id\":\"invalid_request\",\"message\":\"data should have required property 'strategies'\"},{\"id\":\"invalid_request\",\"message\":\"data should have required property 'rule'\"}],\"data\":null,\"events\":{}}") {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

//
//func TestBacktestPostWithInfo(t *testing.T) {
//	turnOnBacktestModel()
//	defer removeBacktestMode()
//
//	err := SetCurrentBarInfo(&backtest.BarInfo{
//		Datetime:   "2020-07-01 21:00:00.000000",
//		Resolution: "1m",
//	})
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	res, err := Get("/accounts")
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	assert.Equal(t, 200, res.StatusCode)
//
//	defer cls(res.Body)
//
//	body, err := io.ReadAll(res.Body)
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	if !strings.Contains(string(body), "{\"status\":200,\"errors\":[],\"data\":[{\"account\":null,\"account_id\":\"backtest\",\"blocked\":false,\"broker\":\"tradologics\",\"buying_power\":null,\"cash\":100000,\"currency\":\"USD\",\"daytrade_count\":null,\"daytrading_buying_power\":null,\"equity\":null,\"initial_margin\":1,\"maintenance_margin\":null,\"multiplier\":null,\"name\":\"paper\",\"pattern_day_trader\":false,\"regt_buying_power\":null,\"shorting_enabled\":false,\"sma\":null,\"status\":null}],\"events\":{}}") {
//		assert.Equal(t, true, false, invalidErrorMsg)
//	}
//}
//
//func TestBacktestNewRequestWithoutInfo(t *testing.T) {
//	turnOnBacktestModel()
//	defer removeBacktestMode()
//
//	res, err := Get("/accounts")
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	assert.Equal(t, 200, res.StatusCode)
//
//	defer cls(res.Body)
//
//	body, err := io.ReadAll(res.Body)
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	if !strings.Contains(string(body), "{\"status\":200,\"errors\":[],\"data\":[{\"account\":null,\"account_id\":\"backtest\",\"blocked\":false,\"broker\":\"tradologics\",\"buying_power\":null,\"cash\":100000,\"currency\":\"USD\",\"daytrade_count\":null,\"daytrading_buying_power\":null,\"equity\":null,\"initial_margin\":1,\"maintenance_margin\":null,\"multiplier\":null,\"name\":\"paper\",\"pattern_day_trader\":false,\"regt_buying_power\":null,\"shorting_enabled\":false,\"sma\":null,\"status\":null}],\"events\":{}}") {
//		assert.Equal(t, true, false, invalidErrorMsg)
//	}
//}
//
//func TestBacktestNewRequestWithInfo(t *testing.T) {
//	turnOnBacktestModel()
//	defer removeBacktestMode()
//
//	err := SetCurrentBarInfo(&backtest.BarInfo{
//		Datetime:   "2020-07-01 21:00:00.000000",
//		Resolution: "1m",
//	})
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	res, err := Get("/accounts")
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	assert.Equal(t, 200, res.StatusCode)
//
//	defer cls(res.Body)
//
//	body, err := io.ReadAll(res.Body)
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	if !strings.Contains(string(body), "{\"status\":200,\"errors\":[],\"data\":[{\"account\":null,\"account_id\":\"backtest\",\"blocked\":false,\"broker\":\"tradologics\",\"buying_power\":null,\"cash\":100000,\"currency\":\"USD\",\"daytrade_count\":null,\"daytrading_buying_power\":null,\"equity\":null,\"initial_margin\":1,\"maintenance_margin\":null,\"multiplier\":null,\"name\":\"paper\",\"pattern_day_trader\":false,\"regt_buying_power\":null,\"shorting_enabled\":false,\"sma\":null,\"status\":null}],\"events\":{}}") {
//		assert.Equal(t, true, false, invalidErrorMsg)
//	}
//}

// TODO PostForm
