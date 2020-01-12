package main

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
)

var decoder = schema.NewDecoder()

func getHealth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var h health
	err := checkHealth()
	if err != nil {
		h.InfluxDB = "fail"
		h.Status = 500
	} else {
		h.InfluxDB = "ok"
		h.Status = 200
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h.Status)
	err = json.NewEncoder(w).Encode(h)
	if err != nil {
		logError("%s: %s", ErrUnableToEncode, err)
	}
}

func getStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s status
	s.GoRoutines = runtime.NumGoroutine()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		logError("%s: %s", ErrUnableToEncode, err)
	}
}

func getVersion(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v := build{
		BuildCommit: buildCommit,
		BuildDate:   buildDate,
		BuildTag:    buildTag,
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		logError("%s: %s", ErrUnableToEncode, err)
	}
}

func getUpdateWeatherStation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Create model from query string
	var d weatherData
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&d, r.URL.Query())
	if err != nil {
		logError("%s: %s", ErrUnableToDecode, err)
		http.Error(w, ErrUnableToDecode.Error(), 400)
	}

	// Write model
	err = d.write()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
