package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestBinNew(t *testing.T) {
	is := is.New(t)

	srv, err := newServer()
	srv.db.Connect()
	defer srv.db.Disconnect()

	is.NoErr(err) // newServer error

	hsrv := httptest.NewServer(srv)
	defer hsrv.Close()

	req, err := http.NewRequest("POST", hsrv.URL+"/new/", nil)
	is.NoErr(err) // NewRequest error
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 1 * time.Minute,
	}

	resp, err := client.Do(req)
	is.NoErr(err) // Client request error

	inspectURL := hsrv.URL + "/inspect/"
	is.Equal(resp.StatusCode, 302) // Redirects
	url, _ := resp.Location()

	is.True(strings.Contains(url.String(), inspectURL)) // Redirects to inspect path
	code1 := url.String()[len(inspectURL):]

	req2, err := http.NewRequest("POST", hsrv.URL+"/new/", nil)
	is.NoErr(err) // NewRequest error

	resp2, err := client.Do(req2)
	is.NoErr(err) // Request2 error

	url2, _ := resp2.Location()
	code2 := url2.String()[len(inspectURL):]
	is.True(code1 != code2) // Unique codes on subsequent requests
}
