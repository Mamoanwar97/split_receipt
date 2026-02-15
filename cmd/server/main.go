package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mamoanwar97/split-receipt/internal/database"
	"github.com/mamoanwar97/split-receipt/internal/handler"
	"github.com/mamoanwar97/split-receipt/internal/service"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/split_receipt?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	queries := database.New(pool)

	friendH := handler.NewFriendHandler(queries)
	receiptH := handler.NewReceiptHandler(queries)
	mealH := handler.NewMealHandler(queries)
	feeH := handler.NewFeeHandler(queries)
	settlementSvc := service.NewSettlementService(queries)
	settlementH := handler.NewSettlementHandler(settlementSvc)

	mux := http.NewServeMux()

	// Friends
	mux.HandleFunc("POST /friends", friendH.Create)
	mux.HandleFunc("GET /friends", friendH.List)
	mux.HandleFunc("GET /friends/{id}", friendH.Get)
	mux.HandleFunc("DELETE /friends/{id}", friendH.Delete)

	// Receipts
	mux.HandleFunc("POST /receipts", receiptH.Create)
	mux.HandleFunc("GET /receipts", receiptH.List)
	mux.HandleFunc("GET /receipts/{id}", receiptH.Get)
	mux.HandleFunc("DELETE /receipts/{id}", receiptH.Delete)

	// Meals
	mux.HandleFunc("POST /receipts/{id}/meals", mealH.Create)
	mux.HandleFunc("GET /receipts/{id}/meals", mealH.List)
	mux.HandleFunc("DELETE /receipts/{id}/meals/{mealId}", mealH.Delete)

	// Meal friends
	mux.HandleFunc("POST /receipts/{id}/meals/{mealId}/friends", mealH.AddFriend)
	mux.HandleFunc("DELETE /receipts/{id}/meals/{mealId}/friends/{friendId}", mealH.RemoveFriend)

	// Fixed fees
	mux.HandleFunc("POST /receipts/{id}/fixed-fees", feeH.CreateFixed)
	mux.HandleFunc("GET /receipts/{id}/fixed-fees", feeH.ListFixed)
	mux.HandleFunc("DELETE /receipts/{id}/fixed-fees/{feeId}", feeH.DeleteFixed)

	// Percentage fees
	mux.HandleFunc("POST /receipts/{id}/percentage-fees", feeH.CreatePercentage)
	mux.HandleFunc("GET /receipts/{id}/percentage-fees", feeH.ListPercentage)
	mux.HandleFunc("DELETE /receipts/{id}/percentage-fees/{feeId}", feeH.DeletePercentage)

	// Settlement
	mux.HandleFunc("GET /receipts/{id}/friends/{friendId}/settlement", settlementH.Get)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}
