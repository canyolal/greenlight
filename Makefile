include .envrc
## ---- ##
## HELPERS ##
## ---- ##

.PHONY: help
help:
	@echo 'Usage:'
	@sed 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## ---- ##
## DEVELOPMENT ##
## ---- ##

.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}

.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

## ---- ##
## QUALITY CONTROL ##
## ---- ##
.PHONY: audit
audit:
	@echo 'Formatting code'
	go fmt ./...
	@echo 'Vetting code'
	go vet ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

## ---- ##
## BUILD ##
## ---- ##

current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -w -X main.buildTime=${current_time} -X main.version=${git_description}'

# -s and -w for omitting DWARF table to reduce binary size
.PHONY: build/api
build/api:
	@echo 'Building cmd/api'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

## ---- ##
## PRODUCTION ##
## ---- ##

production_host_ip = '68.183.67.92'

## production/connect: connect to prod. server
.PHONY: production/connect
production/connect:
	ssh greenlight@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -rP --delete ./bin/linux_amd64/api ./migrations greenlight@${production_host_ip}:~
	ssh -t greenlight@${production_host_ip} 'migrate -path ~/migrations -database $$GREENLIGHT_DB_DSN up'

## production/configure/api.service: configure the production systemd api.service file
.PHONY: production/configure/api.service
production/configure/api.service:
	rsync -P ./remote/production/api.service greenlight@${production_host_ip}:~
	ssh -t greenlight@${production_host_ip} '\
	sudo mv ~/api.service /etc/systemd/system/ \
	&& sudo systemctl enable api \
	&& sudo systemctl restart api \
	'

## production/configure/caddyfile: configure the production Caddyfile
.PHONY: production/configure/caddyfile 
production/configure/caddyfile:
	rsync -P ./remote/production/Caddyfile greenlight@${production_host_ip}:~
	ssh -t greenlight@${production_host_ip} '\
	sudo mv ~/Caddyfile /etc/caddy/ \
	&& sudo systemctl reload caddy \
	'