package service

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.fixb.com/mir/push/storage"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// NewRouter creates and returns the router.
func NewRouter(stg *storage.Storage) *mux.Router {

	hnd := new(Handlers)
	hnd.stg = stg

	router := mux.NewRouter().StrictSlash(true)

	healthHandler := http.HandlerFunc(healthCheck)
	router.
		Methods("GET").
		Path("/health").
		Name("Health Check").
		HandlerFunc(healthHandler)

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

	return router
}
