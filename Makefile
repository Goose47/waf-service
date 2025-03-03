.PHONY: run, update, test, lint

run:
	@docker compose up -d

update:
	@git pull && \
	docker compose build && \
	docker compose down && \
	docker compose up -d

test:
	@echo "Running integration tests"
	@cd tests && go test tests/... && cd -

lint:
	@echo "Running linter for analyzer service"
	@cd analyzer && golangci-lint run --config ../.golangci.yml && cd -
	@echo "Running linter for detection service"
	@cd detection && golangci-lint run --config ../.golangci.yml && cd -
	@echo "Running linter for waf service"
	@cd waf && golangci-lint run --config ../.golangci.yml && cd -

