package protocols

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"strconv"

	"github.com/facebookgo/ensure"
	"github.com/millisecond/huic/shared"
)

var httpHandler = CreateHTTPHandler("testing", "./testcontent/")

func TestBasicHTTP(t *testing.T) {
	s := httptest.NewServer(httpHandler)
	defer s.Close()

	resp, err := http.Get(s.URL + "/text.txt")
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp.StatusCode, 200)

	body, err := ioutil.ReadAll(resp.Body)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, string(body), "some text")
}

func TestHUICUpgrade(t *testing.T) {
	s := httptest.NewServer(httpHandler)
	defer s.Close()

	client := &http.Client{}

	req, err := http.NewRequest("GET", s.URL+"/text.txt", nil)
	req.Header.Add("Accept-Encoding", "huic")
	resp, err := client.Do(req)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp.StatusCode, 200)

	// it's a basic HUIC upgrade so we shouldn't have the file in the body
	body, err := ioutil.ReadAll(resp.Body)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(body), 0)

	//verify headers are set
	ensure.True(t, len(resp.Header.Get(shared.HeaderSessionID)) > 0)
	ensure.DeepEqual(t, len(resp.Header.Get(shared.HeaderSessionKey)), 64)
	fileSize, err := strconv.ParseInt(resp.Header.Get(shared.HeaderFileSizeBytes), 10, 0)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, fileSize, int64(9))
	ensure.True(t, len(resp.Header.Get(shared.HeaderUDPPort)) > 0)
}

func TestHUICUpgradeWithContent(t *testing.T) {
	s := httptest.NewServer(httpHandler)
	defer s.Close()

	client := &http.Client{}

	req, err := http.NewRequest("GET", s.URL+"/text.txt", nil)
	req.Header.Add(shared.HeaderAcceptEncoding, "huic")
	req.Header.Add(shared.HeaderUpgrade, "with-content")
	resp, err := client.Do(req)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp.StatusCode, 200)

	body, err := ioutil.ReadAll(resp.Body)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, string(body), "some text")
}
