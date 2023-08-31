package app

import (
	"avito-user-segmenting/config"
	v1 "avito-user-segmenting/core/controller/http/v1"
	"avito-user-segmenting/core/repo"
	"avito-user-segmenting/core/service"
	"avito-user-segmenting/modules/httpserver"
	"avito-user-segmenting/modules/postgresql"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// @title           Avito User Segmenting
// @version         1.0
// @description     This is a service for managing user slugs.

// @contact.name   Eugene Gladkov
// @contact.email  gladkov.ea@mail.com

// @host      localhost:8080
// @BasePath  /

func Run(configPath string) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	SetLogrus(cfg.Log.Level)

	log.Info("Starting PostgeSQL")
	pg, err := postgresql.New(cfg.PG.URL, postgresql.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	log.Info("Initializing deps")
	repositories := repo.NewRepositories(pg)
	deps := service.ServicesDependencies{
		Repos: repositories,
	}
	services := service.NewServices(deps)
	
	log.Info("Initializing ECHO and router")
	handler := echo.New()
	v1.NewRouter(handler, services)

	log.Info("Starting HTTP server")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	log.Info("Configuring graceful shutdown")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	log.Info("Shutting down")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
