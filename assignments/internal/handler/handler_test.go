package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xsenzel/emblabs-golang/internal/handler"
	"github.com/0xsenzel/emblabs-golang/internal/service"
)

func executeRequest(h *handler.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	http.HandlerFunc(h.Pay).ServeHTTP(rr, req)

	return rr
}

func TestPayHandler_NewTransaction(t *testing.T) {
	svc := service.NewPaymentService()
	h := handler.NewHandler(svc)

	txID := "tx123"
	reqBody := service.Transaction{
		UserID:        "user123",
		Amount:        100,
		TransactionID: txID,
	}
	body, _ := json.Marshal(reqBody)

	rr := executeRequest(h, "POST", "/pay", body)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestPayHandler_IdempotentTransaction(t *testing.T) {
	svc := service.NewPaymentService()
	h := handler.NewHandler(svc)

	txID := "Idempotent-tx123"
	reqBody := service.Transaction{
		UserID:        "user1234",
		Amount:        1000,
		TransactionID: txID,
	}
	body, _ := json.Marshal(reqBody)

	// First call: New transaction
	rr1 := executeRequest(h, http.MethodPost, "/pay", body)
	if rr1.Code != http.StatusCreated {
		t.Fatalf("Setup failed: first call failed.")
	}

	// Second call: Same transaction ID
	rr2 := executeRequest(h, http.MethodPost, "/pay", body)
	if rr2.Code != http.StatusOK {
		t.Fatalf("Idempotent call failed: expected %d (OK), got %d. Body: %s", http.StatusOK, rr2.Code, rr2.Body.String())
	}
}
