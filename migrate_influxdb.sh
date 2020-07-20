#!/usr/bin/env bash

if [ -z ${INFLUXDB_ADDRESS} ]; then
  read -p "InfluxDB Address (https://influxdb.example.com:8086): " -r INFLUXDB_ADDRESS
fi

if [ -z ${INFLUXDB_DATABASE} ]; then
  read -p "InfluxDB Database: " -r INFLUXDB_DATABASE
fi

if [ -z ${INFLUXDB_USERNAME} ]; then
  read -p "InfluxDB Username: " -r INFLUXDB_USERNAME
fi

if [ -z ${INFLUXDB_PASSWORD} ]; then
  read -sp "InfluxDB Password: " -r INFLUXDB_PASSWORD
fi

if [ -z ${POSTGRESQL_DATABASE} ]; then
  read -sp "PostgreSQL Database: " -r POSTGRESQL_DATABASE
fi

if [ -z ${POSTGRESQL_HOSTNAME} ]; then
  read -sp "PostgreSQL Hostname: " -r POSTGRESQL_HOSTNAME
fi

if [ -z ${POSTGRESQL_USERNAME} ]; then
  read -sp "PostgreSQL Username: " -r POSTGRESQL_USERNAME
fi

if [ -z ${POSTGRESQL_PASSWORD} ]; then
  read -sp "PostgreSQL Username: " -r POSTGRESQL_PASSWORD
fi

echo

curl -k "${INFLUXDB_ADDRESS}/query?&u=${INFLUXDB_USERNAME}&p=${INFLUXDB_PASSWORD}" --data-urlencode "q=select * from ${INFLUXDB_DATABASE}.hourly.hourly" | jq --compact-output .results[0].series[0].values > /tmp/weatherproxy.json
psql "host=${POSTGRESQL_HOSTNAME} port=5432 dbname=${POSTGRESQL_DATABASE} user=${POSTGRESQL_USERNAME} password=${POSTGRESQL_PASSWORD}" << EOF
\set content $(cat /tmp/weatherproxy.json)
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
) (
select
  (elems->>0)::timestamp at time zone 'UTC' as created,
  (elems->>1)::real as dew_point_avg,
  (elems->>1)::real as dew_point_max,
  (elems->>1)::real as dew_point_min,
  (elems->>2)::real as humidity_indoor_avg,
  (elems->>2)::real as humidity_indoor_max,
  (elems->>2)::real as humidity_indoor_min,
  (elems->>3)::real as humidity_outdoor_avg,
  (elems->>3)::real as humidity_outdoor_max,
  (elems->>3)::real as humidity_outdoor_min,
  (elems->>4)::real as pressure_absolute_avg,
  (elems->>4)::real as pressure_absolute_max,
  (elems->>4)::real as pressure_absolute_min,
  0::real as pressure_absolute_min,
  0::real as pressure_relative_avg,
  0::real as pressure_relative_max,
  (elems->>5)::real as rain_avg,
  (elems->>5)::real as rain_max,
  (elems->>5)::real as rain_min,
  (elems->>6)::real as solar_radiation_avg,
  (elems->>6)::real as solar_radiation_max,
  (elems->>6)::real as solar_radiation_min,
  (elems->>7)::real as temperature_indoor_avg,
  (elems->>7)::real as temperature_indoor_max,
  (elems->>7)::real as temperature_indoor_min,
  (elems->>8)::real as temperature_outdoor_avg,
  (elems->>8)::real as temperature_outdoor_max,
  (elems->>8)::real as temperature_outdoor_min,
  (elems->>9)::real as uv_avg,
  (elems->>9)::real as uv_max,
  (elems->>9)::real as uv_min,
  (elems->>10)::real as wind_chill_avg,
  (elems->>10)::real as wind_chill_max,
  (elems->>10)::real as wind_chill_min,
  (elems->>11)::real as wind_direction_avg,
  (elems->>12)::real as wind_gust_avg,
  (elems->>12)::real as wind_gust_max,
  (elems->>12)::real as wind_gust_min,
  (elems->>13)::real as wind_speed_avg,
  (elems->>13)::real as wind_speed_max,
  (elems->>13)::real as wind_speed_min
from json_array_elements(:'content') as elems);
EOF
