package handlers

import (
	"net/http"

	"github.com/gamegos/jsend"
	"github.com/gamegos/scotty/server/context"
)

func GetHealth(jw jsend.JResponseWriter, r *http.Request, ctx *context.Context) {
	jw.Status(200).Data("ok")
}
