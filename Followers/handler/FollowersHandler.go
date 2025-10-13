package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"followers.com/service"
	"github.com/gorilla/mux"
)

type FollowHandler struct {
	FollowService *service.FollowService
}

func (h *FollowHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followerId := vars["followerId"]
	followedId := vars["followedId"]

	err := h.FollowService.Follow(followerId, followedId)
	if err != nil {
		http.Error(w, "Error following user", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User followed successfully"})
}

func (h *FollowHandler) IsFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followerId := vars["followerId"]
	followedId := vars["followedId"]

	log.Printf("=== IsFollowing Handler ===")
	log.Printf("Follower ID: %s", followerId)
	log.Printf("Followed ID: %s", followedId)

	isFollowing, err := h.FollowService.IsFollowing(followerId, followedId)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Error checking follow status", http.StatusInternalServerError)
		return
	}

	log.Printf("Result: isFollowing = %v", isFollowing)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"isFollowing": isFollowing})
}
