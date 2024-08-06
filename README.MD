# Asynqmon Auth

A webserver with basic authentication middleware for package [Asynqmon](https://github.com/hibiken/asynqmon)

## Configuration

- Environment variables: checking `.env.example` or table below

```bash
# Server port
SERVER_PORT=8080
# Authentication mode for web-server. Accept modes:
# - `none` - disable authentication
# - `basic` - basic authentication
# - `http` - authentication via http (TODO)
SERVER_AUTH_MODE=basic
# Username & password for basic authentication
SERVER_AUTH_BASIC_USERNAME=
SERVER_AUTH_BASIC_PASSWORD=

# Asynqmon environment variables
ASYNQ_MON_ROOT_PATH=/
ASYNQ_READ_ONLY=false
ASYNQ_REDIS_DSN=redis://127.0.0.1:6379/0
ASYNQ_REDIS_INSECURE_TLS=false
```

## Running package

### Running in source code

``` bash
# Running with make
make run
```

### Running with binary

``` bash
# build binary
make build

# running with binary
./build/asynqmon_auth
```

- Binary accept first argument as server port. This argument will override environment variable `SERVER_PORT`

```bash
./build/asynqmon_auth 8080
```

### Running with docker

- Build/pull image

```bash
# Build docker image in local
make build

# Pull docker image from registry
docker pull simpsons310/asynqmon-auth
```

- Run image

```bash
# Docker command
docker run --rm \
    --name asynqmon-auth \
    --env-file .env
    -p 8080:8080
    simpsons310/asynqmon-auth

# With source code
make docker-run
```

## NOTE

This package is created just for practice purpose. I'm newbie to golang :smile:. I've referenced from [syahidfrd/asynqmon-handler](https://github.com/syahidfrd/asynqmon-handler)

## TODO

- Push image to docker registry
- Versioning
- Authentication via HTTP
