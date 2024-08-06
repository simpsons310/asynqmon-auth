DOCKER_TAG = simpsons310/asynqmon-auth
ENV_FILE = .env
DOCKER_APP_PORT = 8080

.PHONY: run
run:
	go run cmd/httpserver/main.go

.PHONY: build
build:
	go build -v -o ./build/asynqmon_auth cmd/httpserver/main.go

.PHONY: docker-build
docker-build:
	docker build -t ${DOCKER_TAG} .

.PHONY: docker-run
docker-run:
	docker run -p ${DOCKER_APP_PORT}:8080 --env-file=${ENV_FILE} ${DOCKER_TAG}
