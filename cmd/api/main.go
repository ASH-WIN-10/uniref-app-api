package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
	}

	app.registerRoutes(e)

	logger.Info("Starting server", "version", version, "port", cfg.port, "env", cfg.env)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.port)))
}
