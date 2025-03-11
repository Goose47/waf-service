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
	@docker build -t waf-integration-tests-image ./tests >> /dev/null
	@docker run --network=waf-service_waf-network waf-integration-tests-image

lint:
	@echo "Running linter for analyzer service"
	@cd analyzer && golangci-lint run --config ../.golangci.yml && cd -
	@echo "Running linter for detection service"
	@cd detection && golangci-lint run --config ../.golangci.yml && cd -
	@echo "Running linter for waf service"
	@cd waf && golangci-lint run --config ../.golangci.yml && cd -

