package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/SakamotoHiroya/go-cloudrun-todo/internal/api"
	"github.com/golang-jwt/jwt/v5"
)

type userKeyType struct{}

type UserInfo struct {
	id   string
	name string
}

var userKey = userKeyType{}

func (h *Handler) HandleBearerAuth(
	ctx context.Context,
	operationName string,
	t api.BearerAuth,
) (context.Context, error) {

	tokenStr := t.Token
	if tokenStr == "" {
		return ctx, errors.New("missing token")
	}

	secretKey := ""
	if h.cfg != nil {
		secretKey = h.cfg.JWTSecret
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, err
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return ctx, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx, errors.New("invalid claims")
	}

	userID, _ := claims["sub"].(string)
	userName, _ := claims["name"].(string)
	if userID == "" {
		return ctx, errors.New("invalid token subject")
	}

	ctx = context.WithValue(ctx, userKey, UserInfo{
		id:   userID,
		name: userName,
	})

	return ctx, nil
}
