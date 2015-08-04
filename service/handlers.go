package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gamegos/jsend"
	"github.com/gorilla/mux"
	"gitlab.fixb.com/mir/push/storage"
)

// Handlers holds the handler functions to be run with different routes.
type Handlers struct {
	stg storage.Storage
}

type failResponse struct {
	Message string `json:"message"`
}

func (hnd *Handlers) getHealth(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, []byte(nil), 400)
}

func (hnd *Handlers) createApp(w http.ResponseWriter, r *http.Request) {

	var app storage.App

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		jsend.Fail(w, &failResponse{Message: err.Error()}, 400)
		return
	}

	str, _ := json.Marshal(app)
	err := hnd.stg.CreateApp(app.ID, string(str))

	if err != nil {
		jsend.Error(w, err.Error(), 500)
		return
	}

	jsend.Success(w, []byte(nil), 201)
}

func (hnd *Handlers) updateApp(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appID := vars["appId"]

	if !hnd.stg.AppExists(appID) {
		jsend.Fail(w, &failResponse{Message: fmt.Sprintf("App %v does not exist.", appID)}, 400)
		return
	}

	var app storage.App
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		jsend.Fail(w, &failResponse{Message: err.Error()}, 400)
		return
	}

	if appID != app.ID {
		jsend.Fail(w, &failResponse{Message: "AppID mismatch"}, 400)
		return
	}

	str, _ := json.Marshal(app)
	err := hnd.stg.CreateApp(app.ID, string(str))

	if err != nil {
		jsend.Error(w, err.Error(), 500)
		return
	}

	jsend.Success(w, []byte(nil), 200)
}

func (hnd *Handlers) getApp(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appID := vars["appId"]

	appData, err := hnd.stg.GetApp(appID)

	if err != nil {
		jsend.Fail(w, &failResponse{Message: "App not found."}, 404)
		return
	}

	jsend.Success(w, appData, 200)
}

func (hnd *Handlers) addDevice(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appID := vars["appId"]

	var postData storage.AddDeviceRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&postData); err != nil {
		jsend.Fail(w, &failResponse{Message: err.Error()}, 400)
		return
	}

	if !hnd.stg.AppExists(appID) {
		jsend.Fail(w, &failResponse{Message: "App not found."}, 400)
		return
	}

	device := storage.Device{
		Platform:  postData.Platform,
		Token:     postData.Token,
		CreatedAt: int(time.Now().Unix()),
	}
	err := hnd.stg.AddSubscriberDevice(appID, postData.SubscriberID, &device)

	if err != nil {
		jsend.Error(w, err.Error(), 500)
		return
	}

	jsend.Success(w, []byte(nil), 201)
}

func (hnd *Handlers) addSubscriber(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appID := vars["appId"]
	channelID := vars["channelId"]

	var f storage.AddSubscriberRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		jsend.Fail(w, &failResponse{Message: err.Error()}, 400)
		return
	}

	if !hnd.stg.AppExists(appID) {
		jsend.Fail(w, &failResponse{Message: "App not found."}, 400)
		return
	}

	err := hnd.stg.AddSubscriber(appID, channelID, f.SubscriberIds)

	if err != nil {
		jsend.Error(w, err.Error(), 500)
		return
	}

	jsend.Success(w, []byte(nil), 201)
}

func (hnd *Handlers) addChannel(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appID := vars["appId"]

	var f interface{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		jsend.Fail(w, &failResponse{Message: err.Error()}, 400)
		return
	}

	m := f.(map[string]interface{})

	if !hnd.stg.AppExists(appID) {
		jsend.Fail(w, &failResponse{Message: "App not found."}, 400)
		return
	}

	err := hnd.stg.AddChannel(appID, m["id"].(string))

	if err != nil {
		w.WriteHeader(500)
		jsend.Fail(w, &failResponse{Message: err.Error()}, 400)
		return
	}

	jsend.Success(w, []byte(nil), 201)
}

func (hnd *Handlers) deleteChannel(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appID := vars["appId"]
	channelID := vars["channelId"]

	err := hnd.stg.DeleteChannel(appID, channelID)

	if err != nil {
		jsend.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte{})
}
