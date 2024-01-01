package main

import (
	"log"
	"os"

	"github.com/paulocuambe/gloak/internal/config"
	"github.com/paulocuambe/gloak/internal/db"
	"github.com/paulocuambe/gloak/internal/httpserver"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		if err, ok := err.(*config.ErrConfig); ok {
			if len(err.Errors) > 0 {
				for _, err := range err.Errors {
					log.Println("validation error:", err)
				}
				os.Exit(1)
			}

			if len(err.Warnings) > 0 {
				for _, err := range err.Warnings {
					log.Println("warning:", err)
				}
			}
		}
	}

	conn, err := db.ProvideDBConnection(cfg.DatabaseConfig)

	if err != nil {
		log.Fatalf("could not start database: %v\n", err)
	}

	server := httpserver.ProvideHttpServer(cfg, conn)

	log.Fatal(server.Start())
}
