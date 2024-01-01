package httpserver

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (hs *HttpServer) GetRealms(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Write([]byte("get yourself some realms"))
}
