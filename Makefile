OS := $(shell uname -s)
ifeq ($(OS),Linux)
    COMPOSE_FILE+= "docker-compose.linux.yaml"
endif
ifeq ($(OS),Darwin)
    COMPOSE_FILE+= "docker-compose.mac.yaml"
endif

run:
	docker-compose -f $(COMPOSE_FILE) up --build -d; \
	docker image prune -f

stop:
	docker-compose -f $(COMPOSE_FILE) down

# ---------------------
# Run migrations: Goose
# ---------------------

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" create rename_me sql

goose-up:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" up

goose-down:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" down

goose-status:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" status

