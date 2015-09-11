package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gamegos/scotty/storage"
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

type jsonResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

var (
	testServer *Server
)

func init() {
	stg := storage.NewMemStorage()
	testServer = Init(stg)
}

func apiCall(method string, urlStr string, bodyStr string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, urlStr, strings.NewReader(bodyStr))

	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	testServer.router.ServeHTTP(w, req)

	return w, nil
}

func TestCreateApp(t *testing.T) {
	postBody := initialData
	res, err := apiCall("POST", "/apps", postBody)

	if err != nil {
		t.Error(err)
	}

	if res.Code != http.StatusCreated {
		t.Error("App could not be created.", res.Code, res.Body)
	}
}

func TestUpdateApp(t *testing.T) {
	postBody := updatedData
	res, err := apiCall("PUT", "/apps/"+appID, postBody)

	if err != nil {
		t.Error(err)
	}

	if res.Code != http.StatusOK {
		t.Error("App could not be updated.")
	}
}

func TestGetApp(t *testing.T) {
	res, err := apiCall("GET", "/apps/"+appID, "")

	if err != nil {
		t.Error(err)
	}

	var response jsonResponse

	decoder := json.NewDecoder(res.Body)

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
	postBody := `{"subscriberId": "randomSubId", "platform": "gcm", "token": "foo123"}`
	res, err := apiCall("POST", "/apps/"+appID+"/devices", postBody)

	if err != nil {
		t.Error(err)
	}

	if res.Code != http.StatusCreated {
		t.Error("Subscriber device could not be added.")
	}
}

func TestAddChannel(t *testing.T) {
	postBody := `{"id": "` + channelID + `"}`
	res, err := apiCall("POST", "/apps/"+appID+"/channels", postBody)

	if err != nil {
		t.Error(err)
	}

	if res.Code != http.StatusCreated {
		t.Error("Channel could not be added.")
	}
}

func TestAddSubscriber(t *testing.T) {
	postBody := `{"subscribers": ["foo", "bar"]}`
	res, err := apiCall("POST", "/apps/"+appID+"/channels/"+channelID+"/subscribers", postBody)

	if err != nil {
		t.Error(err)
	}

	if res.Code != http.StatusCreated {
		t.Error("Subscriber device could not be added.")
	}
}

func TestDeleteChannel(t *testing.T) {
	res, err := apiCall("DELETE", "/apps/"+appID+"/channels/"+channelID, "")

	if err != nil {
		t.Error(err)
	}

	if res.Code != http.StatusOK {
		t.Error("Channel could not be deleted.")
	}
}
