package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int

	Error error
}

func (r *Response) status(fallback int) int {
	if r != nil && r.Status != 0 {
		return r.Status
	}
	return fallback
}

func (r *Response) error() error {
	if r != nil && r.Error != nil {
		return r.Error
	}
	return nil
}

func (r *Response) handle(w http.ResponseWriter, returnVal any, requireReturn bool) {
	if err := r.error(); err != nil {
		http.Error(w, err.Error(), r.status(http.StatusInternalServerError))
		return
	}

	if returnVal == nil {
		if !requireReturn {
			w.WriteHeader(r.status(http.StatusOK))
			return
		}
		http.Error(w, "no result returned", r.status(http.StatusInternalServerError))
		return
	}

	resultJSON, err := json.Marshal(returnVal)
	if err != nil {
		http.Error(w, err.Error(), r.status(http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.status(http.StatusOK))
	w.Write(resultJSON)
}

type JSONReturnHandler[Accept any, Return any] func(*http.Request, *Accept) (*Return, *Response)

func (j JSONReturnHandler[Accept, Return]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var accept Accept
	if err := json.NewDecoder(r.Body).Decode(&accept); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	returnVal, response := j(r, &accept)
	response.handle(w, returnVal, true)
}

type JSONHandler[Accept any] func(*http.Request, *Accept) *Response

func (j JSONHandler[Accept]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var accept Accept
	if err := json.NewDecoder(r.Body).Decode(&accept); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j(r, &accept).handle(w, nil, false)
}

type GetReturnHandler[Return any] func(*http.Request) (*Return, *Response)

func (g GetReturnHandler[Return]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	returnVal, response := g(r)
	response.handle(w, returnVal, true)
}

type GetHandler func(*http.Request) *Response

func (g GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g(r).handle(w, nil, false)
}
