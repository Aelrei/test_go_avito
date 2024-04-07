GO=go

all:
	FORCE

compile:
	export DB_HOST=localhost; \
	export DB_USER=aelrei; \
	export DB_PASSWORD=123; \
	export DB_NAME=aelrei; \
	export DB_PORT=5432; \
	$(GO) run ./cmd/main.go


docker_compose:
	sudo docker-compose up --build main

docker_compose_up:
	sudo docker-compose up main