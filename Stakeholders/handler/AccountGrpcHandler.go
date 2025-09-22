package handler

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stakeholders.com/dto"
	"stakeholders.com/model"
	"stakeholders.com/proto/stakeholders"
	"stakeholders.com/service"
)

type AccountGrpcHandler struct {
	stakeholders.UnimplementedStakeholdersServiceServer
	Service *service.AccountService
}

func (handler AccountGrpcHandler) RegisterAccount(ctx context.Context, req *stakeholders.RegisterAccountRequest) (*stakeholders.RegisterAccountResponse, error) {
	log.Printf("Register sa grpc")
	account := model.Account{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
	}
	log.Printf("Login attempt: Email=%s, Password=%s, Result=%s, Username=%s, Role=%s ", req.Email, req.Password, req.Username, req.Role)
	account.Id = uuid.New().String()

	if !account.IsValid() {
		return nil, status.Errorf(codes.InvalidArgument, "Account is not valid")
	}

	if account.Role != "guide" && account.Role != "tourist" {
		return nil, status.Errorf(codes.InvalidArgument, "Account role is not valid")
	}

	if err := handler.Service.Create(&account); err != nil {
		return nil, status.Errorf(codes.Internal, "Neuspesno kreiranje naloga: %v", err)
	}

	return &stakeholders.RegisterAccountResponse{}, nil
}

func (handler AccountGrpcHandler) Login(ctx context.Context, req *stakeholders.LoginRequest) (*stakeholders.LoginResponse, error) {
	log.Printf("Login sa grpc")

	loginDetails := dto.LoginDetailsDto{
		Email:    req.Email,
		Password: req.Password,
	}
	str := handler.Service.Login(&loginDetails)

	log.Printf("Login attempt: Email=%s, Password=%s, Result=%s", req.Email, req.Password, str)

	switch str {
	case "Nema":
		return nil, status.Errorf(codes.NotFound, "Nema naloga.")
	case "Lozinka":
		return nil, status.Errorf(codes.InvalidArgument, "Pogresna lozinka")
	case "Blok":
		return nil, status.Errorf(codes.PermissionDenied, "Nalog je blokiran")
	case "Token":
		return nil, status.Errorf(codes.Internal, "Problem prilikom kreiranja tokena")
	}

	return &stakeholders.LoginResponse{Token: str}, nil
}
