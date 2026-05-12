package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SakamotoHiroya/go-cloudrun-todo/db"
	"github.com/SakamotoHiroya/go-cloudrun-todo/internal/api"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

func (h *Handler) AuthWithGoogle(ctx context.Context, req *api.AuthWithGoogleReq) (api.AuthWithGoogleRes, error) {
	if h.cfg == nil {
		return nil, errors.New("config is not initialized")
	}
	if h.repo == nil {
		return nil, errors.New("repository is not initialized")
	}
	if h.cfg.JWTSecret == "" {
		return nil, errors.New("jwt secret is empty")
	}

	audience := ""
	if h.cfg.GoogleClientID != "" {
		audience = h.cfg.GoogleClientID
	}
	payload, err := idtoken.Validate(ctx, req.IdToken, audience)
	if err != nil {
		return nil, fmt.Errorf("failed to validate ID token: %w", err)
	}

	if payload.Issuer != "https://accounts.google.com" &&
		payload.Issuer != "accounts.google.com" {
		return nil, fmt.Errorf("invalid issuer: %s", payload.Issuer)
	}

	sub, _ := payload.Claims["sub"].(string)
	if sub == "" {
		return nil, errors.New("missing sub in google token")
	}

	user, err := h.repo.GetUserByGoogleSub(ctx, sub)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to load user: %w", err)
		}

		name, _ := payload.Claims["name"].(string)
		user, err = h.repo.CreateUserByGoogleSub(ctx, db.CreateUserByGoogleSubParams{
			GoogleSub: sub,
			Name:      sql.NullString{String: name, Valid: name != ""},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	claimName := ""
	if user.Name.Valid {
		claimName = user.Name.String
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        fmt.Sprintf("%d", user.ID),
		"user_id":    user.ID,
		"google_sub": user.GoogleSub,
		"name":       claimName,
		"exp":        expiresAt.Unix(),
		"iat":        time.Now().Unix(),
	})

	signedToken, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign jwt: %w", err)
	}

	cookie := (&http.Cookie{
		Name:     "Authorization",
		Value:    "Bearer " + signedToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	}).String()

	return &api.AuthWithGoogleOK{SetCookie: cookie}, nil
}
