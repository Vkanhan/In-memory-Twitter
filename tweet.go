package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	ID        int
	Username  string
	Followers map[int]bool
	Following map[int]bool
	mu        sync.RWMutex
}

type Tweet struct {
	ID        int
	UserID    int
	Content   string
	Timestamp time.Time
	Likes     int
}

func (h *Handler) CreateTweet(w http.ResponseWriter, r *http.Request) {
	var req CreateTweetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID <= 0 {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tweet, err := h.service.CreateTweet(userID, req.Content)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    tweet,
	})
}

func (h *Handler) DeleteTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tweet ID")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID <= 0 {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.service.DeleteTweet(userID, tweetID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    "Tweet successfully deleted",
	})
}

func (h *Handler) LikeTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tweet ID")
		return
	}

	if err := h.service.LikeTweet(tweetID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    "Tweet liked successfully",
	})
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	query := r.URL.Query()
	limitStr := query.Get("limit")
	limit := 20
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	tweets, err := h.service.GetFeed(userID, limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    tweets,
	})
}
