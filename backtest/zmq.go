package backtest

import (
	"encoding/json"
	"gopkg.in/zeromq/goczmq.v4"
	"log"
)

type ZmqConn struct {
	req *goczmq.Sock
}

func NewZmq(socketUrl string) *ZmqConn {

	// Create a new Req socket and connect it to the router.
	req, err := goczmq.NewReq(socketUrl)
	if err != nil {
		// TODO
		log.Fatal(err)
	}

	return &ZmqConn{
		req,
	}
}

func (z *ZmqConn) SendMsg(msg []byte) error {
	err := z.req.SendFrame(msg, goczmq.FlagNone)

	return err
}

func (z *ZmqConn) SendJSON(src interface{}) error {

	// Convert to json
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
