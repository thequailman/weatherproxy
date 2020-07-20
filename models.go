package main

import (
	"time"
)

type build struct {
	BuildCommit string `json:"buildCommit"`
	BuildDate   string `json:"buildDate"`
	BuildTag    string `json:"buildTag"`
}

type health struct {
	PostgreSQL bool `json:"postgresql"`
	Status     int  `json:"-"`
}

type status struct {
	GoRoutines int `json:"goroutines"`
}

type weatherData struct {
	Created            time.Time `db:"created"`
	DewPoint           float64   `db:"dew_point"`
	HumidityIndoor     float64   `db:"humidity_indoor"`
	HumidityOutdoor    float64   `db:"humidity_outdoor"`
	PressureAbsolute   float64   `db:"pressure_absolute"`
	PressureRelative   float64   `db:"pressure_relative"`
	Rain               float64   `db:"rain"`
	SolarRadiation     float64   `db:"solar_radiation"`
	TemperatureIndoor  float64   `db:"temperature_indoor"`
	TemperatureOutdoor float64   `db:"temperature_outdoor"`
	UV                 float64   `db:"uv"`
	WindChill          float64   `db:"wind_chill"`
	WindDirection      float64   `db:"wind_direction"`
	WindGust           float64   `db:"wind_gust"`
	WindSpeed          float64   `db:"wind_speed"`
}

func (w *weatherData) write() error {
	if w.Created == (time.Time{}) {
		w.Created = time.Now()
	}

	_, err := db.NamedExec(`
insert into raw (
	created,
	dew_point,
	humidity_indoor,
	humidity_outdoor,
	pressure_absolute,
	pressure_relative,
	rain,
	solar_radiation,
	temperature_indoor,
	temperature_outdoor,
	uv,
	wind_chill,
	wind_direction,
	wind_gust,
	wind_speed
) values (
	:created,
	:dew_point,
	:humidity_indoor,
	:humidity_outdoor,
	:pressure_absolute,
	:pressure_relative,
	:rain,
	:solar_radiation,
	:temperature_indoor,
	:temperature_outdoor,
	:uv,
	:wind_chill,
	:wind_direction,
	:wind_gust,
	:wind_speed
)`, w)
	return err
}

func aggregateRaw() error {
	_, err := db.Exec(`
insert into hourly (
	created,
	dew_point_avg,
	dew_point_max,
	dew_point_min,
	humidity_indoor_avg,
	humidity_indoor_max,
	humidity_indoor_min,
	humidity_outdoor_avg,
	humidity_outdoor_max,
	humidity_outdoor_min,
	pressure_absolute_avg,
	pressure_absolute_max,
	pressure_absolute_min,
	pressure_relative_avg,
	pressure_relative_max,
	pressure_relative_min,
	rain_avg,
	rain_max,
	rain_min,
	solar_radiation_avg,
	solar_radiation_max,
	solar_radiation_min,
	temperature_indoor_avg,
	temperature_indoor_max,
	temperature_indoor_min,
	temperature_outdoor_avg,
	temperature_outdoor_max,
	temperature_outdoor_min,
	uv_avg,
	uv_max,
	uv_min,
	wind_chill_avg,
	wind_chill_max,
	wind_chill_min,
	wind_gust_avg,
	wind_gust_max,
	wind_gust_min,
	wind_direction_avg,
	wind_speed_avg,
	wind_speed_max,
	wind_speed_min
)	(select
	date_trunc('hour', raw.created) as hour,
	avg(raw.dew_point),
	max(raw.dew_point),
	min(raw.dew_point),
	avg(raw.humidity_indoor),
	max(raw.humidity_indoor),
	min(raw.humidity_indoor),
	avg(raw.humidity_outdoor),
	max(raw.humidity_outdoor),
	min(raw.humidity_outdoor),
	avg(pressure_absolute),
	max(pressure_absolute),
	min(pressure_absolute),
	avg(pressure_relative),
	max(pressure_relative),
	min(pressure_relative),
	avg(rain),
	max(rain),
	min(rain),
	avg(solar_radiation),
	max(solar_radiation),
	min(solar_radiation),
	avg(temperature_indoor),
	max(temperature_indoor),
	min(temperature_indoor),
	avg(temperature_outdoor),
	max(temperature_outdoor),
	min(temperature_outdoor),
	avg(uv),
	max(uv),
	min(uv),
	avg(wind_chill),
	max(wind_chill),
	min(wind_chill),
	avg(wind_gust),
	max(wind_gust),
	min(wind_gust),
	avg(wind_direction),
	avg(wind_speed),
	max(wind_speed),
	min(wind_speed)
from raw
where
	raw.created < date_trunc('hour', now())
group by
	hour
order by
	hour)
;
`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
delete from raw
where
	date_trunc('hour', created) < date_trunc('hour', now())
`)
	if err != nil {
		return err
	}

	_, err = db.Exec("refresh materialized view daily")
	if err != nil {
		return err
	}

	_, err = db.Exec("refresh materialized view monthly")
	return err
}
