package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mamoanwar97/split-receipt/internal/database"
)

type ReceiptHandler struct {
	queries *database.Queries
}

func NewReceiptHandler(q *database.Queries) *ReceiptHandler {
	return &ReceiptHandler{queries: q}
}

func (h *ReceiptHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RestaurantName string  `json:"restaurant_name"`
		DateTime       string  `json:"date_time"`
		Total          float64 `json:"total"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RestaurantName == "" {
		writeError(w, http.StatusBadRequest, "restaurant_name is required")
		return
	}

	dt, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		writeError(w, http.StatusBadRequest, "date_time must be RFC3339 format")
		return
	}

	receipt, err := h.queries.CreateReceipt(r.Context(), database.CreateReceiptParams{
		RestaurantName: req.RestaurantName,
		DateTime:       pgtype.Timestamptz{Time: dt, Valid: true},
		Total:          numericFromFloat(req.Total),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create receipt")
		return
	}
	writeJSON(w, http.StatusCreated, receipt)
}

func (h *ReceiptHandler) List(w http.ResponseWriter, r *http.Request) {
	receipts, err := h.queries.ListReceipts(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list receipts")
		return
	}
	writeJSON(w, http.StatusOK, receipts)
}

func (h *ReceiptHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	receipt, err := h.queries.GetReceipt(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "receipt not found")
		return
	}
	writeJSON(w, http.StatusOK, receipt)
}

func (h *ReceiptHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.queries.DeleteReceipt(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete receipt")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func numericFromFloat(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%.2f", f))
	return n
}
