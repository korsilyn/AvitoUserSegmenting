package app

import (
	"avito-user-segmenting/modules/httpserver"
	"avito-user-segmenting/modules/postgresql"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	SetLogrus(log.WarnLevel)

	log.Info("Starting PostgeSQL")
	pg, err := postgresql.New(":8593", posgresql.MaxPoolSize(1))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	log.Info("Starting HTTP server")
	log.Debugf("Server port: %s", "8080")
	httpServer := httpserver.New(handler, httpserver.Port("8080"))

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
