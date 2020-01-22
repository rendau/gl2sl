package http_api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (a *Api) createRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("", hMsg)
	r.HandleFunc("/", hMsg)

	return a.middleware(r)
}
