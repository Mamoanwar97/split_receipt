package handler

import (
	"net/http"

	"github.com/mamoanwar97/split-receipt/internal/service"
)

type SettlementHandler struct {
	service *service.SettlementService
}

func NewSettlementHandler(s *service.SettlementService) *SettlementHandler {
	return &SettlementHandler{service: s}
}

func (h *SettlementHandler) Get(w http.ResponseWriter, r *http.Request) {
	receiptID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid receipt id")
		return
	}

	friendID, err := parseUUID(r.PathValue("friendId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid friend id")
		return
	}

	result, err := h.service.Calculate(r.Context(), receiptID, friendID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}
