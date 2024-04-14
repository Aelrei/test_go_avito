all: docker_compose

docker_compose:
	sudo docker-compose up --build main

docker_compose_up:
	sudo docker-compose up main

test:
	go test -v ./tests
