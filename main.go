package main

import (
	"fmt"
	"go-sdk/net/http"
	"io"
	"log"
)

func main() {

	c := http.DefaultClient
	http.SetBacktestMode("2020-07-01 21:00:00.000000", "2020-07-01 21:00:00.000000")
	http.SetToken("test")
	http.SetBacktestMode("1", "2")
	res, err := c.Get("google.com")
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	//Convert the body to type string
	sb := string(body)
	fmt.Printf(sb)

}
