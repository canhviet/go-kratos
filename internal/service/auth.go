// internal/service/auth.go
package service

import (
	"context"

	pb "myapp/api/auth/v1"
	"myapp/internal/biz"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	uc *biz.AuthUsecase
}

func NewAuthService(uc *biz.AuthUsecase) *AuthService {
	return &AuthService{uc: uc}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	err := s.uc.Register(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterReply{Message: "registered"}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	token, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{Token: token}, nil
}