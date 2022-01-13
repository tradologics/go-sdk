package sandbox

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	go_sdk "github.com/tradologics/go-sdk"
	"testing"
	"time"
)

const (
	kindBar      = "bar"
	kindError    = "error"
	kindInvalid  = "foo"
	kindPosition = "position"
	kindOrder    = "order"
)

type tradehookErrorPayload struct {
	Errors []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	} `json:"errors"`
	Data interface{} `json:"data"`
}

type barPayload struct {
	Assets []string               `json:"assets"`
	Bars   map[string]interface{} `json:"bars"`
}

type monitorPayload struct {
	Event    string      `json:"event"`
	Rule     interface{} `json:"rule"`
	Position interface{} `json:"position"`
}

type orderPayload struct {
	OrderID       string      `json:"order_id"`
	Side          string      `json:"side"`
	Type          string      `json:"type"`
	Tif           string      `json:"tif"`
	ExtendedHours bool        `json:"extended_hours"`
	Qty           int         `json:"qty"`
	FilledQty     int         `json:"filled_qty"`
	LimitPrice    int         `json:"limit_price"`
	StopPrice     int         `json:"stop_price"`
	AvgFillPrice  int         `json:"avg_fill_price"`
	Status        string      `json:"status"`
	SubmittedAt   time.Time   `json:"submitted_at"`
	AcceptedAt    time.Time   `json:"accepted_at"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	FilledAt      time.Time   `json:"filled_at"`
	CanceledAt    time.Time   `json:"canceled_at"`
	ExpiredAt     interface{} `json:"expired_at"`
	RejectedAt    interface{} `json:"rejected_at"`
	Comment       interface{} `json:"comment"`
	StrategyID    string      `json:"strategy_id"`
	AccountID     string      `json:"account_id"`
	Asset         interface{} `json:"asset"`
}

func authInit() {
	cfg := go_sdk.GetTestConfig()
	setSandboxURL(cfg.SandboxURL)
	SetToken(cfg.SandboxToken)
}

type sandboxMethod func(strategy func(string, []byte), args map[string]interface{})
type sandboxMethodWithInputKind func(kind string, strategy func(string, []byte), args map[string]interface{})

func validateMethod(t *testing.T, sm sandboxMethod, kind string, p interface{}) {
	sm(
		// strategy function
		func(tradehook string, payload []byte) {
			assert.Equal(t, kind, tradehook)

			if err := json.Unmarshal(payload, &p); err != nil {
				assert.NoError(t, err)
			}
		},
		nil,
	)
}

func validateMethodWithInputKind(t *testing.T, smk sandboxMethodWithInputKind, inputKind, kind string, p interface{}) {
	smk(
		inputKind,
		// strategy function
		func(tradehook string, payload []byte) {
			assert.Equal(t, kind, tradehook)

			if err := json.Unmarshal(payload, &p); err != nil {
				assert.NoError(t, err)
			}
		},
		nil,
	)
}

func TestTradehookValidToken(t *testing.T) {
	authInit()
	validateMethodWithInputKind(t, Tradehook, kindBar, kindBar, barPayload{})
}

func TestTradehookInvalidToken(t *testing.T) {
	cfg := go_sdk.GetTestConfig()
	setSandboxURL(cfg.SandboxURL)
	SetToken("")

	Tradehook(
		kindError,
		// strategy function
		func(tradehook string, payload []byte) {
			var p tradehookErrorPayload

			assert.Equal(t, kindError, tradehook)

			if err := json.Unmarshal(payload, &p); err != nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, 1, len(p.Errors))
			assert.Equal(t, "authentication_error", p.Errors[0].ID)
			assert.Equal(t, "Token cannot be validated. Please make sure you are using a valid and active token.", p.Errors[0].Message)
		},
		nil,
	)
}

func TestTradehookInvalidKind(t *testing.T) {
	authInit()

	Tradehook(
		kindInvalid,
		// strategy function
		func(tradehook string, payload []byte) {
			var p tradehookErrorPayload

			assert.Equal(t, kindInvalid, tradehook)

			if err := json.Unmarshal(payload, &p); err != nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, 1, len(p.Errors))
			assert.Equal(t, "internal_server_error", p.Errors[0].ID)
			assert.Equal(t, "Endpoint /sandbox/foo not found", p.Errors[0].Message)
		},
		nil,
	)
}

func TestBarTradehook(t *testing.T) {
	authInit()
	validateMethod(t, Bar, kindBar, barPayload{})
}

func TestMonitorTradehook(t *testing.T) {
	authInit()
	validateMethodWithInputKind(t, Monitor, kindBar, kindBar, monitorPayload{})
}

func TestMonitorExpiredTradehook(t *testing.T) {
	authInit()
	validateMethodWithInputKind(t, MonitorExpired, kindPosition, kindPosition, monitorPayload{})
}

func TestPositionMonitorTradehook(t *testing.T) {
	authInit()
	validateMethod(t, PositionMonitor, kindPosition, monitorPayload{})
}

func TestPositionMonitorExpiredTradehook(t *testing.T) {
	authInit()
	validateMethod(t, PositionMonitorExpired, kindPosition, monitorPayload{})
}

//// TODO pricing api bug
//func TestPriceMonitorTradehook(t *testing.T) {
//	authInit()
//	validateMethod(t, PriceMonitor, kindPosition, monitorPayload{})
//}
//
//// TODO pricing api bug
//func TestPriceMonitorExpiredTradehook(t *testing.T) {
//	authInit()
//	validateMethod(t, PriceMonitorExpired, kindPosition, monitorPayload{})
//}

func TestErrorTradehook(t *testing.T) {
	authInit()

	Error(
		// strategy function
		func(tradehook string, payload []byte) {
			var p tradehookErrorPayload

			assert.Equal(t, kindError, tradehook)

			if err := json.Unmarshal(payload, &p); err != nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, 1, len(p.Errors))
			assert.Equal(t, "sandbox_error", p.Errors[0].ID)
			assert.Equal(t, "sandbox message error", p.Errors[0].Message)
		},
		//
		nil,
	)
}

func TestOrderTradehook(t *testing.T) {
	authInit()
	validateMethodWithInputKind(t, Order, kindOrder, kindOrder, orderPayload{})
}

func TestOrderReceivedTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderReceived, kindOrder, orderPayload{})
}

func TestOrderPendingTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderPending, kindOrder, orderPayload{})
}

func TestOrderSubmittedTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderSubmitted, kindOrder, orderPayload{})
}

func TestOrderSentTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderSent, kindOrder, orderPayload{})
}

func TestOrderAcceptedTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderAccepted, kindOrder, orderPayload{})
}

func TestOrderPartiallyFilledTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderPartiallyFilled, kindOrder, orderPayload{})
}

func TestOrderFilledTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderFilled, kindOrder, orderPayload{})
}

func TestOrderCanceledTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderCanceled, kindOrder, orderPayload{})
}

func TestOrderExpiredTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderExpired, kindOrder, orderPayload{})
}

func TestOrderPendingCancelTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderPendingCancel, kindOrder, orderPayload{})
}

func TestOrderRejectedTradehook(t *testing.T) {
	authInit()
	validateMethod(t, OrderRejected, kindOrder, orderPayload{})
}
