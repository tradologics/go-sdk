package server

import (
	"fmt"
	"log"
	"net/http"
)

var strategyHandler func(w http.ResponseWriter, r *http.Request)

func postMethodOnlyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		strategyHandler(w, r)
	} else {
		w.WriteHeader(405)
		_, err := w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func router(strategy func(w http.ResponseWriter, r *http.Request), endpoint string) http.Handler {
	strategyHandler = strategy

	router := http.NewServeMux()
	router.HandleFunc(endpoint, postMethodOnlyHandler)

	return router
}

func Start(strategy func(w http.ResponseWriter, r *http.Request), endpoint, host string, port int) {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router(strategy, endpoint))
	if err != nil {
		log.Fatalln(err)
	}
}
