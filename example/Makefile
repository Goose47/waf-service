.PHONY: build, run, stoprm

build:
	@docker build -t waf-example . && \
	docker image ls | head -n 1 && \
	docker image ls | grep waf-example

run:
	@docker run -p "8080:8080" -d --name "waf-example" waf-example && \
	docker ps | head -n 1 && \
	docker ps | grep waf-example

stoprm:
	@docker stop waf-example && \
	docker container rm waf-example