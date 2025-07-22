package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"stakeholders.com/dto"
	"stakeholders.com/service"
)

type ProfileHandler struct {
	ProfileService *service.ProfileService
}

func (handler *ProfileHandler) FindByUsername(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]

	profile, err := handler.ProfileService.FindByUsername(username)

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(profile)
}

func (handler *ProfileHandler) UpdateProfile(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]

	var profileDto dto.ProfileDto
	err := json.NewDecoder(req.Body).Decode(&profileDto)
	if err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedProfile, errr := handler.ProfileService.UpdateProfile(username, profileDto)

	if errr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(updatedProfile)

}
