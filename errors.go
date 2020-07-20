package main

import (
	"errors"
)

// ErrUnableToConnect should be used for errors related to postgresql connectivity issues
var ErrUnableToConnect = errors.New("unable to connect to PostgreSQL")

// ErrUnableToEncode should be used for errors related to encoding data.
var ErrUnableToEncode = errors.New("unable to encode")

// ErrUnableToOpen is used when a config file doesn't exist.
var ErrUnableToOpen = errors.New("unable to open config")

// ErrUnableToParse is used when a query can't be parsed.
var ErrUnableToParse = errors.New("unable to parse query")

// ErrUnableToRead is used when a config file contains invalid json.
var ErrUnableToRead = errors.New("unable to read config")
