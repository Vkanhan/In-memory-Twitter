package main

import (
	"net/http"
	"strconv"
)

func getUserIDFromRequest(r *http.Request) int {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0
	}

	return userID
}
