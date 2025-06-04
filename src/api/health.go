package api

import (
	"net/http"

	"github.com/zackarysantana/overstory/src/api/handlers"
)

func healthCheck(_ *API) handlers.GetHandler {
	return func(r *http.Request) *handlers.Response {
		return &handlers.Response{
			Status: http.StatusTeapot,
		}
	}
}
