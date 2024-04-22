package app

import (
	"GRPC_Calc/internal/app/grpc"
	"GRPC_Calc/internal/services"
	"GRPC_Calc/internal/storage/sqlite"
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

	authService := services.NewAuth(log, storage, storage, tokenTTL)

	calcService := services.NewCalc(log, storage, storage)

	grpcApp := grpcapp.New(log, authService, calcService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
