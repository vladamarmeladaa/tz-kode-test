package response

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	resp := Response{
		Status: StatusOK,
	}
	return resp

}

func Error(msg string) Response {
	resp := Response{
		Status: StatusError,
		Error:  msg,
	}
	return resp
}

func ResponseError(w http.ResponseWriter, msg string) {
	err := json.NewEncoder(w).Encode(Error(msg))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
