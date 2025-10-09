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

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type AccountGrpcHandler struct {
	stakeholders.UnimplementedStakeholdersServiceServer
	Service     *service.AccountService
	Tracer      *sdktrace.TracerProvider
	ServiceName string
}

func (handler AccountGrpcHandler) RegisterAccount(ctx context.Context, req *stakeholders.RegisterAccountRequest) (*stakeholders.RegisterAccountResponse, error) {
	_, span := (trace.TracerProvider)(handler.Tracer).Tracer(handler.ServiceName).Start(ctx, "grpc-register-account")
	defer func() { span.End() }()

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
		span.RecordError(status.Errorf(codes.InvalidArgument, "Account is not valid"), trace.WithAttributes())
		return nil, status.Errorf(codes.InvalidArgument, "Account is not valid")
	}

	if account.Role != "guide" && account.Role != "tourist" {
		span.RecordError(status.Errorf(codes.InvalidArgument, "Account role is not valid"), trace.WithAttributes())
		return nil, status.Errorf(codes.InvalidArgument, "Account role is not valid")
	}

	if err := handler.Service.Create(&account); err != nil {
		span.RecordError(err, trace.WithAttributes())
		return nil, status.Errorf(codes.Internal, "Neuspesno kreiranje naloga: %v", err)
	}

	return &stakeholders.RegisterAccountResponse{}, nil
}

func (handler AccountGrpcHandler) Login(ctx context.Context, req *stakeholders.LoginRequest) (*stakeholders.LoginResponse, error) {
	_, span := (trace.TracerProvider)(handler.Tracer).Tracer(handler.ServiceName).Start(ctx, "grpc-login")
	defer func() { span.End() }()

	log.Printf("Login sa grpc")

	loginDetails := dto.LoginDetailsDto{
		Email:    req.Email,
		Password: req.Password,
	}
	str := handler.Service.Login(&loginDetails)

	log.Printf("Login attempt: Email=%s, Password=%s, Result=%s", req.Email, req.Password, str)

	switch str {
	case "Nema":
		span.RecordError(status.Errorf(codes.NotFound, "Nema naloga."), trace.WithAttributes())
		return nil, status.Errorf(codes.NotFound, "Nema naloga.")
	case "Lozinka":
		span.RecordError(status.Errorf(codes.InvalidArgument, "Pogresna lozinka"), trace.WithAttributes())
		return nil, status.Errorf(codes.InvalidArgument, "Pogresna lozinka")
	case "Blok":
		span.RecordError(status.Errorf(codes.PermissionDenied, "Nalog je blokiran"), trace.WithAttributes())
		return nil, status.Errorf(codes.PermissionDenied, "Nalog je blokiran")
	case "Token":
		span.RecordError(status.Errorf(codes.Internal, "Problem prilikom kreiranja tokena"), trace.WithAttributes())
		return nil, status.Errorf(codes.Internal, "Problem prilikom kreiranja tokena")
	}

	return &stakeholders.LoginResponse{Token: str}, nil
}
