package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/paulocuambe/gloak/internal/models"
)

func (hs *HttpServer) GetRealms(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	realm, err := hs.realmService.GetRealms(context.Background())

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

func (hs *HttpServer) GetRealmById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("realmId")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("provide a valid realmId parameter"))
		return
	}

	realm, err := hs.realmService.GetRealmByID(r.Context(), id)

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
	d := json.NewDecoder(r.Body)
	var cmd models.CreateRealmCommand

	if err := d.Decode(&cmd); err != nil {
		log.Printf("%#v", err)
		if err == io.EOF {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("the body can't be empty"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the json has an invalid  format"))
		return
	}

	errs := cmd.Validate()

	if len(errs) > 0 {
		// new http api error
		// return in json format
		b, err := json.Marshal(errs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("an unexpected error occurred"))
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	realm, err := hs.realmService.Create(r.Context(), &cmd)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	bs, err := realm.ToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("an unknown error occurred"))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(bs))
}

func jsonError(err error) {
	if err == io.EOF {
		// json can't be empty
	}
}
