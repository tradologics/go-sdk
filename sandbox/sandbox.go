package sandbox

import (
	"fmt"
	go_sdk "github.com/tradologics/go-sdk"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var SandboxURL = "https://api.tradologics.com/v1/sandbox"

var Token string

func SetToken(token string) {
	Token = token
}

// setSandboxUrl give an ability to update sandbox address, useful on dev/test
func setSandboxURL(url string) {
	SandboxURL = url
}

func Tradehook(kind string, strategy func(string, []byte), args map[string]interface{}) {
	client := http.DefaultClient
	url := fmt.Sprintf("%s/%s", SandboxURL, strings.ReplaceAll(kind, "_", "/"))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	req.Header.Set("TGX-CLIENT", fmt.Sprintf("go-sdk/%s", go_sdk.Version))

	if len(args) > 0 {
		values := req.URL.Query()

		for k, v := range args {
			switch reflect.TypeOf(v).Kind() {
			case reflect.String:
				values.Add(k, v.(string))
			case reflect.Int:
				values.Add(k, strconv.Itoa(v.(int)))
			case reflect.Bool:
				values.Add(k, strconv.FormatBool(v.(bool)))
			default:
				// Do nothing
			}

			req.URL.RawQuery = values.Encode()
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	strategy(strings.Split(kind, "_")[0], body)
}

func Bar(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("bar", strategy, args)
}

func Monitor(kind string, strategy func(string, []byte), args map[string]interface{}) {
	Tradehook(kind, strategy, args)
}

func MonitorExpired(kind string, strategy func(string, []byte), args map[string]interface{}) {
	Tradehook(fmt.Sprintf("%s_expire", kind), strategy, args)
}

func PositionMonitor(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("position", strategy, args)
}

func PositionMonitorExpired(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("position_expire", strategy, args)
}

func PriceMonitor(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("price", strategy, args)
}

func PriceMonitorExpired(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("price_expire", strategy, args)
}

func Error(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("error", strategy, args)
}

func Order(kind string, strategy func(string, []byte), args map[string]interface{}) {
	Tradehook(fmt.Sprintf("order_%s", kind), strategy, args)
}

func OrderReceived(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_received", strategy, args)
}

func OrderPending(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_pending", strategy, args)
}

func OrderSubmitted(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_submitted", strategy, args)
}

func OrderSent(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_sent", strategy, args)
}

func OrderAccepted(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_accepted", strategy, args)
}

func OrderPartiallyFilled(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_partially_filled", strategy, args)
}

func OrderFilled(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_filled", strategy, args)
}

func OrderCanceled(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_canceled", strategy, args)
}

func OrderExpired(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_expired", strategy, args)
}

func OrderPendingCancel(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_pending_cancel", strategy, args)
}

func OrderRejected(strategy func(string, []byte), args map[string]interface{}) {
	Tradehook("order_rejected", strategy, args)
}
