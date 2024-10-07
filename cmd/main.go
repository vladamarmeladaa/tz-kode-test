package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	config "tz_kode/internal/config"
	createNote "tz_kode/internal/handlers/create_note"
	getNote "tz_kode/internal/handlers/get_notes"
	signInUser "tz_kode/internal/handlers/sign_in_user"
	signUpUser "tz_kode/internal/handlers/sign_up_user"
	auth "tz_kode/internal/services/auth"
	speller "tz_kode/internal/services/speller"
	storage "tz_kode/internal/storage"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	logger := setupLogger()
	logger.Info(
		"starting tz_kode",
	)

	storage, err := storage.New(connStr)
	if err != nil {
		logger.Error("failed to init storage", err)
		os.Exit(1)
	}

	speller := speller.New(cfg.Speller)
	v := validator.New(validator.WithRequiredStructEnabled())

	r := mux.NewRouter()

	r.HandleFunc("/notes", auth.AuthMiddleware(createNote.New(logger, storage, speller, v), logger, storage)).Methods("POST")
	r.HandleFunc("/notes", auth.AuthMiddleware(getNote.New(logger, storage), logger, storage)).Methods("GET")

	r.HandleFunc("/sign-up", signUpUser.New(logger, storage, v)).Methods("POST")
	r.HandleFunc("/login", signInUser.New(logger, storage, v)).Methods("POST")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		logger.Error("failed to start server")
	}
}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	return log
}
