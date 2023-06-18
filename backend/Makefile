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
		docker compose -f docker-compose.testing.yml down
	}
	trap tearDown EXIT
	docker compose -f docker-compose.testing.yml up --build -d
