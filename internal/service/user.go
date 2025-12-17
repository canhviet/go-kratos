package service

import (
	"context"
	pb "myapp/api/user/v1"
	"myapp/internal/biz"
	"myapp/internal/data/model"
)

type UserService struct {
	pb.UnimplementedUserServer
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	err := s.uc.Create(ctx, &model.User{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateReply{Message: "created"}, nil
}

func (s *UserService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListReply, error) {
	users, err := s.uc.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListReply{}
	for _, u := range users {
		resp.Items = append(resp.Items, &pb.UserItem{
			Id:    uint32(u.ID),
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return resp, nil
}
