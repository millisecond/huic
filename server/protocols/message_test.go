package protocols

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestHUICLargeHeader(t *testing.T) {
	uint64Tests := []uint64{0, 1, 128, ^uint64(0)}
	uint32Tests := []uint32{0, 1, 128, ^uint32(0)}
	uint8Tests := []uint8{0, 1, 128, ^uint8(0)}

	for connID := range uint64Tests {
		for packetNumber := range uint64Tests {
			for version := range uint32Tests {
				for delay := range uint8Tests {
					h := &HUICLargeHeader{
						connectionID:    uint64(connID),
						packetNumber:    uint64(packetNumber),
						version:         uint32(version),
						intendedDelayMS: uint8(delay),
					}
					b := h.ToBytes()
					decoded, err := HUICLargeHeaderFromBytes(b)
					ensure.Nil(t, err)
					ensure.DeepEqual(t, decoded.connectionID, uint64(connID))
					ensure.DeepEqual(t, decoded.packetNumber, uint64(packetNumber))
					ensure.DeepEqual(t, decoded.version, uint32(version))
					ensure.DeepEqual(t, decoded.intendedDelayMS, uint8(delay))
				}
			}
		}
	}
}
