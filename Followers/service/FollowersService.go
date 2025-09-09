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
	follower, err := s.StakeholderService.GetProfileById(followerId)
	if err != nil {
		return fmt.Errorf("error checking follower: %w", err)
	}
	if follower == nil {
		return fmt.Errorf("follower not found")
	}

	followed, err := s.StakeholderService.GetProfileById(followedId)
	if err != nil {
		return fmt.Errorf("error checking followed user: %w", err)
	}
	if followed == nil {
		return fmt.Errorf("followed not found")
	}

	return s.FollowRepo.Follow(follower, followed)
}
