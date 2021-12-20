package backtest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/zeromq/goczmq.v4"
	"log"
	"testing"
)

var socketUrl = "tcp://127.0.0.1:3006"
var routerUrl = "tcp://*:3006"

func createZMQRouter() (*goczmq.Poller, *goczmq.Sock) {
	sock := goczmq.NewSock(goczmq.Router)
	if err := sock.Attach(routerUrl, true); err != nil {
		log.Fatal(err)
	}
	router, err := goczmq.NewPoller(sock)
	if err != nil {
		log.Fatal(err)
	}
	return router, sock
}

func TestSendJZMQMessageAndRetrieveResponse(t *testing.T) {

	// Create client
	clientZMQ, err := NewZmq(socketUrl)
	if err != nil {
		assert.Error(t, err)
	}

	// Create server
	routerZMQ, sock := createZMQRouter()
	defer sock.Destroy()

	// Send message from client to server
	err = clientZMQ.SendMsg([]byte("hello"))
	if err != nil {
		assert.Error(t, err)
	}

	// Router receive message
	hitSock := routerZMQ.Wait(3000)

	msg, err := hitSock.RecvMessage()
	if err != nil {
		assert.Error(t, err)
	}
	fmt.Printf("router received from '%v'", msg)

	fmt.Print("test", string(msg[0]), "\n")
	fmt.Print("test 2", string(msg[1]), "\n")
	assert.Equal(t, "", string(msg[0]), "invalid client message")

	// Router sent response to client
	err = hitSock.SendFrame([]byte("world"), goczmq.FlagNone)
	if err != nil {
		assert.Error(t, err)
	}

	// Receive message from server
	msg, err = clientZMQ.ReceiveMsg()
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, "", string(msg[0]), "invalid server message")
}

// TODO
func TestSendZMQMessageJSONAndRetrieveResponse(t *testing.T) {

}
