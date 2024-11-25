package app

import (
	grpcapp "github.com/qu0ta/go-grpc-auth/internal/app/grpc"
	"github.com/qu0ta/go-grpc-auth/internal/services/auth"
	"github.com/qu0ta/go-grpc-auth/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, tokenTTL)
	grpcApp := grpcapp.New(log, grpcPort, authService)
	return &App{
		GRPCServer: grpcApp,
	}

}

func (a *App) Run() error {
	a.GRPCServer.MustRun()
	return nil
}
