package main

import (
	"context"
	"log"

	"github.com/paulocuambe/gloak/internal/config"
	"github.com/paulocuambe/gloak/internal/db"
	"github.com/paulocuambe/gloak/internal/httpserver"
)

func main() {
	cfg, errs, warnings := config.LoadConfig()

	if len(errs) > 0 {
		for _, err := range errs {
			log.Println("validation error:", err)
		}
	}

	if len(warnings) > 0 {
		for _, err := range warnings {
			log.Println("warning:", err)
		}
	}

	conn, err := db.ProvideDBConnection(cfg.DatabaseConfig)
	defer conn.DB.Close()
	if err != nil {
		log.Fatalf("could not start database: %v\n", err)
	}

	err = conn.RunMigrations(context.Background())
	if err != nil {
		log.Fatalf("error while running migrations: %v\n", err)
	}

	server := httpserver.ProvideHttpServer(cfg, conn)

	log.Fatal(server.Start())
}
