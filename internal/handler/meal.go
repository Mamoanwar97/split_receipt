package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mamoanwar97/split-receipt/internal/database"
)

type MealHandler struct {
	queries *database.Queries
}

func NewMealHandler(q *database.Queries) *MealHandler {
	return &MealHandler{queries: q}
}

func (h *MealHandler) Create(w http.ResponseWriter, r *http.Request) {
	receiptID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid receipt id")
		return
	}

	var req struct {
		Name         string  `json:"name"`
		Quantity     int32   `json:"quantity"`
		TotalPrice   float64 `json:"total_price"`
		PricePerUnit float64 `json:"price_per_unit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	meal, err := h.queries.CreateMeal(r.Context(), database.CreateMealParams{
		ReceiptID:    receiptID,
		Name:         req.Name,
		Quantity:     req.Quantity,
		TotalPrice:   numericFromFloat(req.TotalPrice),
		PricePerUnit: numericFromFloat(req.PricePerUnit),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create meal")
		return
	}
	writeJSON(w, http.StatusCreated, meal)
}

func (h *MealHandler) List(w http.ResponseWriter, r *http.Request) {
	receiptID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid receipt id")
		return
	}

	meals, err := h.queries.ListMealsByReceipt(r.Context(), receiptID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list meals")
		return
	}
	writeJSON(w, http.StatusOK, meals)
}

func (h *MealHandler) Delete(w http.ResponseWriter, r *http.Request) {
	mealID, err := parseUUID(r.PathValue("mealId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid meal id")
		return
	}

	if err := h.queries.DeleteMeal(r.Context(), mealID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete meal")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MealHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	mealID, err := parseUUID(r.PathValue("mealId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid meal id")
		return
	}

	var req struct {
		FriendID string `json:"friend_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	friendID, err := parseUUID(req.FriendID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid friend_id")
		return
	}

	if err := h.queries.AddFriendToMeal(r.Context(), database.AddFriendToMealParams{
		MealID:   mealID,
		FriendID: friendID,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add friend to meal")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MealHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	mealID, err := parseUUID(r.PathValue("mealId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid meal id")
		return
	}

	friendID, err := parseUUID(r.PathValue("friendId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid friend id")
		return
	}

	if err := h.queries.RemoveFriendFromMeal(r.Context(), database.RemoveFriendFromMealParams{
		MealID:   mealID,
		FriendID: friendID,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove friend from meal")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
