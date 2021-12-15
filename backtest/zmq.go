package backtest

import (
	"encoding/json"
	"gopkg.in/zeromq/goczmq.v4"
)

type ZmqConn struct {
	req *goczmq.Sock
}

func NewZmq(socketUrl string) (*ZmqConn, error) {

	// Create a new Req socket and connect it to the router.
	req, err := goczmq.NewReq(socketUrl)
	if err != nil {
		return nil, err
	}

	return &ZmqConn{req}, nil
}

func (z *ZmqConn) SendMsg(msg []byte) error {
	err := z.req.SendFrame(msg, goczmq.FlagNone)

	return err
}

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

func (z *ZmqConn) ReceiveMsg() ([][]byte, error) {
	msg, err := z.req.RecvMessage()
	return msg, err
}

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
