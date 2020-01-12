.PHONY: build config get hash lint test 

DOCKER = docker-compose run --rm
ENDPOINT ?= localhost:3000
GO = $(DOCKER) golang:1.13 go
SHELL := /bin/bash

BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_NAME := weatherproxy_linux
BUILD_TAG ?= $(shell git describe --tags)

WEATHERPROXY_INFLUXDB_DATABASE ?= dev
WEATHERPROXY_INFLUXDB_HOSTNAME ?= http://influxdb:8086
WEATHERPROXY_INFLUXDB_PASSWORD ?= dev
WEATHERPROXY_INFLUXDB_USERNAME ?= dev

export GID = $(shell id -g)
export UID = $(shell id -u)

build: $(BUILD_NAME)
	tar -czf $(BUILD_NAME).tar.gz $(BUILD_NAME)
	sha256sum -b $(BUILD_NAME) > $(BUILD_NAME).tar.gz.sha256
config: 
	jq -n '{debug: true, influxdb: {address: "$(WEATHERPROXY_INFLUXDB_HOSTNAME)", database: "$(WEATHERPROXY_INFLUXDB_DATABASE)", username: "$(WEATHERPROXY_INFLUXDB_USERNAME)", password: "$(WEATHERPROXY_INFLUXDB_PASSWORD)"}}' > config.json

get:
	curl -v 'http://$(ENDPOINT)/weatherstation/updateweatherstation.php?indoortempf=74.3&tempf=70.5&dewptf=68.5&windchillf=70.5&indoorhumidity=54&humidity=94&windspeedmph=2.2&windgustmph=2.2&winddir=77&absbaromin=29.12&baromin=29.89&rainin=0.07&dailyrainin=0.03&weeklyrainin=0.03&monthlyrainin=3.77&solarradiation=57.62&UV=0'

hash: $(BUILD_NAME).sha256sum

lint:
	$(DOCKER) golangci-lint

test: config
	docker-compose up -d influxdb
	while ! curl -s http://localhost:8086/ping; do \
		((c++)) && ((c==30)) && exit 1; \
		sleep 2; \
	done;
	$(DOCKER) go test ./... -coverprofile coverage.out -p=1
	$(DOCKER) go tool cover -func=coverage.out

$(BUILD_NAME):
	$(DOCKER) go build -v -ldflags "-X main.buildDate=$(BUILD_DATE) -X main.buildCommit=$(BUILD_COMMIT) -X main.buildTag=$(BUILD_TAG)" -o $(BUILD_NAME)
