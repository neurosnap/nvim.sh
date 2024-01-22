DOCKER_CMD?=docker
REGISTRY?=registry.erock.io

setup:
	$(DOCKER_CMD) tag nvimsh $(REGISTRY)/nvimsh
.PHONY: setup

build:
	$(DOCKER_CMD) build -t $(REGISTRY)/nvimsh .
.PHONY: build

push:
	$(DOCKER_CMD) push $(REGISTRY)/nvimsh:latest
.PHONY: push

bp: build push
.PHONY: bp
