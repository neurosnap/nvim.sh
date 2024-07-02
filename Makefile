DOCKER_CMD?=docker
REGISTRY?=localhost:1338
GCLOUD?="gcr.io/google.com/cloudsdktool/google-cloud-cli"

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

fmt:
	go fmt ./...
.PHONY: fmt

auth:
	docker run -ti --name gcloud-config $(GCLOUD) gcloud auth login
.PHONY: auth

project:
	docker run --rm --volumes-from gcloud-config $(GCLOUD) gcloud config set project neovim-awesome
.PHONY: project

deploy:
	docker run --rm -v .:/app --volumes-from gcloud-config $(GCLOUD) sh -c "cd /app && gcloud app deploy"
.PHONY: deploy

sh:
	docker run --rm -it -v .:/app --volumes-from gcloud-config $(GCLOUD) sh
.PHONY: sh
