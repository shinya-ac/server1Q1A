DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_FILE = docker-compose.yml
DOCKER_COMPOSE_CMD = $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE)
DOCKER_NETWORK = 1q1a_api_network

.PHONY: all build up down restart logs test

all: build

create-network:
	@if [ $$(docker network ls --filter name=$(DOCKER_NETWORK) --format '{{.Name}}') != "$(DOCKER_NETWORK)" ]; then \
		echo "Creating network $(DOCKER_NETWORK)..."; \
		docker network create $(DOCKER_NETWORK); \
	else \
		echo "Network $(DOCKER_NETWORK) already exists."; \
	fi

build: create-network
	$(DOCKER_COMPOSE_CMD) build

build-no-cache: create-network
	$(DOCKER_COMPOSE_CMD) build --no-cache

up: create-network
	$(DOCKER_COMPOSE_CMD) up

down:
	$(DOCKER_COMPOSE_CMD) down

# サービス再起動
restart: down up

logs:
	$(DOCKER_COMPOSE_CMD) logs -f

test:
	cd api/src && go clean -testcache && go test ./...

ut:
	cd api/src && go clean -testcache && go test ./domain/task/...

# モック生成
# generate:
# 	$(DOCKER_COMPOSE_CMD) run --rm api sh -c "cd /app/src && go generate ./..."

# Swaggerドキュメント生成
swagger:
	cd api/src && swag init -g /cmd/main.go -o docs/

config:
	$(DOCKER_COMPOSE_CMD) config

# lint:
# 	$(DOCKER_COMPOSE_CMD) run --rm api sh -c "cd /app/src && golangci-lint run"

# ローカル環境の初期化
setup: create-network
	$(DOCKER_COMPOSE_CMD) up -d db
	# データベースが起動するまで待機
	sleep 10
	$(DOCKER_COMPOSE_CMD) up api
