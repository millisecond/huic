package protocols

import (
	"context"
	"sync/atomic"
	"testing"

	"strconv"

	"github.com/facebookgo/ensure"
)

var port uint32 = 41000

type TestServer struct {
	port   uint32
	cancel context.CancelFunc
}

func (s *TestServer) stop() {
	s.cancel()
}

// Fire up a UDP listener on a unique port to allow parallel testing
func startTestServer() (server TestServer, err error) {
	testPort := atomic.AddUint32(&port, 1)
	ctx, cancel := context.WithCancel(context.Background())
	server = TestServer{
		port:   testPort,
		cancel: cancel,
	}
	err = ListenHUIC(ctx, uint(testPort))
	return
}

func testClient(port uint32) (send chan *HUICMessage, receive chan *HUICMessage, err error) {
	return client(context.Background(), ":"+strconv.Itoa(int(port)))
}

func TestTestServer(t *testing.T) {
	s, err := startTestServer()
	ensure.Nil(t, err)
	defer s.stop()

	send, receive, err := testClient(s.port)
	ensure.Nil(t, err)

	send <- PingMessage()
	response := <-receive

	ensure.NotNil(t, response)
	ensure.DeepEqual(t, response.packetType, PacketType(PONG))
}

func TestStoppingTestServer(t *testing.T) {
	s, err := startTestServer()
	ensure.Nil(t, err)
	s.stop()

	_, _, err = testClient(s.port)
	ensure.NotNil(t, err)
}
