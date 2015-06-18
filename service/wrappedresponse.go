package service

import (
	"encoding/json"
	"net/http"
)

const STATUS_SUCCESS = "success"
const STATUS_ERROR = "error"

type WrappedResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (res *WrappedResponse) WriteSuccess(w http.ResponseWriter, code int, data interface{}) {
	res.Write(w, code, data, "")
}

func (res *WrappedResponse) WriteError(w http.ResponseWriter, code int, msg string) {
	res.Write(w, code, nil, msg)
}

func (res *WrappedResponse) Write(w http.ResponseWriter, code int, data interface{}, msg string) {
	res.Code = code

	if len(msg) > 0 {
		res.Message = msg
	}

	if code >= 200 && code < 300 {
		res.Status = STATUS_SUCCESS
		res.Data = data
	} else {
		res.Status = STATUS_ERROR
	}

	bodyJson, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bodyJson)
}
