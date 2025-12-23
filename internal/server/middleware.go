package server

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"myapp/internal/repository"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string, redisRepo *repository.RedisRepo) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// Skip for login and register
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.New("transport not found")
			}
			httpTr, ok := tr.(http.Transporter)
			if !ok {
				return nil, errors.New("http transport not found")
			}
			path := httpTr.Request().URL.Path
			if path == "/auth/login" || path == "/auth/register" {
				return handler(ctx, req)
			}

			// Get token from header
			authHeader := httpTr.Request().Header.Get("Authorization")
			if authHeader == "" {
				return nil, errors.New("authorization header missing")
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse JWT
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return nil, errors.New("invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, errors.New("invalid claims")
			}

			userID := int(claims["id"].(float64))
			// Check if token exists in Redis (for revocation)
			storedToken, err := redisRepo.Get(ctx, "token:"+strconv.Itoa(userID))
			if err != nil || storedToken != tokenStr {
				return nil, errors.New("token revoked or invalid")
			}

			// Add user info to context if needed
			return handler(ctx, req)
		}
	}
}