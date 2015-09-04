package service

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gamegos/scotty/storage"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// NewRouter creates and returns the router.
func NewRouter(stg storage.Storage) *mux.Router {

	hnd := new(Handlers)
	hnd.stg = stg

	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("GET").
		Path("/health").
		Name("Health Check").
		HandlerFunc(hnd.getHealth)

	router.
		Methods("GET").
		Path("/apps/{appId}").
		Name("Get App").
		HandlerFunc(hnd.getApp)

	router.
		Methods("POST").
		Path("/apps").
		Name("Create App").
		HandlerFunc(hnd.createApp)

	router.
		Methods("PUT").
		Path("/apps/{appId}").
		Name("Update App").
		HandlerFunc(hnd.updateApp)

	router.
		Methods("POST").
		Path("/apps/{appId}/devices").
		Name("Add Device to Subscriber").
		HandlerFunc(hnd.addDevice)

	router.
		Methods("POST").
		Path("/apps/{appId}/channels").
		Name("Add Channel to App").
		HandlerFunc(hnd.addChannel)

	router.
		Methods("DELETE").
		Path("/apps/{appId}/channels/{channelId}").
		Name("Delete Channel from App").
		HandlerFunc(hnd.deleteChannel)

	router.
		Methods("POST").
		Path("/apps/{appId}/channels/{channelId}/subscribers").
		Name("Add Subscriber to Channel").
		HandlerFunc(hnd.addSubscriber)

	router.
		Methods("POST").
		Path("/apps/{appId}/publish").
		Name("Publish a message").
		HandlerFunc(hnd.publishMessage)

	return router
}
