package service

import (
	"context"
	"errors"
	"grpc/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"stakeholders.com/dto"
	"stakeholders.com/model"
	"stakeholders.com/repo"
)

type AccountService struct {
	AccountRepo *repo.AccountRepository
}

func (service *AccountService) FindAll() ([]model.Account, error) {
	accounts, err := service.AccountRepo.FindAll()

	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (service *AccountService) Create(account *model.Account) error {
	err := service.AccountRepo.Create(account)

	if err != nil {
		return err
	}
	return nil
}

func (service *AccountService) FindAccount(accountId string) (bool, error) {
	err := service.AccountRepo.FindAccount(accountId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (service *AccountService) Login(loginDetails *dto.LoginDetailsDto) string {
	str, _ := service.AccountRepo.Login(loginDetails)

	return str
}

func (service *AccountService) BlockAccount(admin *model.Account, userId string) error {
	if admin == nil || admin.Role != "admin" {
		return errors.New("forbidden: admin only")
	}
	return service.AccountRepo.BlockAccount(userId)
}

func (service *AccountService) ParseToken(token string) *dto.Claims {
	return service.AccountRepo.ParseToken(token)
}

func (service *AccountService) LoginGRPC(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	loginDetails := dto.LoginDetailsDto{
		Email:    req.Email,
		Password: req.Password,
	}
	str := service.Login(&loginDetails)

	switch str {
	case "Nema":
		return nil, status.Errorf(codes.NotFound, "Nema naloga.")
	case "Lozinka":
		return nil, status.Errorf(codes.DataLoss, "Pogresna lozinka")
	case "Blok":
		return nil, status.Errorf(codes.PermissionDenied, "Nalog je blokiran")
	case "Token":
		return nil, status.Errorf(codes.Internal, "Problem prilikom kreiranja tokena")
	}

	return &proto.LoginResponse{Token: str}, nil
}

func (service *AccountService) RegisterGRPC(ctx context.Context, req *proto.RegisterAccountRequest) (*proto.RegisterAccountResponse, error) {
	account := model.Account{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
	}
	account.Id = uuid.New().String()

	if !account.IsValid() {
		return nil, status.Errorf(codes.InvalidArgument, "Account is not valid")
	}

	if account.Role != "guide" && account.Role != "tourist" {
		return nil, status.Errorf(codes.InvalidArgument, "Account role is not valid")
	}

	if err := service.Create(&account); err != nil {
		return nil, status.Errorf(codes.Internal, "Neuspesno kreiranje naloga: %v", err)
	}

	return &proto.RegisterAccountResponse{}, nil
}
