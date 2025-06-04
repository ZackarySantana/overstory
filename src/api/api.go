package api

import (
	"context"
	"net/http"

	"github.com/zackarysantana/overstory/src/service"
)

type API struct {
	service *service.Service
}

type CreateUserReq struct{ Name string }
type CreateUserResp struct{ ID string }

func New(ctx context.Context, service *service.Service) http.Handler {
	api := &API{
		service: service,
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/health", healthCheck(api))

	return mux
}
