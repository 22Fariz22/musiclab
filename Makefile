.PHONY: migrate migrate_down migrate_up migrate_version docker prod local swaggo test up down gen

# ==============================================================================
# Docker compose commands

FILES := $(shell docker ps -aq)

up:
	echo "Starting docker environment"
	docker compose -f docker-compose.dev.yml up --build

down:
	docker stop $(FILES)
	docker rm $(FILES)


# ==============================================================================
# Tools commands

swaggo:
	echo "Starting swagger generating"
	swag init -g **/**/*.go

gen:
	echo "Starting generate mock"
	go generate ./...

cover:
	echo "Create html file with cover data"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
