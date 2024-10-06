package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"time"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(log *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	GRPCApp := grpcapp.New(log, port)
	return &App{
		GRPCApp: GRPCApp,
	}
}
