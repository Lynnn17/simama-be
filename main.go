package main

//go:generate go run github.com/swaggo/swag/cmd/swag init
//go:generate go run github.com/google/wire/cmd/wire

// @title Ajinomoto Gate System API
// @version 1.0
// @description API Documentation for Ajinomoto Gate System Backend
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"

	"lms-be/configs"
	"lms-be/infras"
	"lms-be/shared/logger"
	"lms-be/transport/http"
)

var config *configs.Config
var conn *sqlx.DB

type App struct {
	Config *configs.Config
	HTTP   *http.HTTP
	// Crontab crontab.CrontabService
}

func main() {
	// Initialize logger
	logger.InitLogger()

	// Initialize config
	config = configs.Get()

	// Set desired log level
	logger.SetLogLevel(config)

	// Wire everything up (from wire.go)
	app := InitializeService()

	// Run Crontab
	// if app.Crontab != nil {
	// 	go app.Crontab.Start(config)
	// }

	// Init DB read connection (optional, kalau dipakai untuk read-only ops)
	conn = infras.CreatePostgreSQLReadConn(*config)
	defer conn.Close()

	// Run HTTP server
	go app.HTTP.SetupAndServe()

	// Trap SIGINT & SIGTERM untuk graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
