package main

import "github.com/gorilla/mux"

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", h.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}/follow", h.FollowUser).Methods("POST")
	r.HandleFunc("/users/{id}/unfollow", h.UnfollowUser).Methods("POST")
	r.HandleFunc("/tweets", h.CreateTweet).Methods("POST")
	r.HandleFunc("/tweets/{id}", h.DeleteTweet).Methods("DELETE")
	r.HandleFunc("/tweets/{id}/like", h.LikeTweet).Methods("POST")
	r.HandleFunc("/feeds/{id}", h.GetFeed).Methods("GET")
}
