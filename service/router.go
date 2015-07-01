package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.fixb.com/mir/push/storage"
)

func WrapLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		finish := time.Since(start)

		log.Printf("%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			finish,
		)
	})
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func NewRouter(stg *storage.Storage) *mux.Router {

	hnd := new(Handlers)
	hnd.stg = stg

	router := mux.NewRouter().StrictSlash(true)

	healthHandler := http.HandlerFunc(HealthCheck)
	router.
		Methods("GET").
		Path("/health").
		Name("Health Check").
		HandlerFunc(healthHandler)

	router.
		Methods("GET").
		Path("/apps/{appId}").
		Name("Get App").
		HandlerFunc(hnd.GetApp)

	router.
		Methods("POST").
		Path("/apps").
		Name("Create App").
		HandlerFunc(hnd.CreateApp)

	router.
		Methods("PUT").
		Path("/apps/{appId}").
		Name("Update App").
		HandlerFunc(hnd.UpdateApp)

	router.
		Methods("POST").
		Path("/apps/{appId}/devices").
		Name("Add Device to Subscriber").
		HandlerFunc(hnd.AddDevice)

	router.
		Methods("POST").
		Path("/apps/{appId}/channels").
		Name("Add Channel to App").
		HandlerFunc(hnd.AddChannel)

	router.
		Methods("DELETE").
		Path("/apps/{appId}/channels/{channelId}").
		Name("Delete Channel from App").
		HandlerFunc(hnd.DeleteChannel)

	router.
		Methods("POST").
		Path("/apps/{appId}/channels/{channelId}/subscribers").
		Name("Add Subscriber to Channel").
		HandlerFunc(hnd.AddSubscriber)

	return router
}
