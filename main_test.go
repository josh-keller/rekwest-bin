package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExampleTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

func (suite *ExampleTestSuite) SetupTest() {
	resetDB()
}

func (suite *ExampleTestSuite) TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)
	handler.ServeHTTP(rr, req)

	expected := "<h1>Welcome to Rekwest Bin</h1><form method='POST' action='/r/'><button type='submit'>Create a bin</button></form>"
	assert.Equal(t, rr.Body.String(), expected, "should return HTML home page")
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}
