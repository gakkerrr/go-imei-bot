package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go-imei-bot/internal/config"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting imei-bot",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017/imei-bot"))
	if err != nil {
		log.Error("failed to connect to MongoDB", slog.String("error", err.Error()))
		return
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error("failed to disconnect from MongoDB", slog.String("error", err.Error()))
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Error("failed to ping MongoDB", slog.String("error", err.Error()))
		return
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":8000", r)

	// disconnect server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
