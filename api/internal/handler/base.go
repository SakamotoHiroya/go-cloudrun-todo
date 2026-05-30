package handler

import (
	"database/sql"

	"github.com/SakamotoHiroya/go-cloudrun-todo/db"
	"github.com/SakamotoHiroya/go-cloudrun-todo/internal/config"
)

type Handler struct {
	dbConn *sql.DB
	repo   *db.Queries
	cfg    *config.Config
}

func New(sqlDB *sql.DB, cfg *config.Config) *Handler {
	return &Handler{
		dbConn: sqlDB,
		repo:   db.New(sqlDB),
		cfg:    cfg,
	}
}
