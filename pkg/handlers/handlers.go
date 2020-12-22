package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/eze8789/urlshtn-go/pkg/database"
	"github.com/eze8789/urlshtn-go/pkg/database/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type handler struct {
	storage database.Storage
}

func Routes(storage database.Storage) http.Handler {
	rtr := mux.NewRouter()
	h := handler{storage}
	rtr.HandleFunc("/load", responseHand(h.loadURL))
	rtr.PathPrefix("/").HandlerFunc(h.redirect).Methods(http.MethodGet)

	return rtr
}

type response struct {
	Data interface{} `json:"response"`
}

func responseHand(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, code, err := h(w, r)
		if err != nil {
			data = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(response{data}); err != nil {
			logrus.Errorf("could not process request: %v", err)
		}
	}
}

func (h handler) loadURL(w io.Writer, r *http.Request) (interface{}, int, error) {
	in := models.InOut{}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		logrus.Errorf("invalid payload")
		return nil, http.StatusBadRequest, err
	}
	defer r.Body.Close()

	url := strings.TrimSpace(in.URL)
	id, err := h.storage.Insert(url)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	logrus.Infof("saved from %s the URL: %s", r.RemoteAddr, url)
	return "http://localhost:8080/" + strconv.Itoa(*id), http.StatusCreated, nil
}

func (h handler) redirect(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Path[len("/"):]
	url, err := h.storage.RetrieveURL(u)
	if err != nil {
		errResp := new(bytes.Buffer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(errResp).Encode(response{Data: response{"URL not found"}})
		w.Write(errResp.Bytes())
		logrus.Errorf("request from %s, ID %s not found", r.RemoteAddr, u)
		return
	}

	logrus.Infof("redirect customer %s to %s", r.RemoteAddr, url)
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}
