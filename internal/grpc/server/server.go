package server

import (
	"GRPC_Calc/internal/services"
	genv1 "GRPC_Calc/proto/gen"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
	) (token string, err error)

	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

type Calc interface {
	CalculateExpression(ctx context.Context, expr string) (answer string, err error)
}

type serverAPI struct {
	genv1.UnimplementedCalcServer
	auth Auth
	calc Calc
}

func Register(gRPC *grpc.Server, auth Auth, calc Calc) {
	genv1.RegisterCalcServer(gRPC, &serverAPI{auth: auth, calc: calc})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *genv1.LoginRequest,
) (*genv1.LoginResponse, error) {

	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &genv1.LoginResponse{
		Token: token,
	}, nil
}

func validateLogin(req *genv1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateRegister(req *genv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *genv1.RegisterRequest,
) (*genv1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &genv1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Calculate(
	ctx context.Context,
	req *genv1.ExprRequest,
) (*genv1.ExprResponse, error) {
	ans, err := s.calc.CalculateExpression(ctx, req.GetExpr())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to calculate expression")
	}

	return &genv1.ExprResponse{
		Answer: fmt.Sprintf("%f", ans),
	}, nil
}
