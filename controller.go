package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
)

func getHealth(w http.ResponseWriter, r *http.Request) {
	var h health

	err := db.DB.Ping()
	if err != nil {
		logError("%s: %s", ErrUnableToConnect, err)
		h.PostgreSQL = false
		h.Status = 500
	} else {
		h.PostgreSQL = true
		h.Status = 200
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h.Status)

	err = json.NewEncoder(w).Encode(h)
	if err != nil {
		logError("%s: %s", ErrUnableToEncode, err)
	}
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	var s status
	s.GoRoutines = runtime.NumGoroutine()
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		logError("%s: %s", ErrUnableToEncode, err)
	}
}

func getVersion(w http.ResponseWriter, r *http.Request) {
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

func getUpdateWeatherStation(w http.ResponseWriter, r *http.Request) {
	var d weatherData
	var err error

	d.DewPoint, err = strconv.ParseFloat(r.URL.Query().Get("dewptf"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.HumidityIndoor, err = strconv.ParseFloat(r.URL.Query().Get("indoorhumidity"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.HumidityOutdoor, err = strconv.ParseFloat(r.URL.Query().Get("humidity"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.PressureAbsolute, err = strconv.ParseFloat(r.URL.Query().Get("absbaromin"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.PressureRelative, err = strconv.ParseFloat(r.URL.Query().Get("baromin"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.Rain, err = strconv.ParseFloat(r.URL.Query().Get("rainin"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.SolarRadiation, err = strconv.ParseFloat(r.URL.Query().Get("solarradiation"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.TemperatureIndoor, err = strconv.ParseFloat(r.URL.Query().Get("indoortempf"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.TemperatureOutdoor, err = strconv.ParseFloat(r.URL.Query().Get("tempf"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.UV, err = strconv.ParseFloat(r.URL.Query().Get("UV"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.WindChill, err = strconv.ParseFloat(r.URL.Query().Get("windchillf"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.WindDirection, err = strconv.ParseFloat(r.URL.Query().Get("winddir"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.WindGust, err = strconv.ParseFloat(r.URL.Query().Get("windgustmph"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	d.WindSpeed, err = strconv.ParseFloat(r.URL.Query().Get("windspeedmph"), 2)
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
		return
	}

	// Write model
	err = d.write()
	if err != nil {
		logError("%s: %s", ErrUnableToParse, err)
		http.Error(w, err.Error(), 500)
	}

	fmt.Fprintf(w, "success\n")
}
