package main

import (
	"encoding/json"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type build struct {
	BuildCommit string `json:"buildCommit"`
	BuildDate   string `json:"buildDate"`
	BuildTag    string `json:"buildTag"`
}

type health struct {
	InfluxDB string `json:"influxdb"`
	Status   int    `json:"-"`
}

type status struct {
	GoRoutines int `json:"goroutines"`
}

type weatherData struct {
	DewPoint           float64 `schema:"dewptf"`
	HumidityIndoor     float64 `schema:"indoorhumidity,required"`
	HumidityOutdoor    float64 `schema:"humidity,required"`
	PressureAbsolute   float64 `schema:"absbaromin,required"`
	PressureRelative   float64 `schema:"baromin,required"`
	Rain               float64 `schema:"rainin,required"`
	SolarRadiation     float64 `schema:"solarradiation,required"`
	TemperatureIndoor  float64 `schema:"indoortempf,required"`
	TemperatureOutdoor float64 `schema:"tempf,required"`
	UV                 float64 `schema:"UV,required"`
	WindChill          float64 `schema:"windchillf,required"`
	WindDirection      float64 `schema:"winddir,required"`
	WindGust           float64 `schema:"windgustmph,required"`
	WindSpeed          float64 `schema:"windspeedmph,required"`
}

func (w *weatherData) write() error {
	// Create batch points
	bp, err := client.NewBatchPoints(batchConfig)
	if err != nil {
		logError("%s: %s", ErrUnableToConnect, err)
		return ErrUnableToConnect
	}

	// Convert struct into map
	var fields map[string]interface{}
	wData, err := json.Marshal(w)
	if err != nil {
		logError("%s: %s", ErrUnableToDecode, err)
		return ErrUnableToDecode
	}
	err = json.Unmarshal(wData, &fields)
	if err != nil {
		logError("%s: %s", ErrUnableToDecode, err)
		return ErrUnableToDecode
	}

	// Create point
	pt, err := client.NewPoint("second", make(map[string]string), fields, time.Now())
	if err != nil {
		logError("%s: %s", ErrUnableToWrite, err)
		return ErrUnableToWrite
	}

	// Write batch
	bp.AddPoint(pt)
	err = influxdb.Write(bp)
	if err != nil {
		logError("%s: %s", ErrUnableToWrite, err)
		return ErrUnableToWrite
	}
	logDebug("point written")
	return nil
}
