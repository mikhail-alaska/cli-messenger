package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikhail-alaska/cli-messenger/server/internal/config"
	login "github.com/mikhail-alaska/cli-messenger/server/internal/http-server/handlers/auth"
	"github.com/mikhail-alaska/cli-messenger/server/internal/http-server/handlers/messages"
	"github.com/mikhail-alaska/cli-messenger/server/internal/http-server/handlers/users"
	"github.com/mikhail-alaska/cli-messenger/server/internal/http-server/middleware/auth"
	"github.com/mikhail-alaska/cli-messenger/server/internal/http-server/middleware/logger"
	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/logger/handlers/slogpretty"
	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/logger/sl"
	"github.com/mikhail-alaska/cli-messenger/server/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	log.Info("starting url shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/user", users.NewUser(log, storage))

	router.Post("/login", login.LoginHandler(log, storage))

	router.Group(func(r chi.Router) {
		r.Use(auth.JWTAuth)
		r.Get("/users", users.AllUsers(log, storage))
		r.Get("/users/openkey", users.OpenKeyByUserName(log, storage))
		r.Get("/chats", messages.GetChats(log, storage))
        r.Post("/message", messages.NewMessage(log, storage))
		r.Get("/message", messages.GetMessages(log, storage))
	})
	// TODO run server
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("couldnt start server", sl.Err(err))
		os.Exit(1)
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = setupPrettySlog()

	case envProd:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		))
	default:
		log = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
