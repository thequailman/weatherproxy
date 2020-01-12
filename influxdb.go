package main

import (
	"github.com/influxdata/influxdb/client/v2"
)

var batchConfig client.BatchPointsConfig

var influxdb client.Client

func checkHealth() error {
	_, _, err := influxdb.Ping(10)
	return err
}

func getInfluxDB(c *config) error {
	var err error

	// Create InfluxDB client
	h := client.HTTPConfig{
		Addr:     c.InfluxDB.Address,
		Password: c.InfluxDB.Password,
		Username: c.InfluxDB.Username,
	}
	influxdb, err = client.NewHTTPClient(h)
	if err != nil {
		logError("%s: %s", ErrUnableToConnect, err)
		return ErrUnableToConnect
	}
	err = checkHealth()
	if err != nil {
		logError("%s: %s", ErrUnableToConnect, err)
		return ErrUnableToConnect
	}
	batchConfig = client.BatchPointsConfig{
		Precision: "s",
		Database:  c.InfluxDB.Database,
	}
	return nil
}
