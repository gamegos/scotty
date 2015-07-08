package service

import (
	"encoding/json"
	"net/http"
)

const statusSuccess = "success"
const statusError = "error"

// WrappedResponse holds a response data.
type WrappedResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// WriteSuccess sends response in case of success.
func (res *WrappedResponse) WriteSuccess(w http.ResponseWriter, code int, data interface{}) {
	res.Write(w, code, data, "")
}

// WriteError sends response in case of error.
func (res *WrappedResponse) WriteError(w http.ResponseWriter, code int, msg string) {
	res.Write(w, code, nil, msg)
}

// Write writes the response to the ResponseWriter.
func (res *WrappedResponse) Write(w http.ResponseWriter, code int, data interface{}, msg string) {
	res.Code = code

	if len(msg) > 0 {
		res.Message = msg
	}

	if code >= 200 && code < 300 {
		res.Status = statusSuccess
		res.Data = data
	} else {
		res.Status = statusError
	}

	bodyJSON, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bodyJSON)
}
