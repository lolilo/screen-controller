package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"net/url"
)

type PostURLRequest struct {
	DestinationSlaveName string
	URLToLoadInBrowser string
}

func TestFormHandler(t *testing.T) {
	assert.Equal(t, 200, sendGetToFormHandler("/"))
	assert.Equal(t, 301, sendGetToFormHandler("addfs"))
}

func parseJsonSlave(input []byte) (slave PostURLRequest) {
	err := json.Unmarshal(input, &slave)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func parseJsonReply(input []byte) (reply Reply) {
	err := json.Unmarshal(input, &reply)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func sendGetToFormHandler(URL string) (int) {
	TEMPLATE_PATH="templates/"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		request.URL.Path = URL
		formHandler(w,request)
	}))

	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)

	return resp.StatusCode
}

func TestSetDefaultMasterAddress(t *testing.T) {
	defaultUrl := setMasterAddress()

	assert.Equal(t, "http://localhost:5000", defaultUrl)
}

func TestSubmitHandler(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		submitHandler(w,request)
	}))

	client := &http.Client{}
	resp, _ := client.PostForm(testServer.URL, url.Values{"slave-id": {"1"}, "url": {"http://www.google.com"}})

	POSTRequestBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	reply := parseJsonReply(POSTRequestBody).HTML

	assert.Equal(t, true, strings.Contains(reply, "http://www.google.com"))
}

func TestSendConfirmationMessageToUser(t *testing.T) {
	TEMPLATE_PATH="templates/"
	var responseHeader string
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		sendConfirmationMessageToUser(w,"aaaa", "bbbb", "cccc")
		responseHeader=w.Header().Get("Content-Type")
	}))
	client := &http.Client{}
	resp, _ := client.Get(testServer.URL)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", responseHeader)
	POSTRequestBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	reply := parseJsonReply(POSTRequestBody).HTML

	assert.Equal(t, true, strings.Contains(reply, "aaaa"))
	assert.Equal(t, true, strings.Contains(reply, "bbbb"))
	assert.Equal(t, true, strings.Contains(reply, "cccc"))
}

func TestStatusMessageForAvailableSever(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))
	statusMessage := returnStatusMessageFrom(testServer.URL)

	assert.Equal(t, "Success!", statusMessage)
}

func TestStatusMessageForUnavailableServer(t *testing.T) {
	statusMessage := returnStatusMessageFrom("")

	assert.Equal(t, "URL cannot be open :( (HTTP status code 0)", statusMessage)
}

func TestCheckStatusCode(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	responseStatusCode := checkStatusCode(testServer.URL)
	assert.Equal(t, 200, responseStatusCode)
	responseStatusCode = checkStatusCode("")
	assert.Equal(t, 0, responseStatusCode)
}

func TestIfUrlIsValid(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
	}))

	assert.True(t, isUrlValid(testServer.URL))
}

func TestIfUrlIsInvalid(t *testing.T) {
	assert.False(t, isUrlValid(""))
}

func TestSendUrlAndIdToMaster(t *testing.T) {
	var numberOfMessagesSent = 0
	var url = ""
	var id = ""

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		numberOfMessagesSent++
		POSTRequestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		slave := parseJsonSlave(POSTRequestBody)
		url = slave.URLToLoadInBrowser
		id = slave.DestinationSlaveName
	}))
	
	sendUrlAndIdToMaster(testServer.URL, "http://index.hu", "2")
	assert.Equal(t, 1, numberOfMessagesSent)
	assert.Equal(t, "http://index.hu", url)
	assert.Equal(t, "2", id)
}
func TestReceiveAndMapSlaveAddress(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		receiveAndMapSlaveAddress(w,request)
	}))

	client := &http.Client{}
	var testIdList IdList
	testIdList.Id = append(testIdList.Id, "testSlaveId")
	jsonMessage, _ := json.Marshal(testIdList)
	client.Post(testServer.URL, "application/json", strings.NewReader(string(jsonMessage)))

	assert.Equal(t, testIdList, id_list)
}

func TestConfirmationMessage(t *testing.T) {
	TEMPLATE_PATH="templates/"
	answer_string := parseJsonReply(confirmationMessage("aaaa", "bbbb", "cccc")).HTML
	assert.Equal(t, true, strings.Contains(answer_string, "aaaa"))
	assert.Equal(t, true, strings.Contains(answer_string, "bbbb"))
	assert.Equal(t, true, strings.Contains(answer_string, "cccc"))
}
