package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"followers.com/model"
)

type StakeholderService struct {
	BaseURL string
	Client  *http.Client
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

func (s *StakeholderService) GetProfileById(accountId string) (*model.Profile, error) {
	url := fmt.Sprintf("%s/accounts/%s/username", s.BaseURL, accountId)
	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var profile model.Profile
		profile.Id = accountId
		if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}
		return &profile, nil
	case http.StatusNotFound:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
