package httpserver

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (hs *HttpServer) RegisterRoutes() {
	hs.apiRoutes()
	hs.router.GET("/", hs.Home)
}

func (hs *HttpServer) Home(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Write([]byte("Hi there!"))
}
