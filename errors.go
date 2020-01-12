package main

import (
	"errors"
)

// ErrUnableToConnect should be used for errors related to influxdb connectivity issues
var ErrUnableToConnect = errors.New("unable to connect to InfluxDB")

// ErrUnableToDecode should be used for errors related to decoding data.
var ErrUnableToDecode = errors.New("unable to decode")

// ErrUnableToEncode should be used for errors related to encoding data.
var ErrUnableToEncode = errors.New("unable to encode")

// ErrUnableToOpen is used when a config file doesn't exist.
var ErrUnableToOpen = errors.New("unable to open config")

// ErrUnableToRead is used when a config file contains invalid json.
var ErrUnableToRead = errors.New("unable to read config")

// ErrUnableToWrite should be used for errors related to influxdb writes
var ErrUnableToWrite = errors.New("unable to write to InfluxDB")
