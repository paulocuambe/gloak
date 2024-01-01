package httpserver

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/paulocuambe/gloak/internal/config"

	"github.com/julienschmidt/httprouter"
)

type HttpServer struct {
	DB *sql.DB

	cfg    *config.AppConfig
	server *http.Server
	router *httprouter.Router
}

func (h *HttpServer) Start() error {
	log.Printf("%v version %v", h.cfg.Name, h.cfg.Version)
	log.Printf("starting %v server on %v", h.cfg.Name, h.cfg.Addr())
	return h.server.ListenAndServe()
}

func ProvideHttpServer(cfg *config.AppConfig, db *sql.DB) *HttpServer {
	router := httprouter.New()
	router.RedirectTrailingSlash = true

	s := &http.Server{
		Addr:    cfg.Addr(),
		Handler: router,
	}

	hs := &HttpServer{cfg: cfg, server: s, router: router}
	hs.RegisterRoutes()

	return hs
}
