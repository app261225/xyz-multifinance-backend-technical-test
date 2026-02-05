package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"main/config"
	"main/internal/handler"
	"main/internal/middleware"
	"main/internal/repository"
	"main/internal/usecase"
)

func main() {
	// 1. Database Connection
	db := config.ConnectDB()
	log.Println("✓ Database connected successfully")

	// 2. Repository Layer
	consumerRepo := repository.NewConsumerRepository(db)
	consumerLimitRepo := repository.NewConsumerLimitRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// 3. Usecase Layer
	consumerUC := usecase.NewConsumerUsecase(consumerRepo)
	limitUC := usecase.NewConsumerLimitUsecase(consumerLimitRepo)
	transactionUC := usecase.NewTransactionUsecase(transactionRepo, consumerLimitRepo)

	// 4. Handler Layer
	consumerHandler := handler.NewConsumerHandler(consumerUC, limitUC)
	transactionHandler := handler.NewTransactionHandler(transactionUC)

	// 5. Setup Routes with Security Middleware
	mux := http.NewServeMux()

	// Consumer endpoints
	mux.HandleFunc("/api/consumers", consumerHandler.RegisterConsumer)
	mux.HandleFunc("/api/consumers/get", consumerHandler.GetConsumer)
	mux.HandleFunc("/api/consumers/limits", consumerHandler.AssignLimit)
	mux.HandleFunc("/api/consumers/limits/get", consumerHandler.GetConsumerLimits)

	// Transaction endpoints
	mux.HandleFunc("/api/transactions", transactionHandler.CreateTransaction)
	mux.HandleFunc("/api/transactions/get", transactionHandler.GetTransaction)
	mux.HandleFunc("/api/transactions/consumer", transactionHandler.GetConsumerTransactions)
	mux.HandleFunc("/api/transactions/status", transactionHandler.UpdateTransactionStatus)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","version":"1.0.0"}`)
	})

	// Wrap mux with security middleware
	chain := middleware.SecurityHeaders(
		middleware.InputValidation(
			middleware.CORS(mux),
		),
	)

	// 6. Start Server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("✓ Starting API server on port %s\n", port)
	log.Printf("✓ Health check: http://localhost:%s/health\n", port)
	log.Printf("✓ OWASP Security Headers: ENABLED\n")
	log.Printf("✓ Input Validation: ENABLED\n")
	log.Printf("✓ CORS Protection: ENABLED\n")

	if err := http.ListenAndServe(":"+port, chain); err != nil {
		log.Fatal("Server error:", err)
	}
}
