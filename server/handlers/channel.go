package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gamegos/jsend"
	"github.com/gamegos/scotty/server/context"
	"github.com/gorilla/mux"
)

// addSubscriberRequest holds the structure of new subscriber request.
type addSubscriberRequest struct {
	SubscriberIds []string `json:"subscribers"`
}

func AddSubscriber(jw jsend.JResponseWriter, r *http.Request, ctx *context.Context) {
	vars := mux.Vars(r)
	appID := vars["appId"]
	channelID := vars["channelId"]

	var f addSubscriberRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
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

	err := ctx.Storage.AddSubscriber(appID, channelID, f.SubscriberIds)

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(201).Send()
}

func AddChannel(jw jsend.JResponseWriter, r *http.Request, ctx *context.Context) {
	vars := mux.Vars(r)
	appID := vars["appId"]

	var f interface{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		jw.Status(400).Message(err.Error()).Send()
		return
	}

	m := f.(map[string]interface{})

	if app, err := ctx.Storage.GetApp(appID); app == nil {
		if err != nil {
			jw.Status(500).Message(err.Error())
		} else {
			jw.Status(404).Message("App not found.")
		}
		return
	}

	err := ctx.Storage.AddChannel(appID, m["id"].(string))

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(201).Send()
}

func DeleteChannel(jw jsend.JResponseWriter, r *http.Request, ctx *context.Context) {
	vars := mux.Vars(r)
	appID := vars["appId"]
	channelID := vars["channelId"]

	err := ctx.Storage.DeleteChannel(appID, channelID)

	if err != nil {
		jw.Status(500).Message(err.Error()).Send()
		return
	}

	jw.Status(200).Send()
}
