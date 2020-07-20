package main

import (
	"fmt"
	"sort"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

type migration struct {
	Migration string
	Version   int
}

func initDB() (err error) {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?timezone=UTC&sslmode=%s", c.PostgreSQL.Username, c.PostgreSQL.Password, c.PostgreSQL.Hostname, c.PostgreSQL.Port, c.PostgreSQL.Database, c.PostgreSQL.SSLMode)
	db, err = sqlx.Open("postgres", dataSourceName)
	if err != nil {
		logError(err.Error())
		return err
	}

	err = db.DB.Ping()
	return err
}

func migrateDB() error {
	_, err := db.Exec("create table if not exists migration (app text primary key, version int)")
	if err != nil {
		logError(err.Error())
		return err
	}

	var version int
	err = db.Get(&version, "select version from migration where app = 'weatherproxy'")
	if err != nil && err.Error() != "sql: no rows in result set" {
		logError(err.Error())
		return err
	}

	migrations := []migration{
		{Migration: fmt.Sprintf(`
create table if not exists raw (
	id int generated always as identity primary key,
	created timestamp with time zone not null,
	dew_point real not null,
	humidity_indoor real not null,
	humidity_outdoor real not null,
	pressure_absolute real not null,
	pressure_relative real not null,
	rain real not null,
	solar_radiation real not null,
	temperature_indoor real not null,
	temperature_outdoor real not null,
	uv real not null,
	wind_chill real not null,
	wind_direction real not null,
	wind_gust real not null,
	wind_speed real not null
);
create table if not exists hourly (
	id int generated always as identity primary key,
	created timestamp with time zone not null default current_timestamp,
	dew_point_avg real not null,
	dew_point_max real not null,
	dew_point_min real not null,
	humidity_indoor_avg real not null,
	humidity_indoor_max real not null,
	humidity_indoor_min real not null,
	humidity_outdoor_avg real not null,
	humidity_outdoor_max real not null,
	humidity_outdoor_min real not null,
	pressure_absolute_avg real not null,
	pressure_absolute_max real not null,
	pressure_absolute_min real not null,
	pressure_relative_avg real not null,
	pressure_relative_max real not null,
	pressure_relative_min real not null,
	rain_avg real not null,
	rain_max real not null,
	rain_min real not null,
	solar_radiation_avg real not null,
	solar_radiation_max real not null,
	solar_radiation_min real not null,
	temperature_indoor_avg real not null,
	temperature_indoor_max real not null,
	temperature_indoor_min real not null,
	temperature_outdoor_avg real not null,
	temperature_outdoor_max real not null,
	temperature_outdoor_min real not null,
	uv_avg real not null,
	uv_max real not null,
	uv_min real not null,
	wind_chill_avg real not null,
	wind_chill_max real not null,
	wind_chill_min real not null,
	wind_gust_avg real not null,
	wind_gust_max real not null,
	wind_gust_min real not null,
	wind_direction_avg real not null,
	wind_speed_avg real not null,
	wind_speed_max real not null,
	wind_speed_min real not null
);
create materialized view if not exists daily as
	select
		date_trunc('day', hourly.created at time zone '%[1]s') as created,
		avg(dew_point_avg) as dew_point_avg,
		max(dew_point_max) as dew_point_max,
		min(dew_point_min) as dew_point_min,
		avg(humidity_indoor_avg) as humidity_indoor_avg,
		max(humidity_indoor_max) as humidity_indoor_max,
		min(humidity_indoor_min) as humidity_indoor_min,
		avg(humidity_outdoor_avg) as humidity_outdoor_avg,
		max(humidity_outdoor_max) as humidity_outdoor_max,
		min(humidity_outdoor_min) as humidity_outdoor_min,
		avg(pressure_absolute_avg) as pressure_absolute_avg,
		max(pressure_absolute_max) as pressure_absolute_max,
		min(pressure_absolute_min) as pressure_absolute_min,
		avg(pressure_relative_avg) as pressure_relative_avg,
		max(pressure_relative_max) as pressure_relative_max,
		min(pressure_relative_min) as pressure_relative_min,
		avg(rain_avg) as rain_avg,
		max(rain_max) as rain_max,
		min(rain_min) as rain_min,
		avg(solar_radiation_avg) as solar_radiation_avg,
		max(solar_radiation_max) as solar_radiation_max,
		min(solar_radiation_min) as solar_radiation_min,
		avg(temperature_indoor_avg) as temperature_indoor_avg,
		max(temperature_indoor_max) as temperature_indoor_max,
		min(temperature_indoor_min) as temperature_indoor_min,
		avg(temperature_outdoor_avg) as temperature_outdoor_avg,
		max(temperature_outdoor_max) as temperature_outdoor_max,
		min(temperature_outdoor_min) as temperature_outdoor_min,
		avg(uv_avg) as uv_avg,
		max(uv_max) as uv_max,
		min(uv_min) as uv_min,
		avg(wind_chill_avg) as wind_chill_avg,
		max(wind_chill_max) as wind_chill_max,
		min(wind_chill_min) as wind_chill_min,
		avg(wind_gust_avg) as wind_gust_avg,
		max(wind_gust_max) as wind_gust_max,
		min(wind_gust_min) as wind_gust_min,
		avg(wind_direction_avg) as wind_direction_avg,
		avg(wind_speed_avg) as wind_speed_avg,
		max(wind_speed_max) as wind_speed_max,
		min(wind_speed_min) as wind_speed_min
	from hourly
	group by
		date_trunc('day', hourly.created at time zone '%[1]s')
	order by
		date_trunc('day', hourly.created at time zone '%[1]s')
;
create materialized view if not exists monthly as
	select
		date_trunc('month', daily.created at time zone '%[1]s') as created,
		avg(dew_point_avg) as dew_point_avg,
		max(dew_point_max) as dew_point_max,
		min(dew_point_min) as dew_point_min,
		avg(humidity_indoor_avg) as humidity_indoor_avg,
		max(humidity_indoor_max) as humidity_indoor_max,
		min(humidity_indoor_min) as humidity_indoor_min,
		avg(humidity_outdoor_avg) as humidity_outdoor_avg,
		max(humidity_outdoor_max) as humidity_outdoor_max,
		min(humidity_outdoor_min) as humidity_outdoor_min,
		avg(pressure_absolute_avg) as pressure_absolute_avg,
		max(pressure_absolute_max) as pressure_absolute_max,
		min(pressure_absolute_min) as pressure_absolute_min,
		avg(pressure_relative_avg) as pressure_relative_avg,
		max(pressure_relative_max) as pressure_relative_max,
		min(pressure_relative_min) as pressure_relative_min,
		avg(rain_avg) as rain_avg,
		max(rain_max) as rain_max,
		min(rain_min) as rain_min,
		avg(solar_radiation_avg) as solar_radiation_avg,
		max(solar_radiation_max) as solar_radiation_max,
		min(solar_radiation_min) as solar_radiation_min,
		avg(temperature_indoor_avg) as temperature_indoor_avg,
		max(temperature_indoor_max) as temperature_indoor_max,
		min(temperature_indoor_min) as temperature_indoor_min,
		avg(temperature_outdoor_avg) as temperature_outdoor_avg,
		max(temperature_outdoor_max) as temperature_outdoor_max,
		min(temperature_outdoor_min) as temperature_outdoor_min,
		avg(uv_avg) as uv_avg,
		max(uv_max) as uv_max,
		min(uv_min) as uv_min,
		avg(wind_chill_avg) as wind_chill_avg,
		max(wind_chill_max) as wind_chill_max,
		min(wind_chill_min) as wind_chill_min,
		avg(wind_gust_avg) as wind_gust_avg,
		max(wind_gust_max) as wind_gust_max,
		min(wind_gust_min) as wind_gust_min,
		avg(wind_direction_avg) as wind_direction_avg,
		avg(wind_speed_avg) as wind_speed_avg,
		max(wind_speed_max) as wind_speed_max,
		min(wind_speed_min) as wind_speed_min
	from daily
	group by
		date_trunc('month', daily.created at time zone '%[1]s')
	order by
		date_trunc('month', daily.created at time zone '%[1]s')
;

create index if not exists raw_created on raw using brin (created);
create index if not exists hourly_created on hourly using brin (created);
create index if not exists daily_created on daily using brin (created);
create index if not exists monthly_created on monthly using brin (created);
`, c.Timezone),
			Version: 1,
		},
	}

	sort.SliceStable(migrations, func(i, j int) bool { return migrations[i].Version < migrations[j].Version })
	for i := 0; i < len(migrations); i++ {
		if migrations[i].Version > version {
			tx := db.MustBegin()
			tx.MustExec(migrations[i].Migration)
			err := tx.Commit()

			if err != nil {
				logError(err.Error())
				return err
			}

			db.MustExec(`
insert into migration (
app,
version
) values (
'weatherproxy',
$1
) on conflict (
app
) do update set
version = $1
`, migrations[i].Version)
		}
	}

	return nil
}
