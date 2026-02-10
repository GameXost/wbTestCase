package server

import (
	"github.com/GameXost/wbTestCase/internal/errHandle"
	"github.com/GameXost/wbTestCase/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerGetOrderSuccess(t *testing.T) {
	serv := NewMockOrderService(t)
	handler := NewHandler(serv)

	order := &models.Order{OrderUId: "test1"}

	serv.EXPECT().GetOrder(mock.Anything, "test1").Return(order, nil)

	r := chi.NewRouter()
	r.Get("/orders/{order_uid}", handler.GetOrder)

	req := httptest.NewRequest(http.MethodGet, "/orders/test1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerGetOrderNotFound(t *testing.T) {
	serv := NewMockOrderService(t)
	handler := NewHandler(serv)

	serv.EXPECT().GetOrder(mock.Anything, "test2").Return(nil, errHandle.ErrNotFound)

	r := chi.NewRouter()
	r.Get("/orders/{order_uid}", handler.GetOrder)
	req := httptest.NewRequest(http.MethodGet, "/orders/test2", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestHandlerGetOrderServErr(t *testing.T) {
	serv := NewMockOrderService(t)
	handler := NewHandler(serv)

	serv.EXPECT().GetOrder(mock.Anything, "test3").Return(nil, errHandle.ErrServer)

	r := chi.NewRouter()
	r.Get("/orders/{order_uid}", handler.GetOrder)

	req := httptest.NewRequest(http.MethodGet, "/orders/test3", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
