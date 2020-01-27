package protocols

import (
	"log"
	"net/http"

	"context"

	"strings"

	"github.com/millisecond/huic/server/serverutil"
	"github.com/millisecond/huic/shared"
)

func ServeHTTP(address string, directory string) {
	srv := &http.Server{
		Handler: CreateHTTPHandler(address, directory),
		Addr:    address,
	}

	log.Printf("HTTP - START - LISTENING ON %s", address)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("HTTP - START - ERROR: %v", err)
	}
}

func CreateHTTPHandler(address string, directory string) http.Handler {
	dir := http.Dir(directory)
	return &HTTPHandler{
		address:    address,
		dir:        dir,
		fileServer: http.FileServer(dir),
	}
}

type HTTPHandler struct {
	address    string
	dir        http.Dir
	fileServer http.Handler
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Logging and context wrapper.
	ctx := context.WithValue(req.Context(), "address", h.address)
	req = req.WithContext(ctx)
	if serverutil.OnDev {
		path := ""
		if req.URL != nil {
			path = req.URL.String()
		}
		errorHeader := w.Header().Get("Error-Code")
		log.Printf("HTTP - %s %s - error=%s", req.Method, path, errorHeader)
	}
	h.ProcessHTTPRequest(w, req)
}

func (h *HTTPHandler) ProcessHTTPRequest(w http.ResponseWriter, req *http.Request) {
	path := ""
	if req.URL != nil {
		path = req.URL.String()
	}

	switch req.Method {
	case shared.GETMethod:
		acceptEncoding := strings.Split(req.Header.Get(shared.HeaderAcceptEncoding), ",")
		if shared.StringArrayContains(acceptEncoding, "huic") {
			// we can start a HUIC request
			log.Printf("HUIC - GET - %s", path)

			f, err := h.dir.Open(path)
			if err != nil {
				w.Header().Add("Error-Code", err.Error())
				w.WriteHeader(500)
				return
			}

			s, err := shared.StartSession(f)
			if err != nil {
				w.Header().Add("Error-Code", err.Error())
				w.WriteHeader(500)
				return
			}
			s.WriteSessionHeaders(w)
			if req.Header.Get(shared.HeaderUpgrade) == "with-content" {
				h.fileServer.ServeHTTP(w, req)
			} else {
				w.WriteHeader(200)
			}
		} else {
			// not an upgrade, just serve the file over HTTP
			h.fileServer.ServeHTTP(w, req)
		}
	case shared.HUICSTATUSMethod:
		log.Printf("HUIC - STATUS - %s", path)
		session := shared.SessionFromHTTP(req)
		if session == nil {
			w.WriteHeader(404)
			return
		}
		session.WriteSessionHeaders(w)
		w.WriteHeader(200)
	case shared.HUICSTOPMethod:
		log.Printf("HUIC - STOP - %s", path)
		session := shared.SessionFromHTTP(req)
		if session == nil {
			w.WriteHeader(404)
			return
		}
		session.Stop()
		session.WriteSessionHeaders(w)
		w.WriteHeader(200)
	default:
		log.Printf("HTTP - UNKNOWN METHOD - %s %s", req.Method, path)
		w.WriteHeader(500)
	}
}
