package service

import (
	"fmt"

	"followers.com/repo"
)

type FollowService struct {
	FollowRepo         *repo.FollowRepository
	StakeholderService *StakeholderService
}

func (s *FollowService) Follow(followerId, followedId string) error {
	exists, err := s.StakeholderService.FindAccount(followerId)
	if err != nil {
		return fmt.Errorf("error checking follower: %w", err)
	}
	if !exists {
		return fmt.Errorf("follower not found")
	}

	exists, err = s.StakeholderService.FindAccount(followedId)
	if err != nil {
		return fmt.Errorf("error checking followed user: %w", err)
	}
	if !exists {
		return fmt.Errorf("followed user not found")
	}

	return s.FollowRepo.Follow(followerId, followedId)
}
