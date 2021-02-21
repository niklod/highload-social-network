.PHONY: build test

build:
	go build -v ./cmd/highload-social-network

build-db:
	go build -v ./cmd/db-uploader

run:
	./build/highload-social-network

up:
	docker-compose -f deployment/docker-compose.yml up

down:
	docker-compose -f deployment/docker-compose.yml down