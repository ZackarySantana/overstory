package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zackarysantana/overstory/src/clientmux"
	"github.com/zackarysantana/overstory/src/service"
)

type API struct {
	service *service.Service
}

func New(ctx context.Context, service *service.Service) *clientmux.ClientMux {
	api := &API{
		service: service,
	}

	mux := clientmux.New()
	// mux.HandleFunc("/api/v1/health", api.healthCheck)

	type CreateUserReq struct{ Name string }
	type CreateUserResp struct{ ID string }

	clientmux.HandleJSON(
		mux,
		http.MethodPost,
		"/users",
		func(ctx context.Context, in CreateUserReq) (CreateUserResp, error) {
			api.healthCheck(nil, nil)
			return CreateUserResp{ID: "42"}, nil
		},
	)

	return mux
}

func (api *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("Health check passed", r)
}
