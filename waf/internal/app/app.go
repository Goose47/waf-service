package app

import (
	"log/slog"
	grpcapp "waf-waf/internal/app/grpc"
	"waf-waf/internal/services"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	kafkaHost string,
	kafkaPort int,
	kafkaTopic string,
) *App {
	//todo add kafka service to wafService
	wafService := services.NewWAFService(log)

	grpcApp := grpcapp.New(log, wafService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
