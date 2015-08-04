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
	stg *storage.Storage
}

type failResponse struct {
	Message string `json:"message"`
}

func (hnd *Handlers) getHealth(w http.ResponseWriter, r *http.Request) {
	jsend.Wrap(w).Status(200).Send()
}

func (hnd *Handlers) createApp(w http.ResponseWriter, r *http.Request) {

	jw := jsend.Wrap(w)
	var app storage.App

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		jw.Status(400).Data(&failResponse{Message: err.Error()}).Send()
		return
	}

	str, _ := json.Marshal(app)
	err := hnd.stg.CreateApp(app.ID, string(str))

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(201).Send()
}

func (hnd *Handlers) updateApp(w http.ResponseWriter, r *http.Request) {

	jw := jsend.Wrap(w)
	vars := mux.Vars(r)
	appID := vars["appId"]

	if !hnd.stg.AppExists(appID) {
		jw.Status(400).Data(&failResponse{Message: fmt.Sprintf("App %v does not exist.", appID)}).Send()
		return
	}

	var app storage.App
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		jw.Status(400).Data(&failResponse{Message: err.Error()}).Send()
		return
	}

	if appID != app.ID {
		jw.Status(400).Data(&failResponse{Message: "AppID mismatch"}).Send()
		return
	}

	str, _ := json.Marshal(app)
	err := hnd.stg.CreateApp(app.ID, string(str))

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(200).Send()
}

func (hnd *Handlers) getApp(w http.ResponseWriter, r *http.Request) {

	jw := jsend.Wrap(w)
	vars := mux.Vars(r)
	appID := vars["appId"]

	appData, err := hnd.stg.GetApp(appID)

	if err != nil {
		jw.Status(404).Data(&failResponse{Message: "App not found."}).Send()
		return
	}

	jw.Status(200).Data(appData).Send()
}

func (hnd *Handlers) addDevice(w http.ResponseWriter, r *http.Request) {

	jw := jsend.Wrap(w)
	vars := mux.Vars(r)
	appID := vars["appId"]

	var postData storage.AddDeviceRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&postData); err != nil {
		jw.Status(400).Data(&failResponse{Message: err.Error()}).Send()
		return
	}

	if !hnd.stg.AppExists(appID) {
		jw.Status(400).Data(&failResponse{Message: "App not found."}).Send()
		return
	}

	device := storage.Device{
		Platform:  postData.Platform,
		Token:     postData.Token,
		CreatedAt: int(time.Now().Unix()),
	}
	err := hnd.stg.AddSubscriberDevice(appID, postData.SubscriberID, &device)

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(201).Send()
}

func (hnd *Handlers) addSubscriber(w http.ResponseWriter, r *http.Request) {

	jw := jsend.Wrap(w)
	vars := mux.Vars(r)
	appID := vars["appId"]
	channelID := vars["channelId"]

	var f storage.AddSubscriberRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		jw.Status(400).Data(&failResponse{Message: err.Error()}).Send()
		return
	}

	if !hnd.stg.AppExists(appID) {
		jw.Status(400).Data(&failResponse{Message: "App not found."}).Send()
		return
	}

	err := hnd.stg.AddSubscriber(appID, channelID, f.SubscriberIds)

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(201).Send()
}

func (hnd *Handlers) addChannel(w http.ResponseWriter, r *http.Request) {

	jw := jsend.Wrap(w)
	vars := mux.Vars(r)
	appID := vars["appId"]

	var f interface{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		jw.Status(400).Data(&failResponse{Message: err.Error()}).Send()
		return
	}

	m := f.(map[string]interface{})

	if !hnd.stg.AppExists(appID) {
		jw.Status(400).Data(&failResponse{Message: "App not found."}).Send()
		return
	}

	err := hnd.stg.AddChannel(appID, m["id"].(string))

	if err != nil {
		jw.Status(500).Data(&failResponse{Message: err.Error()}).Send()
		return
	}

	jw.Status(201).Send()
}

func (hnd *Handlers) deleteChannel(w http.ResponseWriter, r *http.Request) {

	jw := jsend.Wrap(w)
	vars := mux.Vars(r)
	appID := vars["appId"]
	channelID := vars["channelId"]

	err := hnd.stg.DeleteChannel(appID, channelID)

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(200).Send()
}
