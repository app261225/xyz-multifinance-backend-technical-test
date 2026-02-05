package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"main/internal/model"
	"main/internal/usecase"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
	}
}

// CreateTransaction handles POST /api/transactions - with concurrent transaction handling
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var transaction model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Println("Error decoding request:", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}

	// This method uses mutex to handle concurrent transactions safely
	if err := h.transactionUsecase.CreateTransaction(&transaction); err != nil {
		log.Println("Error creating transaction:", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Transaction created successfully",
		"data":    transaction,
	})
}

// GetTransaction handles GET /api/transactions/{id}
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 32)
	if err != nil || id == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.transactionUsecase.GetTransaction(uint(id))
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Transaction not found"})
		return
	}

	respondJSON(w, http.StatusOK, transaction)
}

// GetConsumerTransactions handles GET /api/consumers/{id}/transactions
func (h *TransactionHandler) GetConsumerTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 32)
	if err != nil || id == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid consumer ID"})
		return
	}

	transactions, err := h.transactionUsecase.GetConsumerTransactions(uint(id))
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Transactions not found"})
		return
	}

	respondJSON(w, http.StatusOK, transactions)
}

// UpdateTransactionStatus handles PUT /api/transactions/{id}/status
func (h *TransactionHandler) UpdateTransactionStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 32)
	if err != nil || id == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid transaction ID"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}

	if err := h.transactionUsecase.UpdateTransactionStatus(uint(id), req.Status); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Transaction status updated successfully"})
}
