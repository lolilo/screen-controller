package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestJsonCanBeParsed(t *testing.T) {
	var inputJson = []byte(`{"ID":"LeftScreen","URL":"http://google.com"}`)

	parsedJson := parseJson(inputJson)

	assert.Equal(t, "LeftScreen", parsedJson.ID)
	assert.Equal(t, "http://google.com", parsedJson.URL)
}

func TestUrlValueMessageIsSent(t *testing.T) {
	var numberOfMessagesSent = 0
	var url = ""

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
		url = request.PostFormValue("url")
	}))

	sendUrlValueMessageToServer(testServer.URL, "http://index.hu")

	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
}

func TestSlaveIpMapIsInitialized(t *testing.T) {
	slaveIPMap := initializeSlaveIPs()

	assert.Equal(t, "http://10.0.0.42:8080", slaveIPMap["1"])
	assert.Equal(t, "http://10.0.0.231:8080", slaveIPMap["2"])
}

func TestDestinationUrl(t *testing.T) {
	slaveIPMap := initializeSlaveIPs()
	destinationURL := destinationUrl("1", slaveIPMap)

	assert.Equal(t, "http://10.0.0.42:8080", destinationURL)
}