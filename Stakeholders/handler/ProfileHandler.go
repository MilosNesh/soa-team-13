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

func (handler *ProfileHandler) FindByAccountId(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	accountId := vars["accountId"]

	profile, err := handler.ProfileService.FindByAccountId(accountId)

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(profile)
}

func (handler *ProfileHandler) UpdateProfile(writer http.ResponseWriter, req *http.Request) {
	var profileDto dto.ProfileDto
	err := json.NewDecoder(req.Body).Decode(&profileDto)
	if err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedProfile, errr := handler.ProfileService.UpdateProfile(profileDto)

	if errr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(updatedProfile)

}
