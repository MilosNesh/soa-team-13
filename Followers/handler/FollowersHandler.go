package handler

import (
	"encoding/json"
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
