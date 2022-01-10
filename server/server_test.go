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
		var payloadData map[string]interface{}

		err := json.Unmarshal(payload, &payloadData)
		if err != nil {
			assert.NoError(t, err)
		}

		assert.Equal(t, "bars", tradehook, invalidErrorMsg)
		assert.NotEqualf(t, 0, len(payloadData), invalidErrorMsg)
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
		Event string            `json:"event"`
		Data  map[string]string `json:"data"`
	}{
		Event: "bars",
		Data:  map[string]string{"foo": "boo", "boo": "foo"},
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
