package service

import (
	"fmt"
	"net/http"
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
