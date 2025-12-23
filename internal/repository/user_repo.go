package repository

import (
	"context"

	"myapp/internal/data"
	"myapp/internal/data/model"
)

type UserRepo struct {
	data *data.Data
}

func NewUserRepo(data *data.Data) *UserRepo {
	return &UserRepo{data: data}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.data.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.data.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}