package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/paulocuambe/gloak/internal/models"
)

func (hs *HttpServer) GetRealms(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	realm, err := hs.realmService.GetRealmByID(r.Context(), "dummy-data")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	d, err := json.Marshal(realm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("an unkwon error occured, try again later"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(d)
}

func (hs *HttpServer) CreateRealm(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	realm, err := hs.realmService.Create(r.Context(), &models.CreateRealmCommand{Id: "dummy-data-5", Name: "internal5"})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	d, err := json.Marshal(realm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("an unknown error occurred"))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(d))
}
