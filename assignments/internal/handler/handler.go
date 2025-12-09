package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0xsenzel/emblabs-golang/internal/service"
)

type Handler struct {
	Service *service.PaymentService
}

func NewHandler(service *service.PaymentService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Pay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req service.Transaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body format", http.StatusBadRequest)
		return
	}
	// ensure request body closed
	defer r.Body.Close()

	result, err := h.Service.ProcessPayment(req)
	if errors.Is(err, service.ErrAlreadyProcessed) {
		w.WriteHeader(http.StatusOK)
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	w.Header().Set("Content-Type", "application/json")
	if encodeErr := json.NewEncoder(w).Encode(result); encodeErr != nil {
		return
	}
}
