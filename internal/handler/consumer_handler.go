package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"main/internal/model"
	"main/internal/usecase"
)

type ConsumerHandler struct {
	consumerUsecase usecase.ConsumerUsecase
	limitUsecase    usecase.ConsumerLimitUsecase
}

func NewConsumerHandler(
	consumerUsecase usecase.ConsumerUsecase,
	limitUsecase usecase.ConsumerLimitUsecase,
) *ConsumerHandler {
	return &ConsumerHandler{
		consumerUsecase: consumerUsecase,
		limitUsecase:    limitUsecase,
	}
}

// RegisterConsumer handles POST /api/consumers
func (h *ConsumerHandler) RegisterConsumer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var consumer model.Consumer
	if err := json.NewDecoder(r.Body).Decode(&consumer); err != nil {
		log.Println("Error decoding request:", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}

	if err := h.consumerUsecase.RegisterConsumer(&consumer); err != nil {
		log.Println("Error registering consumer:", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Consumer registered successfully",
		"data":    consumer,
	})
}

// GetConsumer handles GET /api/consumers/{id}
func (h *ConsumerHandler) GetConsumer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 32)
	if err != nil || id == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid consumer ID"})
		return
	}

	consumer, err := h.consumerUsecase.GetConsumer(uint(id))
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Consumer not found"})
		return
	}

	respondJSON(w, http.StatusOK, consumer)
}

// AssignLimit handles POST /api/consumers/limits
func (h *ConsumerHandler) AssignLimit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var limit model.ConsumerLimit
	if err := json.NewDecoder(r.Body).Decode(&limit); err != nil {
		log.Println("Error decoding request:", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}

	if err := h.limitUsecase.AssignLimit(&limit); err != nil {
		log.Println("Error assigning limit:", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Limit assigned successfully",
		"data":    limit,
	})
}

// GetConsumerLimits handles GET /api/consumers/{id}/limits
func (h *ConsumerHandler) GetConsumerLimits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 32)
	if err != nil || id == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid consumer ID"})
		return
	}

	limits, err := h.limitUsecase.GetConsumerLimits(uint(id))
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Limits not found"})
		return
	}

	respondJSON(w, http.StatusOK, limits)
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
