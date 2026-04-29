package handler

import (
	"github.com/SakamotoHiroya/go-cloudrun-todo/db"
	"github.com/SakamotoHiroya/go-cloudrun-todo/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	dbConn *pgxpool.Pool
	repo   *db.Queries
	cfg    *config.Config
}