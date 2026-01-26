package httpserver

import (
	"net/http"
	"sync/atomic"
)

type Health struct {
	ready atomic.Bool
}

func (h *Health) SetReady(v bool) {
	h.ready.Store(v)
}

func (h *Health) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Health) Readyz(w http.ResponseWriter, _ *http.Request) {
	if h.ready.Load() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	_, _ = w.Write([]byte("not ready"))
}
