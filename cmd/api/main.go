package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	dbsql "github.com/takagiyuuki/grumble-back/internal/infrastructure/db/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("grumble-back api exited: %v", err)
	}
}

func run() error {
	dsn, err := getDatabaseURL()
	if err != nil {
		return err
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dbsql.RunMigrations(ctx, db, "migrations"); err != nil {
		return err
	}

	addr := getListenAddr()
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("starting grumble API on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func getListenAddr() string {
	if v := os.Getenv("GRUMBLE_HTTP_ADDR"); v != "" {
		return v
	}
	return ":8080"
}

func getDatabaseURL() (string, error) {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v, nil
	}
	return "", errors.New("DATABASE_URL must be set")
}
