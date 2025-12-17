package biz

import (
	"context"
	"myapp/internal/data/model"
)

type UserRepo interface {
	Create(ctx context.Context, u *model.User) error
	List(ctx context.Context) ([]model.User, error)
}

type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) Create(ctx context.Context, u *model.User) error {
	return uc.repo.Create(ctx, u)
}

func (uc *UserUsecase) List(ctx context.Context) ([]model.User, error) {
	return uc.repo.List(ctx)
}

