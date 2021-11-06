SHELL := /bin/bash

all: sales-api metrics

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:1.0 \
		--build-arg VCG_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%d%H:%M:%SZ"` \
		.

run:
	go run app/sales-api/main.go

runa:
	go run app/admin/main.go

tidy:
	go mod tidy
	go mod vendor

test:
	go test -v ./... -count=1
	staticcheck ./...
