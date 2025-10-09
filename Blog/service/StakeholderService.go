package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StakeholderService struct {
	BaseURL string
	Client  *http.Client
}

type AccountDetails struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}

func (service *StakeholderService) FindAccount(accountId string) (bool, error) {
	url := fmt.Sprintf("%s/accounts/doesExists/%s", service.BaseURL, accountId)
	resp, err := service.Client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func (service *StakeholderService) GetAuthorUsername(accountId string) (string, error) {
	url := fmt.Sprintf("%s/accounts/%s/details", service.BaseURL, accountId)
	resp, err := service.Client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var details AccountDetails
		if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
			return "", fmt.Errorf("greška prilikom dekodiranja detalja naloga: %w", err)
		}
		return details.Username, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	return "", fmt.Errorf("neočekivani status kod od Stakeholders servisa: %d", resp.StatusCode)
}
