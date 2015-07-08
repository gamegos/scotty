package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"gitlab.fixb.com/mir/push/storage"
)

var appID, initialData, updatedData, channelID string
var respRec *httptest.ResponseRecorder
var router *mux.Router

func setup() {
	appID = "testapp"
	initialData = `{
	    "id": "` + appID + `",
	    "platforms": {
	        "apns": {
	            "certificate": "apnscertificate",
	            "privateKey": "privatekey"
	        },
	        "gcm": {
	            "projectId": "projectid",
	            "apiKey": "apikey"
	        }
	    }
	}`
	updatedData = `
	{
	    "id": "` + appID + `",
	    "platforms": {
	        "apns": {
	            "certificate": "updatedapnscertificate",
	            "privateKey": "updatedprivatekey"
	        },
	        "gcm": {
	            "projectId": "updatedprojectid",
	            "apiKey": "updatedapikey"
	        }
	    }
	}`
	channelID = "someRandomchannelID"

	respRec = httptest.NewRecorder()
	confFile := "../default.conf"

	conf := storage.InitConfig(confFile)
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

	var response WrappedResponse

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

	if response.Data.(string) != string(updatedAppStr) {
		t.Error("Received app data is not the same as updated data.")
	}
}

func TestAddDevice(t *testing.T) {
	setup()

	postBody := strings.NewReader(`{"subscriberId": "randomSubId", "platform": "apns", "token": "foo123"}`)
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
