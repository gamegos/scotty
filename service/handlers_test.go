package service

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gamegos/scotty/storage"
	"github.com/gorilla/mux"
)

var appID = "testapp"
var channelID = "someRandomchannelID"
var initialData = `{
		"id": "` + appID + `",
		"gcm": {
				"projectId": "projectid",
				"apiKey": "apikey"
		}
}`
var updatedData = `
{
		"id": "` + appID + `",
		"gcm": {
				"projectId": "updatedprojectid",
				"apiKey": "updatedapikey"
		}
}`
var confFile = flag.String("config", "../default.conf", "Config file")
var respRec *httptest.ResponseRecorder
var router *mux.Router

type jsonResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

func setup() {

	respRec = httptest.NewRecorder()

	conf := storage.InitConfig(*confFile)
	stg := storage.Init(&conf.Redis)

	router = NewRouter(stg)
}

func TestCreateApp(t *testing.T) {
	setup()

	postBody := strings.NewReader(initialData)
	req, err := http.NewRequest("POST", "/apps", postBody)

	if err != nil {
		t.Error(err)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusCreated {
		t.Error("App could not be created.", respRec.Code, respRec.Body)
	}
}

func TestUpdateApp(t *testing.T) {
	setup()

	postBody := strings.NewReader(updatedData)
	req, err := http.NewRequest("PUT", "/apps/"+appID, postBody)

	if err != nil {
		t.Error(err)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusOK {
		t.Error("App could not be updated.")
	}
}

func TestGetApp(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/apps/"+appID, nil)

	if err != nil {
		t.Error(err)
	}

	router.ServeHTTP(respRec, req)

	var response jsonResponse

	decoder := json.NewDecoder(respRec.Body)

	if err := decoder.Decode(&response); err != nil {
		t.Error(err)
		return
	}

	var app storage.App

	decoder = json.NewDecoder(strings.NewReader(updatedData))

	if err := decoder.Decode(&app); err != nil {
		t.Error(err)
		return
	}

	updatedAppStr, _ := json.Marshal(app)

	if string(response.Data) != string(updatedAppStr) {
		t.Error("Received app data is not the same as updated data.")
	}
}

func TestAddDevice(t *testing.T) {
	setup()

	postBody := strings.NewReader(`{"subscriberId": "randomSubId", "platform": "gcm", "token": "foo123"}`)
	req, err := http.NewRequest("POST", "/apps/"+appID+"/devices", postBody)

	if err != nil {
		t.Error(err)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusCreated {
		t.Error("Subscriber device could not be added.")
	}
}

func TestAddChannel(t *testing.T) {
	setup()

	postBody := strings.NewReader(`{"id": "` + channelID + `"}`)
	req, err := http.NewRequest("POST", "/apps/"+appID+"/channels", postBody)

	if err != nil {
		t.Error(err)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusCreated {
		t.Error("Channel could not be added.")
	}
}

func TestAddSubscriber(t *testing.T) {
	setup()

	postBody := strings.NewReader(`{"subscribers": ["foo", "bar"]}`)
	req, err := http.NewRequest("POST", "/apps/"+appID+"/channels/"+channelID+"/subscribers", postBody)

	if err != nil {
		t.Error(err)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusCreated {
		t.Error("Subscriber device could not be added.")
	}
}

func TestDeleteChannel(t *testing.T) {
	setup()

	req, err := http.NewRequest("DELETE", "/apps/"+appID+"/channels/"+channelID, nil)

	if err != nil {
		t.Error(err)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusOK {
		t.Error("Channel could not be deleted.")
	}
}
