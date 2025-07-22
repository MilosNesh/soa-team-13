package service

import (
	"stakeholders.com/dto"
	"stakeholders.com/model"
	"stakeholders.com/repo"
)

type ProfileService struct {
	ProfileRepo *repo.ProfileRepository
}

func (service *ProfileService) FindByUsername(username string) (*model.Profile, error) {
	profile, err := service.ProfileRepo.FindByUsername(username)

	if err != nil {
		return nil, err
	}

	return profile, nil

}

func (service *ProfileService) UpdateProfile(username string, profileDto dto.ProfileDto) (*model.Profile, error) {
	profile, err := service.FindByUsername(username)

	if err != nil {
		return nil, err
	}

	if profileDto.Name != nil {
		profile.Name = *profileDto.Name
	}
	if profileDto.Surname != nil {
		profile.Surname = *profileDto.Surname
	}
	if profileDto.ProfilePicture != nil {
		profile.ProfilePicture = *profileDto.ProfilePicture
	}
	if profileDto.Bio != nil {
		profile.Bio = *profileDto.Bio
	}
	if profileDto.Motto != nil {
		profile.Motto = *profileDto.Motto
	}

	err = service.ProfileRepo.Save(profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
