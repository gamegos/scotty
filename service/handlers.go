package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.fixb.com/mir/push/storage"
)

// Handlers holds the handler functions to be run with different routes.
type Handlers struct {
	stg *storage.Storage
}

func (hnd *Handlers) createApp(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	var app storage.App

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	str, _ := json.Marshal(app)
	err := hnd.stg.CreateApp(app.ID, string(str))

	if err != nil {
		res.WriteError(w, 500, err.Error())
		return
	}

	res.WriteSuccess(w, 201, "")
}

func (hnd *Handlers) updateApp(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appID := vars["appId"]

	if !hnd.stg.AppExists(appID) {
		res.WriteError(w, 400, fmt.Sprintf("App %v does not exist.", appID))
		return
	}

	var app storage.App
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	if appID != app.ID {
		res.WriteError(w, 400, "AppID mismatch.")
		return
	}

	str, _ := json.Marshal(app)
	err := hnd.stg.CreateApp(app.ID, string(str))

	if err != nil {
		res.WriteError(w, 500, err.Error())
		return
	}

	res.WriteSuccess(w, 200, "")
}

func (hnd *Handlers) getApp(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appID := vars["appId"]

	appData, err := hnd.stg.GetApp(appID)

	if err != nil {
		res.WriteError(w, 404, "App not found")
		return
	}

	res.WriteSuccess(w, 200, appData)
}

func (hnd *Handlers) addDevice(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appID := vars["appId"]

	var postData storage.AddDeviceRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&postData); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	if !hnd.stg.AppExists(appID) {
		res.WriteError(w, 400, "App not found")
		return
	}

	device := storage.Device{
		Platform:  postData.Platform,
		Token:     postData.Token,
		CreatedAt: int(time.Now().Unix()),
	}
	err := hnd.stg.AddSubscriberDevice(appID, postData.SubscriberID, &device)

	if err != nil {
		res.WriteError(w, 500, err.Error())
		return
	}

	res.WriteSuccess(w, 201, "")
}

func (hnd *Handlers) addSubscriber(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appID := vars["appId"]
	channelID := vars["channelId"]

	var f storage.AddSubscriberRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	if !hnd.stg.AppExists(appID) {
		res.WriteError(w, 400, "App not found")
		return
	}

	err := hnd.stg.AddSubscriber(appID, channelID, f.SubscriberIds)

	if err != nil {
		res.WriteError(w, 500, err.Error())
		return
	}

	res.WriteSuccess(w, 201, "")
}

func (hnd *Handlers) addChannel(w http.ResponseWriter, r *http.Request) {
	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appID := vars["appId"]

	var f interface{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	m := f.(map[string]interface{})

	if !hnd.stg.AppExists(appID) {
		res.WriteError(w, 400, "App not found")
		return
	}

	err := hnd.stg.AddChannel(appID, m["id"].(string))

	if err != nil {
		res.WriteError(w, 500, err.Error())
		return
	}

	res.WriteSuccess(w, 201, "")
}

func (hnd *Handlers) deleteChannel(w http.ResponseWriter, r *http.Request) {
	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appID := vars["appId"]
	channelID := vars["channelId"]

	err := hnd.stg.DeleteChannel(appID, channelID)

	if err != nil {
		res.WriteError(w, 500, err.Error())
		return
	}

	res.WriteSuccess(w, 200, "")
}
