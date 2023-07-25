export SHELL:=/bin/bash
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

.ONESHELL:
.PHONY: test
#   Usage:
#       make test norace=1
#           Expl: run WITHOUT race conditions
#       make test ftest=1
#           Explt: run WITH "ftests", long lived non unit tests
#       make test v=1
#           Explt: run in verbose mode
#       make test package=my_awesome_pkg
#           Explt: run tests in a single package only
#           NOTE: if omitted, will run ALL tests
#       make test package=my_awesome_pkg func=TestViews
#           Explt: run a single test func (needs package as well)
test:
	function tearDown {
		docker compose -f docker-compose.testing.yml down
	}
	trap tearDown EXIT
	docker compose -f docker-compose.testing.yml up --build -d
	$(eval ARGS=)
# by default, run without race conditions
ifndef race
	$(eval ARGS=${ARGS})
else
	$(eval ARGS=${ARGS} -race)
endif
# by default, run in "short" mode
ifndef ftest
	$(eval ARGS=${ARGS} -short)
endif
ifdef v
	$(eval ARGS=${ARGS} -v)
endif
# if package provided, run the package
ifdef package
	$(eval ARGS=${ARGS} ./${package}/...)
else
	$(eval ARGS=${ARGS} ./...)
endif
# if func provided, run the func only
ifdef func
	$(eval ARGS=${ARGS} -run=${func})
endif
	go test ${ARGS} -count=1 -cover

.PHONY: migrate-up
migrate-up:
	migrate -path=./migrations -database=postgres://cookie:cookie@localhost:5432/cookie?sslmode=disable up

.PHONY: migrate-down
migrate-down:
	migrate -path=./migrations -database=postgres://cookie:cookie@localhost:5432/cookie?sslmode=disable down

.PHONY: start-local-dev
start-local-dev:
	go build -o svr cmd/api/*.go
	./svr -db-dsn=postgres://cookie:cookie@localhost:5432/cookie?sslmode=disable

.PHONY: start-docker-dev
start-docker-dev:
	function tearDown {
		docker compose -f docker-compose.development.yml down
	}
	trap tearDown EXIT
	docker compose -f docker-compose.development.yml up --build

.PHONY: start-docker-testing
start-docker-testing:
	function tearDown {
		docker compose -f docker-compose.testing.yml down
	}
	trap tearDown EXIT
	docker compose -f docker-compose.testing.yml up --build

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

production_host_ip = "project4495.canadacentral.cloudapp.azure.com"

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh cookie@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api cookie@${production_host_ip}:~
	rsync -rP --delete ./migrations cookie@${production_host_ip}:~
	rsync -P ./remote/production/Caddyfile cookie@${production_host_ip}:~
	rsync -P ./remote/production/api.service cookie@${production_host_ip}:~
	ssh -t cookie@${production_host_ip} '\
		migrate -path ~/migrations -database $$COOKIE_DB_DSN up \
		&& sudo mv ~/api.service /etc/systemd/system/ \
		&& sudo systemctl enable api \
		&& sudo systemctl restart api \
		&& sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
	'