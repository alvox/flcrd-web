package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

func extractSearchTerms(r *http.Request) []string {
	query := r.URL.Query()
	if len(query) == 0 {
		return nil
	}
	q := query.Get("q")
	if len(q) == 0 {
		return nil
	}
	return strings.Split(q, ",")
}
