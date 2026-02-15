package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mamoanwar97/split-receipt/internal/database"
)

type FeeHandler struct {
	queries *database.Queries
}

func NewFeeHandler(q *database.Queries) *FeeHandler {
	return &FeeHandler{queries: q}
}

// Fixed fees

func (h *FeeHandler) CreateFixed(w http.ResponseWriter, r *http.Request) {
	receiptID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid receipt id")
		return
	}

	var req struct {
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	fee, err := h.queries.CreateFixedFee(r.Context(), database.CreateFixedFeeParams{
		ReceiptID: receiptID,
		Name:      req.Name,
		Amount:    numericFromFloat(req.Amount),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create fixed fee")
		return
	}
	writeJSON(w, http.StatusCreated, fee)
}

func (h *FeeHandler) ListFixed(w http.ResponseWriter, r *http.Request) {
	receiptID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid receipt id")
		return
	}

	fees, err := h.queries.ListFixedFeesByReceipt(r.Context(), receiptID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list fixed fees")
		return
	}
	writeJSON(w, http.StatusOK, fees)
}

func (h *FeeHandler) DeleteFixed(w http.ResponseWriter, r *http.Request) {
	feeID, err := parseUUID(r.PathValue("feeId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid fee id")
		return
	}

	if err := h.queries.DeleteFixedFee(r.Context(), feeID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete fixed fee")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Percentage fees

func (h *FeeHandler) CreatePercentage(w http.ResponseWriter, r *http.Request) {
	receiptID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid receipt id")
		return
	}

	var req struct {
		Name       string   `json:"name"`
		Percentage float64  `json:"percentage"`
		CapAmount  *float64 `json:"cap_amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	params := database.CreatePercentageFeeParams{
		ReceiptID:  receiptID,
		Name:       req.Name,
		Percentage: numericFromFloat(req.Percentage),
	}
	if req.CapAmount != nil {
		params.CapAmount = numericFromFloat(*req.CapAmount)
	}

	fee, err := h.queries.CreatePercentageFee(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create percentage fee")
		return
	}
	writeJSON(w, http.StatusCreated, fee)
}

func (h *FeeHandler) ListPercentage(w http.ResponseWriter, r *http.Request) {
	receiptID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid receipt id")
		return
	}

	fees, err := h.queries.ListPercentageFeesByReceipt(r.Context(), receiptID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list percentage fees")
		return
	}
	writeJSON(w, http.StatusOK, fees)
}

func (h *FeeHandler) DeletePercentage(w http.ResponseWriter, r *http.Request) {
	feeID, err := parseUUID(r.PathValue("feeId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid fee id")
		return
	}

	if err := h.queries.DeletePercentageFee(r.Context(), feeID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete percentage fee")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
