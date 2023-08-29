package app

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func SetLogrus(level string) {
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(logLevel)
	}

	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.SetOutput(os.Stdout)
}
