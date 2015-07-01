package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.fixb.com/mir/push/storage"
)

type Handlers struct {
	stg *storage.Storage
}

func (hnd *Handlers) CreateApp(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	var app storage.App

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	str, _ := json.Marshal(app)
	hnd.stg.CreateApp(app.Id, string(str))

	res.WriteSuccess(w, 201, "")
}

func (hnd *Handlers) UpdateApp(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appID := vars["appId"]

	if hnd.stg.AppExists(appID) {
		var app storage.App
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&app); err != nil {
			res.WriteError(w, 400, err.Error())
			return
		}

		str, _ := json.Marshal(app)
		hnd.stg.CreateApp(app.Id, string(str))
	}

	res.WriteSuccess(w, 200, "")
}

func (hnd *Handlers) GetApp(w http.ResponseWriter, r *http.Request) {

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

func (hnd *Handlers) AddDevice(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appID := vars["appId"]

	var postData storage.AddDeviceRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&postData); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	if hnd.stg.AppExists(appID) {
		device := storage.Device{
			Platform:  postData.Platform,
			Token:     postData.Token,
			CreatedAt: int(time.Now().Unix()),
		}
		err := hnd.stg.AddSubscriberDevice(appID, postData.SubscriberId, &device)

		if err != nil {
			res.WriteError(w, 500, err.Error())
			return
		}
	} else {
		res.WriteError(w, 400, "App not found")
		return
	}

	res.WriteSuccess(w, 201, "")
}

func (hnd *Handlers) AddSubscriber(w http.ResponseWriter, r *http.Request) {

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

	if hnd.stg.AppExists(appID) {
		hnd.stg.AddSubscriber(appID, channelID, f.SubscriberIds)
	} else {
		res.WriteError(w, 400, "App not found")
		return
	}

	res.WriteSuccess(w, 201, "")
}

func (hnd *Handlers) AddChannel(w http.ResponseWriter, r *http.Request) {
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

	if hnd.stg.AppExists(appID) {
		hnd.stg.AddChannel(appID, m["id"].(string))
	} else {
		res.WriteError(w, 400, "App not found")
		return
	}

	res.WriteSuccess(w, 201, "")
}

func (hnd *Handlers) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appID := vars["appId"]
	channelID := vars["channelId"]

	hnd.stg.DeleteChannel(appID, channelID)

	res.WriteSuccess(w, 200, "")
}
