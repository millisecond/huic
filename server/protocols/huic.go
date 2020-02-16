package protocols

import (
	"context"
	"fmt"
	"net"
	"time"

	"bytes"
	"io"
)

const maxBufferSize = 1024
const timeout = time.Second

// server wraps all the UDP echo server functionality.
// ps.: the server is capable of answering to a single
// client at a time.
func ListenHUIC(ctx context.Context, port uint) (err error) {
	address := fmt.Sprintf(":%d", port)

	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("HUIC Server UDP packet-received: bytes=%d from=%s\n", n, addr.String())

			// configure the send timeout, waiting for write queue
			deadline := time.Now().Add(timeout)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			// Echo the packet's contents back to the client.
			n, err = pc.WriteTo(PongMessage().Bytes(), addr)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("HUIC Server UDP packet-written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("HUIC Server: Done", ctx.Err())
		case err = <-doneChan:
			fmt.Println("HUIC Server Error: ", err)
		}
		pc.Close()
	}()

	return
}

// client wraps the whole functionality of a UDP client that sends
// a message and waits for a response coming back from the server
// that it initially targetted.
func client(ctx context.Context, address string) (send chan *HUICMessage, receive chan *HUICMessage, err error) {
	// Resolve the UDP address so that we can make use of DialUDP
	// with an actual IP and port instead of a name (in case a
	// hostname is specified).
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	// Although we're not in a connection-oriented transport,
	// the act of `dialing` is analogous to the act of performing
	// a `connect(2)` syscall for a socket of type SOCK_DGRAM:
	// - it forces the underlying socket to only read and write
	//   to and from a specific remote address.
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return
	}

	doneChan := make(chan error, 1)
	send = make(chan *HUICMessage, 100)
	receive = make(chan *HUICMessage, 100)

	go func() {
		// It is possible that this action blocks, although this
		// should only occur in very resource-intensive situations:
		// - when you've filled up the socket buffer and the OS
		//   can't dequeue the queue fast enough.
		m := <-send
		n, err := io.Copy(conn, bytes.NewReader(m.Bytes()))
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("HUIC Client packet-written: bytes=%d\n", n)

		buffer := make([]byte, maxBufferSize)

		// Set a deadline for the ReadOperation so that we don't
		// wait forever for a server that might not respond on
		// a resonable amount of time.
		deadline := time.Now().Add(timeout)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			doneChan <- err
			return
		}

		nRead, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("HUIC Client packet-received: bytes=%d from=%s\n", nRead, addr.String())

		receive <- HUICMessageFromBytes(buffer[:nRead])
	}()

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("HUIC Client: Done", ctx.Err())
		case err = <-doneChan:
			fmt.Println("HUIC Client Error: ", err)
		}
		close(send)
		close(receive)
		conn.Close()
	}()

	return
}
