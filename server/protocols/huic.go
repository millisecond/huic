package protocols

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/millisecond/huic/shared"
)

func ListenHUIC(ports *shared.Range) {

}

const maxBufferSize = 1024
const timeout = time.Second

// server wraps all the UDP echo server functionality.
// ps.: the server is capable of answering to a single
// client at a time.
func listen(ctx context.Context, port uint) (err error) {
	address := fmt.Sprintf(":%d", port)

	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}
	defer pc.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("UDP packet-received: bytes=%d from=%s\n", n, addr.String())

			// configure the send timeout, waiting for write queue
			deadline := time.Now().Add(timeout)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			// Echo the packet's contents back to the client.
			n, err = pc.WriteTo(buffer[:n], addr)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("UDP packet-written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}
