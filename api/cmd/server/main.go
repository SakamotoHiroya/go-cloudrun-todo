package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/SakamotoHiroya/go-cloudrun-todo/internal/api"
	"github.com/SakamotoHiroya/go-cloudrun-todo/internal/config"
	"github.com/SakamotoHiroya/go-cloudrun-todo/internal/handler"
)

func main() {
	cfg := config.Load()

	sqlDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer sqlDB.Close()

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	h := handler.New(sqlDB, cfg)

	srv, err := api.NewServer(h, h)
	if err != nil {
		log.Fatalf("api server: %v", err)
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000"
	}

	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, "ok")
			return
		}
		withCORS(corsOrigin, withCookieBearer(srv)).ServeHTTP(w, r)
	})

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, root); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

// withCookieBearer copies Cookie "Authorization" into the Authorization header when missing,
// so browsers can send the token set by /auth/google.
func withCookieBearer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			if c, err := r.Cookie("Authorization"); err == nil && strings.TrimSpace(c.Value) != "" {
				r.Header.Set("Authorization", c.Value)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func withCORS(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
