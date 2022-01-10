package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const invalidErrorMsg = "invalid response"

func cls(b io.ReadCloser) {
	b.Close()
}

func TestRunServerWithStrategy(t *testing.T) {

	// Create handler
	strategy := func(tradehook string, payload []byte) {
		var payloadData struct {
			Assets []string    `json:"assets"`
			Bars   interface{} `json:"bars"`
		}

		err := json.Unmarshal(payload, &payloadData)
		if err != nil {
			assert.NoError(t, err)
		}

		assert.Equal(t, "bars", tradehook, invalidErrorMsg)
		assert.NotEqualf(t, 0, len(payloadData.Assets), invalidErrorMsg)
	}

	// Start server
	server := httptest.NewServer(router(strategy, "/"))
	defer server.Close()

	// Get request
	res, err := http.Get(fmt.Sprintf("%s/", server.URL))
	if err != nil {
		assert.NoError(t, err)
	}

	body, err := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusText(http.StatusMethodNotAllowed), string(body), "invalid response")
	assert.Equal(t, 405, res.StatusCode, invalidErrorMsg)
	cls(res.Body)

	// Post request
	data, err := json.Marshal(struct {
		Tradehook string
		Payload   []byte
	}{
		Tradehook: "bars",
		Payload:   []byte("{\n  \"assets\": [\"AAPL:US\", \"BTCUSD:BMEX\", \"...\"],\n  \"bars\": {\n    \"YYYY-MM-01\": {\n      \"AAPL\": {\n        \"o\": \"113.79\",\n        \"h\": \"117.26\",\n        \"l\": \"113.62\",\n        \"c\": \"115.56\",\n        \"v\": 136210200,\n        \"t\": 782692,\n        \"w\": \"116.11665173913042\"\n      },\n      \"BTCUSD:BMEX\": {\n        \"o\": \"...\",\n        \"h\": \"...\",\n        \"l\": \"...\",\n        \"c\": \"...\",\n        \"v\": \"...\",\n        \"t\": \"...\",\n        \"w\": \"...\"\n      }\n    }\n  }\n}"),
	})

	res, err = http.Post(fmt.Sprintf("%s/", server.URL), "", ioutil.NopCloser(bytes.NewBuffer(data)))
	if err != nil {
		assert.NoError(t, err)
	}

	body, err = io.ReadAll(res.Body)
	assert.Equal(t, "OK", string(body), invalidErrorMsg)
	assert.Equal(t, 200, res.StatusCode, invalidErrorMsg)
	cls(res.Body)
}
