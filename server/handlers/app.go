package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gamegos/jsend"
	"github.com/gamegos/scotty/server/context"
	"github.com/gamegos/scotty/storage"
	"github.com/gorilla/mux"
)

// addDeviceRequest holds the structure of new device request.
type addDeviceRequest struct {
	SubscriberID string `json:"subscriberId"`
	Platform     string `json:"platform"`
	Token        string `json:"token"`
}

func CreateApp(jw jsend.JResponseWriter, r *http.Request, ctx *context.Context) {
	var app *storage.App

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		jw.Status(400).Message(err.Error()).Send()
		return
	}

	err := ctx.Storage.PutApp(app)

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(201).Send()
}

func UpdateApp(jw jsend.JResponseWriter, r *http.Request, ctx *context.Context) {
	vars := mux.Vars(r)
	appID := vars["appId"]

	if !ctx.Storage.AppExists(appID) {
		jw.Status(400).Message(fmt.Sprintf("App %v does not exist.", appID)).Send()
		return
	}

	var app *storage.App
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&app); err != nil {
		jw.Status(400).Message(err.Error()).Send()
		return
	}

	if appID != app.ID {
		jw.Status(400).Message("AppID mismatch").Send()
		return
	}

	err := ctx.Storage.PutApp(app)

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(200).Send()
}

func GetApp(jw jsend.JResponseWriter, r *http.Request, ctx *context.Context) {
	vars := mux.Vars(r)
	appID := vars["appId"]

	appData, err := ctx.Storage.GetApp(appID)

	if err != nil {
		jw.Status(404).Message("App not found.").Send()
		return
	}

	jw.Status(200).Data(appData).Send()
}

func AddDevice(jw jsend.JResponseWriter, r *http.Request, ctx *context.Context) {
	vars := mux.Vars(r)
	appID := vars["appId"]

	var postData addDeviceRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&postData); err != nil {
		jw.Status(400).Message(err.Error()).Send()
		return
	}

	if !ctx.Storage.AppExists(appID) {
		jw.Status(400).Message("App not found.").Send()
		return
	}

	device := storage.Device{
		Platform:  postData.Platform,
		Token:     postData.Token,
		CreatedAt: int(time.Now().Unix()),
	}
	err := ctx.Storage.AddSubscriberDevice(appID, postData.SubscriberID, &device)

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(201).Send()
}
