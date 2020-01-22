package http_api

import (
	"context"
	"github.com/rendau/gl2sl/internal/domain/core"
	"net/http"
)

const (
	CoreCtxKey = "core"
)

func uReqSetCore(r *http.Request, cr *core.St) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), CoreCtxKey, cr))
}

func uReqGetCore(r *http.Request) *core.St {
	return r.Context().Value(CoreCtxKey).(*core.St)
}
