package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GameXost/wbTestCase/internal/errHandle"
	"github.com/GameXost/wbTestCase/models"
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
	orderUID := chi.URLParam(r, "order_uid")
	if orderUID == "" {
		http.Error(w, "empty order_id", http.StatusBadRequest)
		return
	}
	order, err := h.Service.GetOrder(r.Context(), orderUID)
	if err != nil {
		switch {
		case errors.Is(err, errHandle.ErrNotFound):
			http.Error(w, "not found", http.StatusNotFound)
		default:
			log.Printf("internal server error getting order: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)

		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
