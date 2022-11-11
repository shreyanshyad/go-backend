package utils

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Success  bool        `json:"success"`
	RespData interface{} `json:"data,omitempty"`
	Error    string      `json:"error,omitempty"`
	Message  string      `json:"message,omitempty"`
}

func WriteSuccessResponse(w http.ResponseWriter, code int, data interface{}) error {
	resp := response{
		Success:  true,
		RespData: data,
	}
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	return e.Encode(resp)
}

func WriteSuccessResponseMsg(w http.ResponseWriter, code int, msg string) error {
	resp := response{
		Success: true,
		Message: msg,
	}
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	return e.Encode(resp)
}

func WriteFailureResponse(w http.ResponseWriter, code int, err string) error {
	resp := response{
		Success: false,
		Error:   err,
	}
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	return e.Encode(resp)
}
