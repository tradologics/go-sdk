# Tradologics Golang SDK

This is the initial version of Tradologics' Golang SDK.

The library supports a wrapper for the `net/http` that will automatically:

- prepend the full endpoint url to your calls
- attach your token to the request headers
- add `datetime` to your order when in backtesting mode

### Install

---

#### Install dependencies first

- [libzmq](DEPENDENCIES.md)
- [CZMQ](DEPENDENCIES.md)

#### Install the library

```sh
go get github/tradologics/go-sdk
```

### Using the library:

---

```golang
package main

import (
	"bytes"
	"encoding/json"
	"github.com/tradologics/go-sdk/net/http"
	"io"
	"log"
)

func main() {
	http.SetToken("YOUR TOKEN")

	data, err := json.Marshal("YOUR DATA STRUCT")

	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.Post("/orders", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalln(err)
	}
	...
}
```

### Running your own server:

---

```golang
package main

import (
	"github.com/tradologics/go-sdk/server"
	"net/http"
)

func strategyHandler(w http.ResponseWriter, r *http.Request) {
	...
}

func main() {
	server.Start(strategyHandler, "/my-strategy", "0.0.0.0", 5000)
}
```
