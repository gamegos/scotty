package service

import (
	"encoding/json"
	// "fmt"
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
}

func (hnd *Handlers) UpdateApp(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	appId := vars["appId"]

	if hnd.stg.AppExists(appId) {
		hnd.CreateApp(w, r)
	}
}

func (hnd *Handlers) GetApp(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appId := vars["appId"]

	appData, err := hnd.stg.GetApp(appId)

	if err != nil {
		res.WriteError(w, 404, "App not found")
		return
	}

	res.WriteSuccess(w, 200, appData)
}

func (hnd *Handlers) AddDevice(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appId := vars["appId"]

	var postData storage.AddDeviceRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&postData); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	if hnd.stg.AppExists(appId) {
		device := storage.Device{postData.Platform, postData.Token, int(time.Now().Unix())}
		hnd.stg.AddSubscriberDevice(appId, postData.SubscriberId, &device)
	} else {
		res.WriteError(w, 400, "App not found")
		return
	}

	res.WriteSuccess(w, 200, "")
}

func (hnd *Handlers) AddSubscriber(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appId := vars["appId"]
	channelId := vars["channelId"]

	var f storage.AddSubscriberRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	if hnd.stg.AppExists(appId) {
		hnd.stg.AddSubscriber(appId, channelId, f.SubscriberIds)
	} else {
		res.WriteError(w, 400, "App not found")
		return
	}
}

func (hnd *Handlers) AddChannel(w http.ResponseWriter, r *http.Request) {

	res := new(WrappedResponse)
	vars := mux.Vars(r)
	appId := vars["appId"]

	var f interface{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		res.WriteError(w, 400, err.Error())
		return
	}

	m := f.(map[string]interface{})

	if hnd.stg.AppExists(appId) {
		hnd.stg.AddChannel(appId, m["id"].(string))
	} else {
		res.WriteError(w, 400, "App not found")
		return
	}

}

func (hnd *Handlers) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appId := vars["appId"]
	channelId := vars["channelId"]

	hnd.stg.DeleteChannel(appId, channelId)
}
