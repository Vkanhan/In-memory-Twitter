package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *TwitterService
}

type CreateUserRequest struct {
	Username string `json:"username"`
}

type CreateTweetRequest struct {
	Content string `json:"content"`
}

func NewHandler(service *TwitterService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.service.CreateUser(req.Username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    user,
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	profile, err := h.service.GetUserProfile(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    profile,
	})
}

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID <= 0 {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.service.FollowUser(userID, targetID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    "Successfully followed user",
	})
}

func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid target user ID")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID <= 0 {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.service.UnfollowUser(userID, targetID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    "Successfully unfollowed user",
	})
}
