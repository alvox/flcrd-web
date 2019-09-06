package main

import (
	"encoding/json"
	"net/http"
)

func reply(w http.ResponseWriter, status int, obj interface{}) {
	w.WriteHeader(status)
	if obj != nil {
		writeJsonResponse(w, obj)
	}
}

func writeJsonResponse(w http.ResponseWriter, obj interface{}) {
	out, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
