package backtest

import (
	"encoding/json"
	"gopkg.in/zeromq/goczmq.v4"
)

type ZmqConn struct {
	req *goczmq.Sock
}

// NewZmq create new Req socket and connect it to the router
func NewZmq(socketUrl string) (*ZmqConn, error) {
	req, err := goczmq.NewReq(socketUrl)
	if err != nil {
		return nil, err
	}

	return &ZmqConn{req}, nil
}

// SendMsg sends a byte array via the socket
func (z *ZmqConn) SendMsg(msg []byte) error {
	err := z.req.SendFrame(msg, goczmq.FlagNone)

	return err
}

// SendJSON convert data to json and sends a byte array via the socket
func (z *ZmqConn) SendJSON(src interface{}) error {
	srcJSON, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = z.SendMsg(srcJSON)
	if err != nil {
		return err
	}
	return nil
}

// ReceiveMsg receives a full message from the socket and returns it as an array of byte arrays
func (z *ZmqConn) ReceiveMsg() ([][]byte, error) {
	msg, err := z.req.RecvMessage()
	return msg, err
}

// ReceiveJSON receives a full message from the socket and parse JSON-encoded data into selected struct
func (z *ZmqConn) ReceiveJSON(dst interface{}) error {
	msg, err := z.req.RecvMessage()
	if err != nil {
		return err
	}

	err = json.Unmarshal(msg[0], dst)
	if err != nil {
		return err
	}
	return nil
}

func (z *ZmqConn) Close() {
	z.req.Destroy()
}
