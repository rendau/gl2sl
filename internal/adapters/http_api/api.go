package http_api

import (
	"context"
	"github.com/rendau/gl2sl/internal/domain/core"
	"net/http"
	"time"
)

type Api struct {
	listen string
	server *http.Server
	cr     *core.St
}

func NewApi(listen string, cr *core.St) *Api {
	return &Api{
		listen: listen,
		cr:     cr,
	}
}

func (a *Api) Start() error {
	a.server = &http.Server{
		Addr:         a.listen,
		Handler:      a.createRouter(),
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Second,
	}
	return a.server.ListenAndServe()
}

func (a *Api) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
