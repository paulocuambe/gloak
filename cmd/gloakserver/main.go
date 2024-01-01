package main

import (
	"log"
	"os"

	"github.com/paulocuambe/gloak/internal/config"
	"github.com/paulocuambe/gloak/internal/db"
	"github.com/paulocuambe/gloak/internal/httpserver"
)

func main() {
	cfg, err, warnings := config.LoadConfig()

	if err != nil {
		if err, ok := err.(config.ConfigErr); ok {
			for _, err := range err.Errs {
				log.Println("validation error:", err)
			}
		} else {
			log.Println(err)
		}
		os.Exit(1)
	}

	if warnings != nil {
		if w, ok := warnings.(config.ConfigWarnigs); ok {
			for _, err := range w.Warnings {
				log.Println("warning:", err)
			}
		} else {
			log.Println(err)
		}
	}

	conn, err := db.ProvideDBConnection(cfg.DatabaseConfig)

	if err != nil {
		log.Fatalf("could not start database: %v\n", err)
	}

	server := httpserver.ProvideHttpServer(cfg, conn)

	log.Fatal(server.Start())
}
