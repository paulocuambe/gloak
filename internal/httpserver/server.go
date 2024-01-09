package httpserver

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/paulocuambe/gloak/internal/config"
	"github.com/paulocuambe/gloak/internal/db"
	"github.com/paulocuambe/gloak/internal/models"
	"github.com/paulocuambe/gloak/internal/services/realm"

	"github.com/julienschmidt/httprouter"
)

type HttpServer struct {
	store *db.DB

	cfg    *config.AppConfig
	server *http.Server
	router *httprouter.Router

	realmService models.RealmService
}

func (h *HttpServer) DB() *sql.DB {
	return h.store.DB
}

func (h *HttpServer) Start() error {
	log.Printf("%v version %v", h.cfg.Name, h.cfg.Version)
	log.Printf("starting %v server on %v", h.cfg.Name, h.cfg.Addr())
	return h.server.ListenAndServe()
}

func ProvideHttpServer(cfg *config.AppConfig, db *db.DB) *HttpServer {
	router := httprouter.New()
	router.RedirectTrailingSlash = true

	s := &http.Server{
		Addr:    cfg.Addr(),
		Handler: router,
	}

	hs := &HttpServer{cfg: cfg, server: s, router: router}
	hs.realmService = realm.ProvideService(db)
	hs.RegisterRoutes()

	return hs
}
