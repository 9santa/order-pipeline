package orderapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Handlers struct {
	svc *Service
}

func NewHandlers(svc *Service) *Handlers {
	return &Handlers{svc: svc}
}

type createOrderReq struct {
	CustomerID  string          `json:"customer_id"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Currency    string          `json:"currency"`
}

type createOrderResp struct {
	Order   Order  `json:"order"`
	EventID string `json:"event_id"`
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("/orders", h.handleOrders)
	mux.HandleFunc("/orders/", h.handleOrderByID)
}

func (h *Handlers) handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createOrder(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) createOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Decode JSON into createOrderReq struct
	var req createOrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.CustomerID) == "" || req.TotalAmount.LessThan(decimal.NewFromInt(0)) {
		http.Error(w, "customer_id and total_amount required", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		req.Currency = "RUB"
	}

	order, eventID, err := h.svc.CreateOrder(ctx, req.CustomerID, req.TotalAmount, req.Currency)
	if err != nil {
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	resp := createOrderResp{Order: order, EventID: eventID}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handlers) handleOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/orders/")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	order, ok := h.svc.GetOrder(r.Context(), id)
	if !ok {
		http.Error(w, "order with this id not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, order)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
