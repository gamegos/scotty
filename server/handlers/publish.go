package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gamegos/gcmlib"
	"github.com/gamegos/jsend"
	"github.com/gamegos/scotty/server/context"
	"github.com/gorilla/mux"
)

// publishRequest represents http body of "publish" requests.
type publishRequest struct {
	Subscribers []string `json:"subscribers"`
	Channels    []string `json:"channels"`
	Message     *gcmlib.Message
}

func PublishMessage(jw jsend.JResponseWriter, r *http.Request, ctx *context.Context) {
	vars := mux.Vars(r)

	app, err := ctx.Storage.GetApp(vars["appId"])
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
		subscriberDevices, err := ctx.Storage.GetSubscriberDevices(app.ID, subscriberID)
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
