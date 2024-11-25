.PHONY: migrate migrate_down migrate_up migrate_version docker prod local swaggo test

# ==============================================================================
# Docker compose commands

FILES := $(shell docker ps -aq)

dev:
	echo "Starting docker environment"
	docker-compose -f docker-compose.dev.yml up --build

dev-down:
	docker stop $(FILES)
	docker rm $(FILES)