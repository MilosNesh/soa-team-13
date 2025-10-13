package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type FollowerService struct {
	BaseURL string
	Client  *http.Client
}

func (service *FollowerService) IsFollowing(followerId, followedId string) (bool, error) {
	// LOG 1: Na ulazu u funkciju
	log.Printf("INFO: IsFollowing pozvan. followerId: %s, followedId: %s", followerId, followedId)

	url := fmt.Sprintf("%sisFollowing/%s/%s", service.BaseURL, followerId, followedId)

	// LOG 2: Provera generisanog URL-a
	log.Printf("INFO: IsFollowing gađa URL: %s", url)

	resp, err := service.Client.Get(url)
	if err != nil {
		// LOG 3: Pad pri HTTP pozivu
		log.Printf("ERROR: IsFollowing - Greška pri HTTP GET pozivu: %v", err)
		return false, err
	}
	defer resp.Body.Close()

	// LOG 4: Provera Status Koda
	log.Printf("INFO: IsFollowing - Dobijen HTTP status: %d", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		var result map[string]bool
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			// Log već postoji ovde: log.Printf("Greska je: %v", err)
			return false, err
		}
		log.Printf("IsFollowing result: %v", result["isFollowing"])

		// LOG 5: Uspešno dekodiranje
		log.Printf("SUCCESS: IsFollowing - Vraća rezultat: %v", result["isFollowing"])
		return result["isFollowing"], nil
	}

	// LOG 6: Pad zbog neočekivanog status koda (npr. 404, 500)
	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
