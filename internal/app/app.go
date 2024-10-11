package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	authservice "sso/internal/services/auth"
	"sso/internal/storage/sqlite"

	"time"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(log *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := authservice.New(log, storage, storage, storage, tokenTTL)
	GRPCApp := grpcapp.New(log, authService, port)
	return &App{
		GRPCApp: GRPCApp,
	}
}
