package main

import (
	"testing"
	"time"
)

type HourlyData struct {
	Created     time.Time `db:"created"`
	DewPointAvg float64   `db:"dew_point_avg"`
	DewPointMax float64   `db:"dew_point_max"`
	DewPointMin float64   `db:"dew_point_min"`
}

func TestWeatherDataWrite(t *testing.T) {
	w := weatherData{
		Created: time.Now(),
	}

	for i := 1; i <= 50; i++ {
		w.Created = w.Created.Add(-30 * time.Minute)
		w.DewPoint += .5
		w.HumidityIndoor += .5
		w.HumidityOutdoor += .5
		w.PressureAbsolute += .5
		w.PressureRelative += .5
		w.Rain += .5
		w.SolarRadiation += .5
		w.TemperatureIndoor += .5
		w.TemperatureOutdoor += .5
		w.UV += .5
		w.WindChill += .5
		w.WindDirection += .5
		w.WindGust += .5
		w.WindSpeed += .5

		err := w.write()
		if err != nil {
			t.Error(err)
		}
	}
}

func TestAggregateRaw(t *testing.T) {
	err := aggregateRaw()
	if err != nil {
		t.Error(err)
	}

	var hourly []HourlyData

	err = db.Select(&hourly, "select created, dew_point_avg, dew_point_max, dew_point_min from hourly")
	if err != nil {
		t.Error(err)
	}

	if len(hourly) < 25 {
		t.Error("hourly data should be at least 25 records")
	}

	first := hourly[0]

	if first.DewPointAvg == 0 {
		t.Errorf("got %f, want > 0", first.DewPointAvg)
	}

	if first.DewPointMax == 0 {
		t.Errorf("got %f, want > 0", first.DewPointMax)
	}

	if first.DewPointMin == 0 {
		t.Errorf("got %f, want > 0", first.DewPointMin)
	}

	var data []weatherData

	err = db.Select(&data, "select created from raw")
	if err != nil {
		t.Error(err)
	}

	if len(data) > 20 {
		t.Errorf("raw data should be less than 20, got %d", len(data))
	}

	data = []weatherData{}

	err = db.Select(&data, "select created from daily")
	if err != nil {
		t.Error(err)
	}

	if len(data) > 2 || len(data) == 0 {
		t.Errorf("daily data should be less than 3, got %d", len(data))
	}

	data = []weatherData{}

	err = db.Select(&data, "select created from monthly")
	if err != nil {
		t.Error(err)
	}

	if len(data) > 2 || len(data) == 0 {
		t.Errorf("monthly data should be less than 3, got %d", len(data))
	}
}
