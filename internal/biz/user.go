package biz

import (
	"context"
	"errors"
	"strconv"
	"time"

	"myapp/internal/data/model"
	"myapp/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo     *repository.UserRepo
	redis    *repository.RedisRepo // Assuming RedisRepo for token storage
	secret   string                // JWT secret
	tokenExp int                   // Token expiration in minutes
}

func NewAuthUsecase(repo *repository.UserRepo, redis *repository.RedisRepo, secret string, tokenExp int) *AuthUsecase {
	return &AuthUsecase{repo: repo, redis: redis, secret: secret, tokenExp: tokenExp}
}

func (uc *AuthUsecase) Register(ctx context.Context, username, password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	return uc.repo.Create(ctx, user)
}

func (uc *AuthUsecase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"id":       user.ID,
		"exp":      time.Now().Add(time.Minute * time.Duration(uc.tokenExp)).Unix(),
	})

	signedToken, err := token.SignedString([]byte(uc.secret))
	if err != nil {
		return "", err
	}

	// Store token in Redis for validation/revocation
	err = uc.redis.Set(ctx, "token:"+strconv.Itoa(int(user.ID)), signedToken, time.Minute*time.Duration(uc.tokenExp))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}