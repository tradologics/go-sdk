package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tradologics/go-sdk/backtest"
	"github.com/tradologics/go-sdk/config"
	"io"
	"log"
	"reflect"
	"strings"
	"testing"
)

const invalidErrorMsg = "invalid response"

func authInit() {
	cfg := config.GetTestConfig("../../.env")
	SetToken(cfg.SandboxToken)
}

func turnOnBacktestModel() {
	err := SetBacktestMode("2021-01-01 21:00:00.000000", "2021-01-08 21:00:00.000000")
	if err != nil {
		log.Fatal(err)
	}
}

func setCurrentBarInfo() {
	err := SetCurrentBarInfo(&backtest.BarInfo{
		Datetime:   "2020-07-01 21:00:00.000000",
		Resolution: "1d",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func removeBacktestMode() {
	Backtest = nil
}

func removeToken() {
	SetToken("")
}

func cls(b io.ReadCloser) {
	b.Close()
}

type MonitorDataRule struct {
	Type   string `json:"type"`
	Target int    `json:"target"`
}

type MonitorData struct {
	Type       string          `json:"type"`
	Rule       MonitorDataRule `json:"rule"`
	Asset      string          `json:"asset"`
	Price      float64         `json:"price"`
	Strategies []string        `json:"strategies"`
}

func createMonitorsPostData(strategies []string) ([]byte, error) {
	data, err := json.Marshal(MonitorData{
		Type: "price",
		Rule: MonitorDataRule{
			Type:   "above",
			Target: 10,
		},
		Asset:      "AAPL",
		Price:      123.0,
		Strategies: strategies,
	})
	return data, err
}

func TestHttpGetHost(t *testing.T) {
	res, err := Get("https://google.com")
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
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
		"title":     "my simple data",
		"completed": false,
	})

	if err != nil {
		assert.NoError(t, err)
	}

	res, err := Post("https://google.com", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		assert.NoError(t, err)
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
	authInit()
	defer removeToken()

	res, err := Get("/me")
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(string(body), "{\"errors\":[],\"data\":{\"name\""))
}

func TestTradologicsGetWithInvalidToken(t *testing.T) {
	SetToken("Bearer")
	defer removeToken()

	res, err := Get("/me")
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 401, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[{\"id\":\"authentication_error\",\"message\":\"Token cannot be validated. Please make sure you are using a valid and active token.\"}],\"data\":{}}"),
	)
}

func TestTradologicsPostWithEmptyToken(t *testing.T) {
	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.NoError(t, err)
	}

	_, err = Post("/accounts", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		assert.Equal(t, "please use `SetToken(...)` first", err.Error())
	} else {
		assert.Equal(t, true, false, invalidErrorMsg)
	}
}

func TestTradologicsPostWithToken(t *testing.T) {
	authInit()
	defer removeToken()

	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.NoError(t, err)
	}

	res, err := Post("/accounts", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 400, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[{\"id\":\"invalid_request\",\"message\":\"data should have required property"),
	)
}

func TestTradologicsPostWithInvalidToken(t *testing.T) {
	SetToken("Bearer")
	defer removeToken()

	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.NoError(t, err)
	}

	res, err := Post("/accounts", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 401, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[{\"id\":\"authentication_error\",\"message\":\"Token cannot be validated. Please make sure you are using a valid and active token.\"}],\"data\":{}}"),
	)
}

func TestTradologicsNewRequestWithEmptyToken(t *testing.T) {
	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.NoError(t, err)
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
	authInit()
	defer removeToken()

	payload, err := json.Marshal(map[string]interface{}{
		"test": true,
	})

	if err != nil {
		assert.NoError(t, err)
	}

	req, err := NewRequest(MethodGet, "/me", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	c := DefaultClient
	res, err := c.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		assert.NoError(t, err)
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
		assert.NoError(t, err)
	}

	req, err := NewRequest(MethodGet, "/me", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	c := DefaultClient
	res, err := c.Do(req)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 401, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[{\"id\":\"authentication_error\",\"message\":\"Token cannot be validated. Please make sure you are using a valid and active token.\"}],\"data\":{}}"),
	)
}

func TestHttpNewRequestWithEmptyToken(t *testing.T) {
	req, err := NewRequest(MethodGet, "https://google.com/", nil)
	req.Header.Set("Content-Type", "application/json")

	c := DefaultClient
	res, err := c.Do(req)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"google"),
	)
}

func TestSetCurrentBarInfoWithoutBacktest(t *testing.T) {
	err := SetCurrentBarInfo(&backtest.BarInfo{
		Datetime:   "2020-07-01 21:00:00.000000",
		Resolution: "1m",
	})
	if err != nil {
		assert.Equal(t, "please set backtest mode first", err.Error())
	} else {
		assert.Error(t, err)
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
		assert.NoError(t, err)
	}

	assert.Equal(t, 502, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[{\"id\":\"internal_server_error\",\"message\":\"Endpoint not found\"}],\"data\":{}}"),
	)
}

func TestBacktestGetWithoutInfo(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	res, err := Get("/accounts")
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":[{\"account\":null,\"account_id\":\"backtest\",\"blocked\":false,\"broker\":\"tradologics\",\"buying_power\":null,\"cash\":100000,\"currency\":\"USD\",\"daytrade_count\":null,\"daytrading_buying_power\":null,\"equity\":null,\"initial_margin\":1,\"maintenance_margin\":null,\"multiplier\":null,\"name\":\"paper\",\"pattern_day_trader\":false,\"regt_buying_power\":null,\"shorting_enabled\":false,\"sma\":null,\"status\":null}]}"),
	)
}

func TestBacktestGetWithInfo(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	setCurrentBarInfo()

	res, err := Get("/accounts")
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 200, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":[{\"account\":null,\"account_id\":\"backtest\",\"blocked\":false,\"broker\":\"tradologics\",\"buying_power\":null,\"cash\":100000,\"currency\":\"USD\",\"daytrade_count\":null,\"daytrading_buying_power\":null,\"equity\":null,\"initial_margin\":1,\"maintenance_margin\":null,\"multiplier\":null,\"name\":\"paper\",\"pattern_day_trader\":false,\"regt_buying_power\":null,\"shorting_enabled\":false,\"sma\":null,\"status\":null}]}"),
	)
}

func TestGetRuntimeEventsWithoutBacktest(t *testing.T) {
	_, err := GetRuntimeEvents()
	if err != nil {
		assert.Equal(t, "please set backtest mode first", err.Error())
	} else {
		assert.Error(t, err)
	}
}

func TestGetRuntimeEventsWithBacktest(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	data, err := GetRuntimeEvents()
	if err != nil {
		assert.Equal(t, true, false, invalidErrorMsg)
	} else {
		assert.IsType(t, reflect.Interface, reflect.TypeOf(data).Kind())
	}
}

func TestBacktestPostWithoutInfoAndValidDataWithoutBarInfo(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	data, err := createMonitorsPostData([]string{"demo-strategy"})
	if err != nil {
		assert.NoError(t, err)
	}

	res, err := Post("/monitors", "application/json", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 201, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":{\"active_at\":\"\",\"asset\":{\"currency\":\"USD \",\"exchange\":\"XNAS\",\"figi\":\"BBG000B9XRY4\",\"name\":\"Apple Inc\",\"security_type\":\"CS\",\"ticker\":\"AAPL\",\"tsid\":\"TXS0005PKIKN\",\"tuid\":\"TXU000BB2K0H\"},\"canceled_at\":null,\"comment\":null,\"id\":\""),
	)
}

func TestBacktestPostWithoutInfoAndValidDataWithBarInfo(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	setCurrentBarInfo()

	data, err := createMonitorsPostData([]string{"demo-strategy"})
	if err != nil {
		assert.NoError(t, err)
	}

	res, err := Post("/monitors", "application/json", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 201, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":{\"active_at\":\"2020-07-01 21:00:00.000000\",\"asset\":{\"currency\":\"USD \",\"exchange\":\"XNAS\",\"figi\":\"BBG000B9XRY4\",\"name\":\"Apple Inc\",\"security_type\":\"CS\",\"ticker\":\"AAPL\",\"tsid\":\"TXS0005PKIKN\",\"tuid\":\"TXU000BB2K0H\"},\"canceled_at\":null,\"comment\":null,\"id\""),
	)
}

// TODO timeout error?
//func TestBacktestPostWithoutInfoAndInvalidStrategyData(t *testing.T) {
//	turnOnBacktestModel()
//	defer removeBacktestMode()
//
//	data, err := createMonitorsPostData([]string{"foo"})
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	res, err := Post("/monitors", "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	assert.Equal(t, 201, res.StatusCode)
//
//	defer cls(res.Body)
//
//	body, err := io.ReadAll(res.Body)
//	if err != nil {
//		assert.Error(t, err)
//	}
//
//	fmt.Println(string(body))
//	if !strings.Contains(string(body), "{\"errors\":[],\"data\":{\"active_at\":\"2020-07-01 21:00:00.000000\",\"asset\":{\"currency\":\"USD \",\"exchange\":\"XNAS\",\"figi\":\"BBG000B9XRY4\",\"name\":\"Apple Inc\",\"security_type\":\"CS\",\"ticker\":\"AAPL\",\"tsid\":\"TXS0005PKIKN\",\"tuid\":\"TXU000BB2K0H\"},\"canceled_at\":null,\"comment\":null,\"id\"") {
//		assert.Equal(t, true, false, invalidErrorMsg)
//	}
//}

func TestBacktestPostWithoutInfoAndValidAndInvalidStrategyData(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	data, err := createMonitorsPostData([]string{"demo-strategy", "foo"})
	if err != nil {
		assert.NoError(t, err)
	}

	res, err := Post("/monitors", "application/json", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 201, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":{\"active_at\":\"\",\"asset\":{\"currency\":\"USD \",\"exchange\":\"XNAS\",\"figi\":\"BBG000B9XRY4\",\"name\":\"Apple Inc\",\"security_type\":\"CS\",\"ticker\":\"AAPL\",\"tsid\":\"TXS0005PKIKN\",\"tuid\":\"TXU000BB2K0H\"},\"canceled_at\":null,\"comment\":null,\"id\""),
	)
}

func TestBacktestPostWithInfoAndValidAndInvalidStrategyData(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	setCurrentBarInfo()

	data, err := createMonitorsPostData([]string{"demo-strategy", "foo"})
	if err != nil {
		assert.NoError(t, err)
	}

	res, err := Post("/monitors", "application/json", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 201, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":{\"active_at\":\"2020-07-01 21:00:00.000000\",\"asset\":{\"currency\":\"USD \",\"exchange\":\"XNAS\",\"figi\":\"BBG000B9XRY4\",\"name\":\"Apple Inc\",\"security_type\":\"CS\",\"ticker\":\"AAPL\",\"tsid\":\"TXS0005PKIKN\",\"tuid\":\"TXU000BB2K0H\"},\"canceled_at\":null,\"comment\":null,\"id\""),
	)
}

func TestBacktestPostWithInvalidFieldType(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	setCurrentBarInfo()

	type InvalidFiendTypeMonitorData struct {
		Type       int             `json:"type"`
		Rule       MonitorDataRule `json:"rule"`
		Asset      string          `json:"asset"`
		Price      float64         `json:"price"`
		Strategies []string        `json:"strategies"`
	}

	data, err := json.Marshal(InvalidFiendTypeMonitorData{
		// Invalid type
		Type: 100500,

		Rule: MonitorDataRule{
			Type:   "above",
			Target: 10,
		},
		Asset:      "AAPL",
		Price:      123.0,
		Strategies: []string{"demo-strategy"},
	})

	if err != nil {
		assert.NoError(t, err)
	}

	res, err := Post("/monitors", "application/json", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 400, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[{\"id\":\"invalid_request\",\"message\":\"data.type should be string\"},{\"id\":\"invalid_request\",\"message\":\"data.type should be equal to one of the allowed values\"}],\"data\":null}"),
	)
}

func TestBacktestPostWithoutInfoAndInvalidFieldValue(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	setCurrentBarInfo()

	data, err := json.Marshal(MonitorData{
		// Invalid value
		Type: "string",

		Rule: MonitorDataRule{
			Type:   "above",
			Target: 10,
		},
		Asset:      "AAPL",
		Price:      123.0,
		Strategies: []string{"demo-strategy"},
	})

	if err != nil {
		assert.NoError(t, err)
	}

	res, err := Post("/monitors", "application/json", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 400, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[{\"id\":\"invalid_request\",\"message\":\"data.type should be equal to one of the allowed values\"}],\"data\":null}"),
	)
}

// TODO PostForm

func TestBacktestNewPostRequestWithoutInfoAndValidDataWithoutBarInfo(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	data, err := createMonitorsPostData([]string{"demo-strategy"})
	if err != nil {
		assert.NoError(t, err)
	}

	req, err := NewRequest("POST", "/monitors", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	c := &Client{}
	res, err := c.Do(req)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 201, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":{\"active_at\":\"\",\"asset\":{\"currency\":\"USD \",\"exchange\":\"XNAS\",\"figi\":\"BBG000B9XRY4\",\"name\":\"Apple Inc\",\"security_type\":\"CS\",\"ticker\":\"AAPL\",\"tsid\":\"TXS0005PKIKN\",\"tuid\":\"TXU000BB2K0H\"},\"canceled_at\":null,\"comment\":null,\"id\":\""),
	)
}

func TestBacktestNewPostRequestWithoutInfoAndValidDataWithBarInfo(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	setCurrentBarInfo()

	data, err := createMonitorsPostData([]string{"demo-strategy"})
	if err != nil {
		assert.NoError(t, err)
	}

	req, err := NewRequest("POST", "/monitors", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	c := &Client{}
	res, err := c.Do(req)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 201, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":{\"active_at\":\"2020-07-01 21:00:00.000000\",\"asset\":{\"currency\":\"USD \",\"exchange\":\"XNAS\",\"figi\":\"BBG000B9XRY4\",\"name\":\"Apple Inc\",\"security_type\":\"CS\",\"ticker\":\"AAPL\",\"tsid\":\"TXS0005PKIKN\",\"tuid\":\"TXU000BB2K0H\"},\"canceled_at\":null,\"comment\":null,\"id\""),
	)
}

func TestBacktestNewPostRequestWithoutInfoAndValidAndInvalidStrategyData(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	data, err := createMonitorsPostData([]string{"demo-strategy", "foo"})
	if err != nil {
		assert.NoError(t, err)
	}

	req, err := NewRequest("POST", "/monitors", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	c := &Client{}
	res, err := c.Do(req)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 201, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":{\"active_at\":\"\",\"asset\":{\"currency\":\"USD \",\"exchange\":\"XNAS\",\"figi\":\"BBG000B9XRY4\",\"name\":\"Apple Inc\",\"security_type\":\"CS\",\"ticker\":\"AAPL\",\"tsid\":\"TXS0005PKIKN\",\"tuid\":\"TXU000BB2K0H\"},\"canceled_at\":null,\"comment\":null,\"id\""),
	)
}

func TestBacktestNewPostRequestWithInfoAndValidAndInvalidStrategyData(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	setCurrentBarInfo()

	data, err := createMonitorsPostData([]string{"demo-strategy", "foo"})
	if err != nil {
		assert.NoError(t, err)
	}

	req, err := NewRequest("POST", "/monitors", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	c := &Client{}
	res, err := c.Do(req)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 201, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[],\"data\":{\"active_at\":\"2020-07-01 21:00:00.000000\",\"asset\":{\"currency\":\"USD \",\"exchange\":\"XNAS\",\"figi\":\"BBG000B9XRY4\",\"name\":\"Apple Inc\",\"security_type\":\"CS\",\"ticker\":\"AAPL\",\"tsid\":\"TXS0005PKIKN\",\"tuid\":\"TXU000BB2K0H\"},\"canceled_at\":null,\"comment\":null,\"id\""),
	)
}

func TestBacktestNewPostRequestWithInvalidFieldType(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	setCurrentBarInfo()

	type InvalidFiendTypeMonitorData struct {
		Type       int             `json:"type"`
		Rule       MonitorDataRule `json:"rule"`
		Asset      string          `json:"asset"`
		Price      float64         `json:"price"`
		Strategies []string        `json:"strategies"`
	}

	data, err := json.Marshal(InvalidFiendTypeMonitorData{
		// Invalid type
		Type: 100500,

		Rule: MonitorDataRule{
			Type:   "above",
			Target: 10,
		},
		Asset:      "AAPL",
		Price:      123.0,
		Strategies: []string{"demo-strategy"},
	})

	if err != nil {
		assert.NoError(t, err)
	}

	req, err := NewRequest("POST", "/monitors", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	c := &Client{}
	res, err := c.Do(req)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 400, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[{\"id\":\"invalid_request\",\"message\":\"data.type should be string\"},{\"id\":\"invalid_request\",\"message\":\"data.type should be equal to one of the allowed values\"}],\"data\":null}"),
	)
}

func TestBacktestNewPostRequestWithoutInfoAndInvalidFieldValue(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	setCurrentBarInfo()

	data, err := json.Marshal(MonitorData{
		// Invalid value
		Type: "string",

		Rule: MonitorDataRule{
			Type:   "above",
			Target: 10,
		},
		Asset:      "AAPL",
		Price:      123.0,
		Strategies: []string{"demo-strategy"},
	})

	if err != nil {
		assert.NoError(t, err)
	}

	req, err := NewRequest("POST", "/monitors", bytes.NewBuffer(data))
	if err != nil {
		assert.NoError(t, err)
	}

	c := &Client{}
	res, err := c.Do(req)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 400, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"errors\":[{\"id\":\"invalid_request\",\"message\":\"data.type should be equal to one of the allowed values\"}],\"data\":null}"),
	)
}

func TestBacktestPostInvalidJSON(t *testing.T) {
	turnOnBacktestModel()
	defer removeBacktestMode()

	res, err := Post("/monitors", "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 502, res.StatusCode)

	defer cls(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, strings.Contains(
		string(body),
		"{\"status\":502,\"errors\":[{\"id\":\"internal_server_error\",\"message\":\"Invalid JSON\"}],\"data\":{}"),
	)
}
