package main

import (
	"fmt"
	"go-sdk/net/http"
	"io"
)

type User struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

func main() {
	// todo timeout???

	c := http.DefaultClient

	http.SetBacktestMode("2020-07-01 21:00:00.000000", "2020-07-01 21:00:00.000000")
	http.SetToken("***REMOVED***")

	res, _ := c.Get("/accounts")
	//payload, err := json.Marshal(map[string]interface{}{
	//	"title":     "my simple todo",
	//	"completed": false,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Todo new request and request
	//req, err := _http.NewRequest(_http.MethodPatch, "https://postb.in/1639571630983-4007084534969", bytes.NewBuffer(payload))
	//req.Header.Set("Content-Type", "application/json")

	//res, err := c.Do(req)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	//Convert the body to type string
	sb := string(body)
	fmt.Print(sb, "\n")

}
