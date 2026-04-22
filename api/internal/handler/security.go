package handler

import(
	"fmt"
	"context"
	"errors"
	"os"
	"github.com/golang-jwt/jwt/v5"
	"github.com/SakamotoHiroya/go-cloudrun-todo/internal/api"
)

type SecurityHandler struct{}
type userKeyType struct{}

type UserInfo struct {
    id string
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

	secretKey := os.Getenv("SECRET_KEY")

    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, err
		}
		return secretKey, nil;
    })
    if err != nil || !token.Valid {
        return ctx, errors.New("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return ctx, errors.New("invalid claims")
    }

    userID := claims["sub"].(string)
	userName := claims["name"].(string)

    ctx = context.WithValue(ctx, userKey, UserInfo{
		id: userID,
		name: userName,
	})

    return ctx, nil
}