package handlers

import (
	"encoding/json"
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

	if app, err := ctx.Storage.GetApp(appID); app == nil {
		if err != nil {
			jw.Status(500).Message(err.Error())
		} else {
			jw.Status(404).Message("App not found.")
		}
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

	app, err := ctx.Storage.GetApp(appID)

	if err != nil {
		jw.Status(500).Message(err.Error())
		return
	}

	if app == nil {
		jw.Status(404).Message("App not found.")
		return
	}

	jw.Data(app)
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

	if app, err := ctx.Storage.GetApp(appID); app == nil {
		if err != nil {
			jw.Status(500).Message(err.Error())
		} else {
			jw.Status(404).Message("App not found.")
		}
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
