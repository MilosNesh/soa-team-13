package service

import (
	"errors"

	"gorm.io/gorm"
	"stakeholders.com/dto"
	"stakeholders.com/model"
	"stakeholders.com/repo"
)

type ProfileService struct {
	ProfileRepo *repo.ProfileRepository
}

func (service *ProfileService) FindByAccountId(accountId string) (*model.Profile, error) {
	profile, err := service.ProfileRepo.FindByAccountId(accountId)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (service *ProfileService) UpdateProfile(profileDto dto.ProfileDto) (*model.Profile, error) {
	var profile *model.Profile
	var err error

	profile, err = service.FindByAccountId(profileDto.AccountId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		profile = &model.Profile{
			AccountId: profileDto.AccountId,
		}
	} else if err != nil {
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

func (service *ProfileService) HasSufficientBalance(accountId string, amount float64) (bool, error) {
	profile, err := service.ProfileRepo.FindByAccountId(accountId)
	if err != nil {
		return false, err
	}
	return float64(profile.Balance) >= amount, nil
}

func (service *ProfileService) DeductBalance(accountId string, amount float64) error {
	return service.ProfileRepo.UpdateBalance(accountId, -amount)
}

func (service *ProfileService) RefundBalance(accountId string, amount float64) error {
	return service.ProfileRepo.UpdateBalance(accountId, +amount)
}
