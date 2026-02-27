package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GameXost/wbTestCase/internal/apperror"
	"github.com/GameXost/wbTestCase/internal/models"
	"github.com/GameXost/wbTestCase/metrics"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderUID string) (*models.Order, error)
}

type Handler struct {
	Service OrderService
}

func NewHandler(srv OrderService) *Handler {
	return &Handler{Service: srv}
}
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {

	metrics.RequestsTotal.Inc()

	orderUID := chi.URLParam(r, "order_uid")
	if orderUID == "" {
		handleHTTPErr(w, apperror.ErrOrderUIDMissing)
		return
	}
	order, err := h.Service.GetOrder(r.Context(), orderUID)
	if err != nil {
		handleHTTPErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
	metrics.RequestsSuccess.Inc()

}

func handleHTTPErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		metrics.RequestsNotFound.Inc()
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, apperror.ErrOrderUIDMissing):
		metrics.RequestsBadRequest.Inc()
		http.Error(w, "empty order_id", http.StatusBadRequest)
	default:
		log.Printf("internal server error getting order: %v", err)
		metrics.RequestsServerError.Inc()
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}
