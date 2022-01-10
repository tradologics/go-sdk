package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var strategyHandler func(tradehook string, payload []byte)

// strategyWrapper retrieve tradehook information from response body and send it to customer strategy
func strategyWrapper(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Tradehook string `json:"tradehook"`
		Payload   []byte `json:"payload"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Fatal(err)
	}

	strategyHandler(requestBody.Tradehook, requestBody.Payload)

	w.WriteHeader(200)
	_, err = w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Fatal(err)
	}
}

// postMethodOnlyHandler validate request method and execute only 'POST'
func postMethodOnlyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		strategyWrapper(w, r)
	} else {
		w.WriteHeader(405)
		_, err := w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		if err != nil {
			log.Fatal(err)
		}
	}
}

// router returns http server mux with selected handler and URL path
func router(strategy func(tradehook string, payload []byte), endpoint string) http.Handler {
	strategyHandler = strategy

	router := http.NewServeMux()
	router.HandleFunc(endpoint, postMethodOnlyHandler)

	return router
}

// Start create new server with selected host and port, and use strategy as request handler
func Start(strategy func(tradehook string, payload []byte), endpoint, host string, port int) {

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router(strategy, endpoint))
	if err != nil {
		log.Fatalln(err)
	}
}
