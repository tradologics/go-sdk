package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func cls(b io.ReadCloser) {
	b.Close()
}

func TestRunServerWithStrategy(t *testing.T) {

	// Create handler
	strategy := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("ok"))
		if err != nil {
			assert.Error(t, err, "handler error")
		}
	}

	// Start server
	server := httptest.NewServer(router(strategy, "/"))
	defer server.Close()

	// Get request
	res, err := http.Get(fmt.Sprintf("%s/", server.URL))
	if err != nil {
		assert.Error(t, err, "server error")
	}

	body, err := io.ReadAll(res.Body)
	assert.NotEqual(t, "ok", string(body), "invalid response")
	cls(res.Body)

	// Post request
	res, err = http.Post(fmt.Sprintf("%s/", server.URL), "", nil)
	if err != nil {
		assert.Error(t, err, "server error")
	}

	body, err = io.ReadAll(res.Body)
	assert.Equal(t, "ok", string(body), "invalid response")
	cls(res.Body)
}
