package http_api

import (
	"context"
	"log"
	"net/http"
)

func (a *Api) middleware(h http.Handler) http.Handler {
	h = a.mwCoreCtx(h)
	h = a.mwRecovery(h)

	return h
}

func (a *Api) mwRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancelCtx, cancel := context.WithCancel(r.Context())
		r = r.WithContext(cancelCtx)
		defer func() {
			if err := recover(); err != nil {
				cancel()
				w.WriteHeader(http.StatusInternalServerError)
				log.Panicln(
					"Panic in http handler:",
					err,
					"method:", r.Method,
					"path:", r.URL,
				)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func (a *Api) mwCoreCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = uReqSetCore(r, a.cr)
		h.ServeHTTP(w, r)
	})
}
