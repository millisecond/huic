package shared

import (
	"net/http"
	"sync"
	"time"

	"encoding/hex"
	"strconv"

	"github.com/millisecond/huic/shared/huiccrypt"
)

var liveSessions = make(map[string]*HUICSession)
var liveSessionsMutex = &sync.RWMutex{}

type HUICSession struct {
	ID  string
	Key []byte

	fileSizeBytes int64
	start         time.Time
	finish        time.Time
	finished      bool

	name    string
	udpPort int

	lock sync.RWMutex
}

func (session *HUICSession) WriteSessionHeaders(w http.ResponseWriter) {
	w.Header().Add(HeaderSessionID, session.ID)
	w.Header().Add(HeaderSessionKey, hex.EncodeToString(session.Key))
	w.Header().Add(HeaderFileSizeBytes, strconv.FormatInt(session.fileSizeBytes, 10))
	w.Header().Add(HeaderUDPPort, strconv.Itoa(session.udpPort))
}

func (session *HUICSession) Stop() {
	liveSessionsMutex.Lock()
	defer liveSessionsMutex.Unlock()
	delete(liveSessions, session.ID)
}

func StartSession(f http.File) (session *HUICSession, err error) {
	uuid, err := huiccrypt.PseudoUUID()
	if err != nil {
		return
	}
	key, err := huiccrypt.AESKey()
	if err != nil {
		return
	}
	session = &HUICSession{}
	session.ID = uuid
	session.Key = key
	session.udpPort = 41000
	session.start = time.Now()
	session.finished = false

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	session.name = stat.Name()
	session.fileSizeBytes = stat.Size()
	return session, nil
}

func SessionFromHTTP(req *http.Request) *HUICSession {
	sessionID := req.Header.Get(HeaderSessionID)
	if len(sessionID) == 0 {
		return nil
	}

	liveSessionsMutex.RLock()
	defer liveSessionsMutex.RUnlock()
	sess, ok := liveSessions[sessionID]
	if ok {
		return sess
	}

	return nil
}
