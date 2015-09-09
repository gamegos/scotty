package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gamegos/gcmlib"
	"github.com/gamegos/jsend"
	"github.com/gamegos/scotty/storage"
	"github.com/gorilla/mux"
)

// addDeviceRequest holds the structure of new device request.
type addDeviceRequest struct {
	SubscriberID string `json:"subscriberId"`
	Platform     string `json:"platform"`
	Token        string `json:"token"`
}

// addSubscriberRequest holds the structure of new subscriber request.
type addSubscriberRequest struct {
	SubscriberIds []string `json:"subscribers"`
}

// publishRequest represents http body of "publish" requests.
type publishRequest struct {
	Subscribers []string `json:"subscribers"`
	Channels    []string `json:"channels"`
	Message     *gcmlib.Message
}

// Handlers holds the handler functions to be run with different routes.
type Handlers struct {
	stg storage.Storage
}

func (hnd *Handlers) getHealth(w http.ResponseWriter, r *http.Request) {
	jsend.Wrap(w).Status(200).Send()
}

func (hnd *Handlers) createApp(w http.ResponseWriter, r *http.Request) {

	jw := jsend.Wrap(w)
	var app storage.App

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		jw.Status(400).Message(err.Error()).Send()
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
		jw.Status(400).Message(fmt.Sprintf("App %v does not exist.", appID)).Send()
		return
	}

	var app storage.App
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		jw.Status(400).Message(err.Error()).Send()
		return
	}

	if appID != app.ID {
		jw.Status(400).Message("AppID mismatch").Send()
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
		jw.Status(404).Message("App not found.").Send()
		return
	}

	jw.Status(200).Data(appData).Send()
}

func (hnd *Handlers) addDevice(w http.ResponseWriter, r *http.Request) {

	jw := jsend.Wrap(w)
	vars := mux.Vars(r)
	appID := vars["appId"]

	var postData addDeviceRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&postData); err != nil {
		jw.Status(400).Message(err.Error()).Send()
		return
	}

	if !hnd.stg.AppExists(appID) {
		jw.Status(400).Message("App not found.").Send()
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

	var f addSubscriberRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		jw.Status(400).Message(err.Error()).Send()
		return
	}

	if !hnd.stg.AppExists(appID) {
		jw.Status(400).Message("App not found.").Send()
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
		jw.Status(400).Message(err.Error()).Send()
		return
	}

	m := f.(map[string]interface{})

	if !hnd.stg.AppExists(appID) {
		jw.Status(400).Message("App not found.").Send()
		return
	}

	err := hnd.stg.AddChannel(appID, m["id"].(string))

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
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

func (hnd *Handlers) publishMessage(w http.ResponseWriter, r *http.Request) {
	jw := jsend.Wrap(w)
	vars := mux.Vars(r)

	app, err := hnd.stg.GetApp(vars["appId"])
	if app == nil {
		log.Println("App not found:", app, err)
		jw.Status(404).Message("App not found").Send()
		return
	}

	publishReq := new(publishRequest)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&publishReq); err != nil {
		log.Println("Could not decode message, ", err)
		jw.Status(400).Message("Could not decode message: " + err.Error()).Send()
		return
	}

	deviceTokens := make([]string, 0, len(publishReq.Subscribers))

	for _, subscriberID := range publishReq.Subscribers {
		log.Println(subscriberID)
		subscriberDevices, err := hnd.stg.GetSubscriberDevices(app.ID, subscriberID)
		if err != nil {
			log.Println("Error, ", err)
		}

		for _, device := range subscriberDevices {
			deviceTokens = append(deviceTokens, device.Token)
		}

		log.Printf("Devices %#v\n", deviceTokens)
	}

	client := gcmlib.NewClient(gcmlib.Config{
		APIKey: app.GCM.APIKey,
	})

	msg := publishReq.Message
	msg.RegistrationIDs = deviceTokens

	if err := msg.Validate(); err != nil {
		jw.Status(400).Message(err.Error()).Send()
		return
	}

	result, gcmErr := client.Send(msg)
	log.Printf("GCM Request: %#v, %#v\n", result, err)

	if gcmErr != nil {
		jw.Status(400).Message(gcmErr.Error()).Send()
		return
	}

	jw.Data(result).Send()
}
