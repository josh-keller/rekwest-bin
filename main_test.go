package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestMain(t *testing.T) {
	is := is.New(t)

	srv, err := newServer()
	is.NoErr(err) // newServer error

	hsrv := httptest.NewServer(srv)
	defer hsrv.Close()

	req, err := http.NewRequest("get", hsrv.URL+"/", nil)
	is.NoErr(err) // NewRequest error
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	resp, err := client.Do(req)
	is.NoErr(err) // Client request error
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	is.NoErr(err)                                                                       // Body read error
	is.True(strings.Contains(string(b), `<h1>Welcome`))                                 // Contains welcome
	is.True(strings.Contains(string(b), `<button type='submit'>Create a bin</button>`)) // Contains button
}
