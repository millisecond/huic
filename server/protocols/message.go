package protocols

import (
	"encoding/binary"
	"errors"
	"strconv"
)

// PacketType is the first byte of every HUIC message.
type PacketType byte

const (
	PING = 1 + iota
	PONG
)

const LargeHeaderByteCount = 21

var Endian = binary.BigEndian

// HUICLargeHeader is the standard, include everything fully, header.
//
// The packetType is controlled by the HUICMessage so it can select between header types.
type HUICLargeHeader struct {
	connectionID    uint64
	packetNumber    uint64
	version         uint32
	intendedDelayMS uint8
}

func HUICLargeHeaderFromBytes(b []byte) (h *HUICLargeHeader, err error) {
	if len(b) != LargeHeaderByteCount {
		err = errors.New("Invalid byte count to parse a header: " + strconv.Itoa(len(b)))
		return
	}
	h = &HUICLargeHeader{
		connectionID:    Endian.Uint64(b[0:8]),
		packetNumber:    Endian.Uint64(b[8:16]),
		version:         Endian.Uint32(b[16:20]),
		intendedDelayMS: b[20],
	}
	return
}

func (h *HUICLargeHeader) ToBytes() (b []byte) {
	b = make([]byte, LargeHeaderByteCount)
	//connectionIDBuf := make([]byte, 8)
	//versionBuf := make([]byte, 4)
	//packetNumberBuf := make([]byte, 8)
	Endian.PutUint64(b, h.connectionID)
	Endian.PutUint64(b[8:], h.packetNumber)
	Endian.PutUint32(b[16:], h.version)
	b[20] = h.intendedDelayMS
	return
}

// HUICMessage is the standard, include everything, message.
type HUICMessage struct {
	packetType PacketType
	body       []byte
}

func (m *HUICMessage) Bytes() []byte {
	b := make([]byte, 1)
	b[0] = byte(m.packetType)
	return b
}

func HUICMessageFromBytes(b []byte) (msg *HUICMessage) {
	msg = &HUICMessage{
		packetType: PacketType(b[0]),
	}
	return
}

func PingMessage() (msg *HUICMessage) {
	msg = &HUICMessage{
		packetType: PING,
	}
	return
}

func PongMessage() (msg *HUICMessage) {
	msg = &HUICMessage{
		packetType: PONG,
	}
	return
}
