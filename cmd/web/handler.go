package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func extractPaging(r *http.Request) (int, int, int, error) {
	query := r.URL.Query()
	if len(query) == 0 {
		return 1, 0, 5, nil // default
	}
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		return 0, 0, 0, err
	}
	limit, err := strconv.Atoi(query.Get("per_page"))
	if err != nil {
		return 0, 0, 0, err
	}
	offset := page
	if offset > 0 {
		offset = offset - 1
	}
	offset = offset * limit
	return page, offset, limit, nil
}

func addLinkHeader(w http.ResponseWriter, page, limit, total int) {
	s := `<https://flashcards.rocks/v0/public/decks?page=%d&per_page=%d>; rel="%s"`
	var res []string
	if page != 1 {
		res = append(res, fmt.Sprintf(s, page-1, limit, "prev"))
		res = append(res, fmt.Sprintf(s, 1, limit, "first"))
	}
	if total > (page * limit) {
		res = append(res, fmt.Sprintf(s, page+1, limit, "next"))
		if limit == 1 {
			res = append(res, fmt.Sprintf(s, total/limit, limit, "last"))
		} else {
			res = append(res, fmt.Sprintf(s, (total/limit)+1, limit, "last"))
		}
	}
	if len(res) > 0 {
		h := strings.Join(res, ",")
		w.Header().Set("Link", h)
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
}

func modelError(app *application, err error, w http.ResponseWriter, model string) bool {
	if err == models.ErrNoRecord {
		app.notFound(w, model)
		return true
	}
	if err != nil {
		app.serverError(w, err)
		return true
	}
	return false
}
