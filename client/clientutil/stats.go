package clientutil

import "time"

type DownloadSpeed struct {
	start           time.Time
	end             time.Time
	inProgress      bool
	bytesDownloaded uint64
}

type UDPPacketLoss struct {
	seqReceived []uint64
}
