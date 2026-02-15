package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mamoanwar97/split-receipt/internal/database"
)

type FriendHandler struct {
	queries *database.Queries
}

func NewFriendHandler(q *database.Queries) *FriendHandler {
	return &FriendHandler{queries: q}
}

func (h *FriendHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CreateFriend] received request: method=%s content-type=%s", r.Method, r.Header.Get("Content-Type"))

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[CreateFriend] failed to decode request body: %v", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	log.Printf("[CreateFriend] decoded request: name=%q", req.Name)

	if req.Name == "" {
		log.Printf("[CreateFriend] validation failed: name is empty")
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	friend, err := h.queries.CreateFriend(r.Context(), req.Name)
	if err != nil {
		log.Printf("[CreateFriend] database error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create friend")
		return
	}
	log.Printf("[CreateFriend] created successfully: %+v", friend)
	writeJSON(w, http.StatusCreated, friend)
}

func (h *FriendHandler) List(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ListFriends] hit: method=%s url=%s", r.Method, r.URL.String())
	friends, err := h.queries.ListFriends(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list friends")
		return
	}
	writeJSON(w, http.StatusOK, friends)
}

func (h *FriendHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	friend, err := h.queries.GetFriend(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "friend not found")
		return
	}
	writeJSON(w, http.StatusOK, friend)
}

func (h *FriendHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.queries.DeleteFriend(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete friend")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parseUUID(s string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}
