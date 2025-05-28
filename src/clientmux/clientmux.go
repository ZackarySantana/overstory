package clientmux

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"sync"
)

type RouteMeta struct {
	Method   string
	Pattern  string
	ReqType  reflect.Type
	RespType reflect.Type
	Handler  http.Handler
}

type ClientMux struct {
	mux    *http.ServeMux
	routes []RouteMeta
	mu     sync.RWMutex
}

func New() *ClientMux { return &ClientMux{mux: http.NewServeMux()} }

func (cm *ClientMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cm.mux.ServeHTTP(w, r)
}

func (cm *ClientMux) Routes() []RouteMeta {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	out := make([]RouteMeta, len(cm.routes))
	copy(out, cm.routes)
	return out
}

// Generic handler type.
type Handler[TReq any, TResp any] func(ctx context.Context, in TReq) (TResp, error)

// HandleJSON registers a JSON endpoint with compile-time request/response types.
// It’s a **function**, not a method — so we can use type parameters.
func HandleJSON[TReq any, TResp any](
	cm *ClientMux,
	method, pattern string,
	h Handler[TReq, TResp],
) {
	// Wrapper adds (un)marshalling and error handling.
	wrapped := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req TReq
		if method != http.MethodGet && method != http.MethodHead {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		resp, err := h(r.Context(), req)
		if err != nil {
			status := http.StatusInternalServerError
			var he interface{ Status() int }
			if errors.As(err, &he) {
				status = he.Status()
			}
			http.Error(w, err.Error(), status)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodHead {
			return
		}
		_ = json.NewEncoder(w).Encode(resp)
	}

	// Register and record metadata.
	cm.mux.Handle(pattern, http.HandlerFunc(wrapped))

	var zeroReq TReq
	var zeroResp TResp
	meta := RouteMeta{
		Method:   method,
		Pattern:  pattern,
		ReqType:  reflect.TypeOf(zeroReq),
		RespType: reflect.TypeOf(zeroResp),
		Handler:  http.HandlerFunc(wrapped),
	}

	cm.mu.Lock()
	cm.routes = append(cm.routes, meta)
	cm.mu.Unlock()
}
