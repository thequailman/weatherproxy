ENDPOINT ?= localhost:3000
SHELL := /bin/bash

BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_NAME := weatherproxy
BUILD_TAG ?= $(shell git describe --tags)

DOCKER_LOGOPTS = --log-opt max-file=1 --log-opt max-size=100k
DOCKER_POSTGRESQL = docker exec -it weatherproxy_postgresql psql -U postgres weatherproxy -c

WEATHERPROXY_POSTGRESQL_DATABASE ?= weatherproxy
WEATHERPROXY_POSTGRESQL_HOSTNAME ?= localhost
WEATHERPROXY_POSTGRESQL_PASSWORD ?= weatherproxy
WEATHERPROXY_POSTGRESQL_PORT ?= 5432
WEATHERPROXY_POSTGRESQL_USERNAME ?= weatherproxy

export GID = $(shell id -g)
export UID = $(shell id -u)

include bin.mk

.PHONY: build
build: $(BUILD_NAME)
	tar -czf $(BUILD_NAME).tar.gz $(BUILD_NAME)
	sha256sum -b $(BUILD_NAME) > $(BUILD_NAME).tar.gz.sha256

.PHONY: build_proxy
build_proxy: $(BUILD_NAME)

.PHONY: config
config:
	jq -n '{debug: true, postgresql: {database: "$(WEATHERPROXY_POSTGRESQL_DATABASE)", hostname: "$(WEATHERPROXY_POSTGRESQL_HOSTNAME)", password: "$(WEATHERPROXY_POSTGRESQL_PASSWORD)", port: $(WEATHERPROXY_POSTGRESQL_PORT), username: "$(WEATHERPROXY_POSTGRESQL_USERNAME)"}}' > config.json

clean_postgresql:
	$(DOCKER_POSTGRESQL) "drop owned by weatherproxy"
	$(DOCKER_POSTGRESQL) 'grant all privileges on schema public to "weatherproxy"'

.PHONY: get
get:
	curl -v 'http://$(ENDPOINT)/weatherstation/updateweatherstation.php?indoortempf=74.3&tempf=70.5&dewptf=68.5&windchillf=70.5&indoorhumidity=54&humidity=94&windspeedmph=2.2&windgustmph=2.2&winddir=77&absbaromin=29.12&baromin=29.89&rainin=0.07&dailyrainin=0.03&weeklyrainin=0.03&monthlyrainin=3.77&solarradiation=57.62&UV=0'

.PHONY: hash
hash: $(BUILD_NAME).sha256sum

.PHONY: lint
lint: golangci-lint
	golangci-lint run

.PHONY: postgresql
postgresql:
	mkdir -p .tmp/api
	docker run \
		-d \
		-e POSTGRES_PASSWORD=postgres \
		$(DOCKER_LOGOPTS) \
		--name weatherproxy_postgresql \
		-p 127.0.0.1:5432:5432 \
		--restart always \
		-v `pwd`/initdb.sql/:/docker-entrypoint-initdb.d/initdb.sql \
		-v weatherproxy_postgresql:/var/lib/postgresql/data \
		postgres:12 \
		-c log_statement=all

postgresql_stop:
	docker rm -f weatherproxy_postgresql || true
	docker volume rm weatherproxy_postgresql || true

test: go
	go test ./... -coverprofile coverage.out -p=1
	go tool cover -func=coverage.out

.PHONY: $(BUILD_NAME)
$(BUILD_NAME):
	go build -v -ldflags "-X main.buildDate=$(BUILD_DATE) -X main.buildCommit=$(BUILD_COMMIT) -X main.buildTag=$(BUILD_TAG)" -o $(BUILD_NAME)
