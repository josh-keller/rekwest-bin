package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestBinInspect(t *testing.T) {
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
		Timeout: 1 * time.Minute,
	}

	resp, err := client.Do(req)
	is.NoErr(err)                  // Client request error
	is.Equal(resp.StatusCode, 200) // Making new bin gives 200

	inspectURL := hsrv.URL + "/inspect/"
	binID := resp.Request.URL.String()[len(inspectURL):]

	var reqBody = []byte(`{"request":"Golang test request"}`)
	req, err = http.NewRequest("POST", hsrv.URL+"/r/"+binID, bytes.NewBuffer(reqBody))
	is.NoErr(err) // New Request error
	req.Header.Set("X-Custom-Header", "custom value 2468")
	req.Header.Set("Content-type", "application/json")
	resp, err = client.Do(req)
	is.NoErr(err)                  // Client request error
	is.Equal(resp.StatusCode, 200) // Initial request returns 200

	req, err = http.NewRequest("GET", inspectURL+binID, nil)
	is.NoErr(err) // New Request error
	resp, err = client.Do(req)
	is.NoErr(err)                  // Client request error
	is.Equal(resp.StatusCode, 200) // Inspect returns 200
	body, err := ioutil.ReadAll(resp.Body)
	is.NoErr(err) // body read error
	is.True(strings.Contains(string(body), "custom value 2468"))
	is.True(strings.Contains(string(body), "Golang test request"))
}
