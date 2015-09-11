package server

import (
	"net/http"

	"github.com/gamegos/jsend"
	"github.com/gamegos/scotty/server/context"
	"github.com/gamegos/scotty/server/handlers"
	"github.com/gamegos/scotty/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	ctx    *context.Context
	//addr   string
}

// Init initializes a scotty http server.
func Init(stg storage.Storage) *Server {
	s := &Server{}
	//s.addr = addr
	s.ctx = &context.Context{stg}
	s.router = initRouter(s.ctx)

	return s
}

// Run starts a scotty http server.
func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

type handlerFunc func(w jsend.JResponseWriter, r *http.Request, ctx *context.Context)

type mainHandler struct {
	ctx *context.Context
	f   handlerFunc
}

func (h *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jw := jsend.Wrap(w)
	h.f(jw, r, h.ctx)
	jw.Send()
}

// initRouter creates and returns the router.
func initRouter(ctx *context.Context) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	wrap := func(f handlerFunc) *mainHandler {
		return &mainHandler{ctx, f}
	}

	router.
		Methods("GET").
		Path("/health").
		Name("Health Check").
		Handler(wrap(handlers.GetHealth))

	router.
		Methods("GET").
		Path("/apps/{appId}").
		Name("Get App").
		Handler(wrap(handlers.GetApp))

	router.
		Methods("POST").
		Path("/apps").
		Name("Create App").
		Handler(wrap(handlers.CreateApp))

	router.
		Methods("PUT").
		Path("/apps/{appId}").
		Name("Update App").
		Handler(wrap(handlers.UpdateApp))

	router.
		Methods("POST").
		Path("/apps/{appId}/devices").
		Name("Add Device to Subscriber").
		Handler(wrap(handlers.AddDevice))

	router.
		Methods("POST").
		Path("/apps/{appId}/channels").
		Name("Add Channel to App").
		Handler(wrap(handlers.AddChannel))

	router.
		Methods("DELETE").
		Path("/apps/{appId}/channels/{channelId}").
		Name("Delete Channel from App").
		Handler(wrap(handlers.DeleteChannel))

	router.
		Methods("POST").
		Path("/apps/{appId}/channels/{channelId}/subscribers").
		Name("Add Subscriber to Channel").
		Handler(wrap(handlers.AddSubscriber))

	router.
		Methods("POST").
		Path("/apps/{appId}/publish").
		Name("Publish a message").
		Handler(wrap(handlers.PublishMessage))

	router.NotFoundHandler = wrap(handlers.NotfoundHandler)

	return router
}
