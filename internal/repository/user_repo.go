package repository

import (
	"context"
	"myapp/internal/data"
	"myapp/internal/data/model"
)

type userRepo struct {
	data *data.Data
}

func NewUserRepo(data *data.Data) *userRepo {
	return &userRepo{data: data}
}

func (r *userRepo) Create(ctx context.Context, u *model.User) error {
	return r.data.DB.WithContext(ctx).Create(u).Error
}

func (r *userRepo) List(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.data.DB.WithContext(ctx).Find(&users).Error
	return users, err
}
